package cmdutil

import (
	"github.com/spf13/cobra"
	"github.com/planitaicojp/estat-cli/internal/config"
)

func GetFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	if f := config.EnvOr(config.EnvFormat, ""); f != "" {
		return f
	}
	cfg, err := config.Load()
	if err != nil {
		return config.DefaultFormat
	}
	return cfg.Format
}
