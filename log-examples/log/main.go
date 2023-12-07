package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("Cookie 🍪")
	log.Info("Hello World")

	err := fmt.Errorf("too much sugar")
	log.Error("failed to bake cookies", "err", err)

	log.Print("Baking ...")

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		//TimeFormat:      time.Kitchen,
		Prefix: "Baking 🍪",
	})

	logger.Info("Starting Oven!", "degree", 375)
	time.Sleep(1 * time.Second)
	logger.Info("Finished ...")
}
