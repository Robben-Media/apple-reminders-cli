package commands

import "github.com/spf13/cobra"

var (
	appVersion = "dev"
	buildTime  = "unknown"
)

func SetBuildInfo(version, built string) {
	if version != "" {
		appVersion = version
	}
	if built != "" {
		buildTime = built
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeJSON(map[string]string{
				"version":    appVersion,
				"build_time": buildTime,
			})
		},
	}
}
