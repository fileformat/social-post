package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	UserAgent string
	vi        VersionInfo
)

type VersionInfo struct {
	Commit  string `json:"commit"`
	LastMod string `json:"lastmod"`
	Version string `json:"version"`
	Builder string `json:"builder"`
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "version",
	Short: "Version & build info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("social-post v%s (%s)\n", vi.Version, vi.LastMod)
	},
}

func InitVersion(versionInfo VersionInfo) {
	rootCmd.AddCommand(versionCmd)
	vi = versionInfo
	UserAgent = fmt.Sprintf("social-post/%s (%s/%s)", vi.Version, vi.Commit, vi.LastMod)
}
