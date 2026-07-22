package main

import (
	"errors"
	"testing"
)

func TestParseOperation(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "up", args: []string{"up"}, want: "up"},
		{name: "down one", args: []string{"down", "1"}, want: "down-one"},
		{name: "version", args: []string{"version"}, want: "version"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseOperation(test.args)
			if err != nil || got != test.want {
				t.Fatalf("parseOperation(%v) = (%q, %v)", test.args, got, err)
			}
		})
	}
}

func TestParseOperationRejectsOtherArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"down"}, {"down", "2"}, {"up", "extra"}, {"unknown"}} {
		if _, err := parseOperation(args); !errors.Is(err, errUsage) {
			t.Fatalf("parseOperation(%v) error = %v, want errUsage", args, err)
		}
	}
}
