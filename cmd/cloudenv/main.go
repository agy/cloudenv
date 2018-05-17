package main

import (
	"fmt"
	"os"
	"time"

	"github.com/agy/cloudenv"
)

func main() {
	timeout := 1 * time.Second

	cfg := cloudenv.Discover(timeout)

	if cfg == nil {
		fmt.Fprintln(os.Stderr, "invalid cloud provider or discovery timed out")
	}

	fmt.Println(cfg)
}
