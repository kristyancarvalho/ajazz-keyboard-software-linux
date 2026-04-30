package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/kristyancarvalho/ak820pro/internal/gui"
)

func main() {
	go func() {
		if err := run(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run() error {
	logger, closeLogger, err := newLogger()
	if err != nil {
		return err
	}
	defer closeLogger()

	var window app.Window
	window.Option(
		app.Title("Ajazz AK820 Pro"),
		app.Size(unit.Dp(900), unit.Dp(600)),
	)

	application, err := gui.New(&window, logger)
	if err != nil {
		return err
	}
	return application.Run()
}

func newLogger() (*slog.Logger, func() error, error) {
	logDir, err := logDirectory()
	if err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, nil, err
	}
	file, err := os.OpenFile(filepath.Join(logDir, "ak820pro.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return slog.New(slog.NewTextHandler(file, nil)), file.Close, nil
}

func logDirectory() (string, error) {
	if dir := os.Getenv("AK820PRO_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "ak820pro"), nil
}