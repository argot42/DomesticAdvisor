package main

import (
	"fmt"
	"github.com/argot42/money/config"
	"os"
)

func main() {
	cfg, err := config.GetConfig(os.Args)
	if err != nil {
		if err == config.ErrConfigFilePath {
			config.Usage()
			return
		}
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(1)
	}
	fmt.Println(cfg)
}
