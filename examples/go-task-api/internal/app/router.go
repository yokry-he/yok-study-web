package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
	taskdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/task"
	userdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

const maxRouterRequestBytes int64 = 1 << 20

type Pinger interface {
	PingContext(context.Context) error
}

type Readiness struct {
	db           Pinger
	shuttingDown atomic.Bool
}

func NewReadiness(db Pinger) *Readiness {
	return &Readiness{db: db}
}

func (r *Readiness) StartShutdown() {
	if r != nil {
		r.shuttingDown.Store(true)
	}
}

func (r *Readiness) Live(w http.ResponseWriter, request *http.Request) {
	_ = httpx.WriteData(w, http.StatusOK, map[string]string{"status": "alive"}, httpx.RequestIDFromContext(request.Context()))
}

func (r *Readiness) Ready(w http.ResponseWriter, request *http.Request) {
	requestID := httpx.RequestIDFromContext(request.Context())
	if r == nil || r.shuttingDown.Load() || r.db == nil {
		_ = httpx.WriteError(w, readinessError(), requestID)
		return
	}
	if err := r.db.PingContext(request.Context()); err != nil {
		_ = httpx.WriteError(w, readinessError(), requestID)
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, map[string]string{"status": "ready"}, requestID)
}

func readinessError() error {
	return httpx.NewAPIError(http.StatusServiceUnavailable, "NOT_READY", "服务尚未就绪", nil)
}

type routeSpec struct {
	method  string
	path    string
	handler http.HandlerFunc
}

type routeDispatch struct {
	mux         *http.ServeMux
	pathMux     *http.ServeMux
	pathMethods map[string][]string
}

func NewRouter(
	logger *slog.Logger,
	readiness *Readiness,
	users *userdomain.Handler,
	tasks *taskdomain.Handler,
	requestTimeout time.Duration,
) http.Handler {
	routes := []routeSpec{
		{method: http.MethodGet, path: "/health/live", handler: readiness.Live},
		{method: http.MethodGet, path: "/health/ready", handler: readiness.Ready},
		{method: http.MethodGet, path: "/api/users", handler: users.List},
		{method: http.MethodPost, path: "/api/users", handler: users.Create},
		{method: http.MethodGet, path: "/api/users/{id}", handler: users.Get},
		{method: http.MethodPatch, path: "/api/users/{id}/status", handler: users.ChangeStatus},
		{method: http.MethodGet, path: "/api/tasks", handler: tasks.List},
		{method: http.MethodPost, path: "/api/tasks", handler: tasks.Create},
		{method: http.MethodGet, path: "/api/tasks/{id}", handler: tasks.Get},
		{method: http.MethodPut, path: "/api/tasks/{id}", handler: tasks.Update},
		{method: http.MethodPatch, path: "/api/tasks/{id}/status", handler: tasks.ChangeStatus},
		{method: http.MethodDelete, path: "/api/tasks/{id}", handler: tasks.Delete},
	}

	mux := http.NewServeMux()
	pathMux := http.NewServeMux()
	pathMethods := make(map[string][]string)
	for _, route := range routes {
		mux.HandleFunc(route.method+" "+route.path, route.handler)
		if _, exists := pathMethods[route.path]; !exists {
			pathMux.Handle(route.path, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		}
		pathMethods[route.path] = append(pathMethods[route.path], route.method)
	}

	dispatch := &routeDispatch{mux: mux, pathMux: pathMux, pathMethods: pathMethods}
	return httpx.Chain(
		dispatch,
		httpx.RequestID,
		httpx.AccessLog(logger),
		httpx.Recover(logger),
		httpx.Deadline(requestTimeout),
		httpx.RequireJSON(maxRouterRequestBytes),
	)
}

func (d *routeDispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d == nil || d.mux == nil || d.pathMux == nil {
		_ = httpx.WriteError(w, errors.New("router is not initialized"), httpx.RequestIDFromContext(r.Context()))
		return
	}
	_, pathPattern := d.pathMux.Handler(r)
	if pathPattern == "" {
		NotFound().ServeHTTP(w, r)
		return
	}
	methods := d.pathMethods[pathPattern]
	if !allowsMethod(methods, r.Method) {
		MethodNotAllowed(methods...).ServeHTTP(w, r)
		return
	}
	d.mux.ServeHTTP(w, r)
}

func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httpx.WriteError(w, httpx.NewAPIError(
			http.StatusNotFound,
			httpx.CodeNotFound,
			"接口不存在",
			nil,
		), httpx.RequestIDFromContext(r.Context()))
	})
}

func MethodNotAllowed(methods ...string) http.Handler {
	allow := stableMethods(methods)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		_ = httpx.WriteError(w, httpx.NewAPIError(
			http.StatusMethodNotAllowed,
			httpx.CodeMethodNotAllowed,
			"请求方法不受支持",
			nil,
		), httpx.RequestIDFromContext(r.Context()))
	})
}

func allowsMethod(methods []string, method string) bool {
	for _, allowed := range methods {
		if method == allowed || (method == http.MethodHead && allowed == http.MethodGet) {
			return true
		}
	}
	return false
}

func stableMethods(methods []string) []string {
	order := map[string]int{
		http.MethodGet: 0, http.MethodPost: 1, http.MethodPut: 2,
		http.MethodPatch: 3, http.MethodDelete: 4,
	}
	seen := make(map[string]struct{}, len(methods))
	result := make([]string, 0, len(methods))
	for _, method := range methods {
		if _, exists := seen[method]; exists || method == "" {
			continue
		}
		seen[method] = struct{}{}
		result = append(result, method)
	}
	sort.Slice(result, func(i, j int) bool {
		left, leftKnown := order[result[i]]
		right, rightKnown := order[result[j]]
		if leftKnown && rightKnown {
			return left < right
		}
		if leftKnown != rightKnown {
			return leftKnown
		}
		return result[i] < result[j]
	})
	return result
}
