// +build integration

package paint_test

import (
	"fmt"
	paintservice "github.com/ykamo001/ai/rpc/paint"
	"testing"
)

var tests map[string]func(*testing.T)
var verbose bool

func init() {
	tests = make(map[string]func(*testing.T))
	tests["FindPath"] = FindPath

	verbose = false
}

func TestProvider(t *testing.T) {
	for name, test := range tests {
		t.Run(name, test)
	}
}

func FindPath(t *testing.T) {
	tests := map[string]func(*testing.T){
		"invalid x input": func(t *testing.T) {
			ctx := makeContext()
			request := paintservice.FillInRequest{
				Matrix: []*paintservice.InternalArray{
					{Array: []string{"X", "_", "X", "X", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "X", "X", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "X", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
				},
				Value: "y",
				X:     -2,
				Y:     2,
			}

			_, err := provider.FillIn(ctx, &request)
			if err == nil {
				t.Error("should have errored, but did not.")
			}
		},
		"invalid y input": func(t *testing.T) {
			ctx := makeContext()
			request := paintservice.FillInRequest{
				Matrix: []*paintservice.InternalArray{
					{Array: []string{"X", "_", "X", "X", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "X", "X", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "X", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
				},
				Value: "y",
				X:     2,
				Y:     -2,
			}

			_, err := provider.FillIn(ctx, &request)
			if err == nil {
				t.Error("should have errored, but did not.")
			}
		},
		"invalid x and y input": func(t *testing.T) {
			ctx := makeContext()
			request := paintservice.FillInRequest{
				Matrix: []*paintservice.InternalArray{
					{Array: []string{"X", "_", "X", "X", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "X", "X", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "X", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
				},
				Value: "y",
				X:     -2,
				Y:     -2,
			}

			_, err := provider.FillIn(ctx, &request)
			if err == nil {
				t.Error("should have errored, but did not.")
			}
		},
		"success": func(t *testing.T) {
			ctx := makeContext()
			request := paintservice.FillInRequest{
				Matrix: []*paintservice.InternalArray{
					{Array: []string{"X", "_", "X", "X", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"X", "_", "_", "_", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"X", "X", "X", "X", "_", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "X", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "X", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
					{Array: []string{"_", "_", "_", "_", "X", "_", "_", "_", "_", "_"}},
				},
				Value: "y",
				X:     2,
				Y:     2,
			}

			res, err := provider.FillIn(ctx, &request)
			if err != nil {
				t.Errorf("should not have errored, but did. traceId:%v err:%v", ctx.Value("id"), err.Error())
			} else {
				if !verbose {
					return
				}
				fmt.Println("Before")
				printPicture(request.Matrix)
				fmt.Println("After")
				printPicture(res.Matrix)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}

func printPicture(matrix []*paintservice.InternalArray) {
	for i, row := range matrix {
		if i == 0 {
			fmt.Print("|")
			for l := 0; l < len(row.Array); l++ {
				fmt.Print("-")
			}
			fmt.Println("|")
		}
		fmt.Print("|")
		for _, value := range row.Array {
			fmt.Printf("%v", value)
		}
		fmt.Println("|")
		if i == len(matrix)-1 {
			fmt.Print("|")
			for l := 0; l < len(row.Array); l++ {
				fmt.Print("-")
			}
			fmt.Println("|")
		}
	}
}
