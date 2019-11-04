// +build integration

package eightpuzzle_test

import (
	eightpuzzleservice "github.com/ykamo001/ai/rpc/eightpuzzle"
	"testing"
)

var tests map[string]func(*testing.T)

func init() {
	tests = make(map[string]func(*testing.T))
	tests["FindPath"] = FindPath
}

func TestProvider(t *testing.T) {
	for name, test := range tests {
		t.Run(name, test)
	}
}

func FindPath(t *testing.T) {
	tests := map[string]func(*testing.T){
		"repeated values in initial matrix": func(t *testing.T) {
			ctx := makeContext()
			request := eightpuzzleservice.FindPathRequest{
				Initial: "1 1 2, 3 4 5, 6 7 8",
				Goal:    "1 2 3, 4 5 6, 7 8 0",
			}

			_, err := provider.FindPath(ctx, &request)
			if err == nil {
				t.Error("should have errored, but did not.")
			}
		},
		"repeated values in goal matrix": func(t *testing.T) {
			ctx := makeContext()
			request := eightpuzzleservice.FindPathRequest{
				Initial: "1 2 3, 4 5 6, 7 8 0",
				Goal:    "1 2 3, 4 5 6, 4 8 0",
			}

			_, err := provider.FindPath(ctx, &request)
			if err == nil {
				t.Error("should have errored, but did not.")
			}
		},
		"success": func(t *testing.T) {
			ctx := makeContext()
			request := eightpuzzleservice.FindPathRequest{
				Initial: "1 2 3, 4 5 6, 7 8 0",
				Goal:    "1 2 3, 4 5 6, 7 8 0",
			}

			_, err := provider.FindPath(ctx, &request)
			if err != nil {
				t.Errorf("should not have errored, but did. traceId:%v err:%v", ctx.Value("id"), err.Error())
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}