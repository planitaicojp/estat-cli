package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/planitaicojp/estat-cli/cmd/catalog"
	"github.com/planitaicojp/estat-cli/cmd/dataset"
	"github.com/planitaicojp/estat-cli/cmd/get"
	"github.com/planitaicojp/estat-cli/cmd/meta"
	"github.com/planitaicojp/estat-cli/cmd/search"
	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "estat",
	Short:         "e-Stat API CLI",
	Long:          "e-Stat（政府統計の総合窓口）APIを操作するコマンドラインツール",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().String("app-id", "", "アプリケーションID")
	rootCmd.PersistentFlags().String("format", "", "出力形式: table, json, csv")
	rootCmd.PersistentFlags().String("lang", "", "言語: J(日本語), E(英語)")
	rootCmd.PersistentFlags().Bool("no-color", false, "色出力を無効にする")
	rootCmd.PersistentFlags().Bool("verbose", false, "詳細出力")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(meta.Cmd)
	rootCmd.AddCommand(dataset.Cmd)
	rootCmd.AddCommand(catalog.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
