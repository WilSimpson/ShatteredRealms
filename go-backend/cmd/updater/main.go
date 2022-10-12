package main

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/updater"
)

func main() {
	err := internal.SetupConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	updater.Execute()
}
