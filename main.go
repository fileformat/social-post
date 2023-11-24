/*
Copyright © 2023 fileformat@gmail.com
*/
package main

import (
	"github.com/fileformat/social-post/cmd"
)

var (
	version = "0.0.0"
	commit  = "local"
	date    = "local"
	builtBy = "unknown"
)

func main() {
	cmd.SetVersionInfo(cmd.VersionInfo{Commit: commit, Version: version, LastMod: date, Builder: builtBy})

	cmd.Execute()

}
