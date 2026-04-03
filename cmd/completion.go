package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "シェル補完スクリプトを生成",
	Long: `指定したシェルの補完スクリプトを標準出力に出力します。

使用例:
  # bash
  estat completion bash > /etc/bash_completion.d/estat

  # zsh
  estat completion zsh > "${fpath[1]}/_estat"

  # fish
  estat completion fish > ~/.config/fish/completions/estat.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
