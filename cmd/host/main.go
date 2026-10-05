package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/host"
)

func main() {
	ctx := context.Background()

	err := host.Run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
