package main

import (
	"context"
	"os"
	"os/exec"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {

}

func main() {
	ctx = context.TODO()
	ctx, cancel = context.WithCancel(ctx)
	switch os.Args[1] {
	case "build":
		exec.CommandContext(ctx, "go", "build", "-x", "-v")
	}
}
