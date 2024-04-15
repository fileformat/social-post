/*
Copyright © 2023 fileformat@gmail.com
*/
package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/fileformat/social-post/cmd"
)

var (
	version = "0.0.0"
	commit  = "local"
	date    = "local"
	builtBy = "unknown"
)

func initLogger() {

	logLevel := os.Getenv("LOG_LEVEL")
	lvl := slog.LevelInfo
	var err error
	if logLevel != "" {
		ilvl, atoiErr := strconv.Atoi(logLevel)
		if atoiErr == nil {
			lvl = slog.Level(ilvl)
		} else {
			err = lvl.UnmarshalText([]byte(logLevel))
			if err != nil {
				lvl = slog.LevelInfo
			}
		}
	}

	// get log format from the environment
	logFormat := os.Getenv("LOG_FORMAT")
	var handler slog.Handler
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	default:
		//LATER: how to set level???
		// this seems way too hard: https://cs.opensource.google/go/go/+/refs/tags/go1.21.4:src/log/slog/example_level_handler_test.go
		//defaultHandler := slog.Default().Handler()
		//handler = defaultHandler
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err != nil {
		slog.Error("unable to set log level", "error", err)
	}

}

func main() {
	initLogger()
	cmd.InitVersion(cmd.VersionInfo{Commit: commit, Version: version, LastMod: date, Builder: builtBy})
	//cmd.InitRoot()
	cmd.InitEmail()
	cmd.InitFacebook()
	cmd.InitMastodon()
	cmd.InitSlack()

	cmd.Execute()

}
