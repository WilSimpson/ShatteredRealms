package main

import (
	"fmt"
	"github.com/ShatteredRealms/UpdaterCLI/internal"
	"github.com/ShatteredRealms/UpdaterCLI/pkg/updater"
)

func main() {
	err := internal.SetupConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	updater.Execute()
}
