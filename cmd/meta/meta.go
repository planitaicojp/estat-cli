package meta

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "meta <statsDataId>",
	Short: "メタ情報を取得",
	Long:  "指定した統計表IDのメタ情報（分類項目、地域コードなど）を取得します。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat meta は今後のバージョンで実装予定です")
	},
}

func init() {
	Cmd.Flags().String("class", "", "特定の分類のみ取得（例: area, time）")
}
