package get

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get <statsDataId>",
	Short: "統計データを取得",
	Long:  "指定した統計表IDの統計データを取得します。",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat get は今後のバージョンで実装予定です")
	},
}

func init() {
	Cmd.Flags().String("area", "", "地域コード")
	Cmd.Flags().String("time", "", "時間軸コード")
	Cmd.Flags().String("category", "", "分類項目フィルタ（例: cat01=A01）")
	Cmd.Flags().Bool("section-header", false, "セクションヘッダを含める")
	Cmd.Flags().Bool("bulk", false, "一括取得モード")
}
