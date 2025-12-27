package main

import (
	"fmt"
	"os"

	"github.com/YashIIT0909/TRexT/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TRexT: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TRexT: %v\n", err)
		os.Exit(1)
	}
}
