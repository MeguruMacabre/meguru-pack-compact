package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/client"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: client <manifest-url>")
		os.Exit(2)
	}

	ctx := context.Background()

	if err := client.Run(ctx, os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
