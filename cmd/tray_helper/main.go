package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/KusoKaihatsuSha/tray_helper/internal/handlers"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func init() {
	fmt.Print(
		fmt.Sprintf("Build version: %s\n", buildVersion),
		fmt.Sprintf("Build date: %s\n", buildDate),
		fmt.Sprintf("Build commit: %s\n", buildCommit),
	)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	handlers.Init(ctx, stop)
}
