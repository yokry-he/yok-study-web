package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := run(ctx, os.Getenv("HEALTHCHECK_URL"), &http.Client{Timeout: 2 * time.Second}); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, rawURL string, client *http.Client) error {
	if ctx == nil {
		return errors.New("健康检查 context 不能为空")
	}
	url := strings.TrimSpace(rawURL)
	if url == "" {
		return errors.New("HEALTHCHECK_URL 不能为空")
	}
	if client == nil {
		return errors.New("健康检查 HTTP client 不能为空")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.New("HEALTHCHECK_URL 格式无效")
	}
	response, err := client.Do(request)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return errors.New("健康检查请求失败")
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4<<10))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("健康检查返回非成功状态：%d", response.StatusCode)
	}
	return nil
}
