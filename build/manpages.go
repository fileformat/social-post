package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fileformat/social-post/cmd"
	mango "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
)

var (
	version = "0.0.0"
	commit  = "local"
	date    = "local"
	builtBy = "unknown"
)

func main() {

	cmd.InitVersion(cmd.VersionInfo{Commit: commit, Version: version, LastMod: date, Builder: builtBy})
	//cmd.InitRoot()
	cmd.InitEmail()
	cmd.InitFacebook()
	cmd.InitMastodon()

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	})

	slogger := slog.New(handler)
	slog.SetDefault(slogger)

	manPage, mangoErr := mango.NewManPage(1, cmd.GetRoot())
	if mangoErr != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to generate man page: %v", mangoErr)
		os.Exit(1)
	}

	_, _ = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
}
