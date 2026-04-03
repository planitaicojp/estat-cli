package dataset

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "dataset",
	Short: "データセットを管理",
	Long:  "データセットの登録・参照を行います。",
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "データセットを登録",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat dataset register は今後のバージョンで実装予定です")
	},
}

var showCmd = &cobra.Command{
	Use:   "show <datasetId>",
	Short: "データセットを参照",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat dataset show は今後のバージョンで実装予定です")
	},
}

func init() {
	registerCmd.Flags().String("stats-data-id", "", "統計表ID")
	registerCmd.Flags().String("filter", "", "フィルタ条件")
	Cmd.AddCommand(registerCmd)
	Cmd.AddCommand(showCmd)
}
