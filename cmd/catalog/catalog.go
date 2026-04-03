package catalog

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "catalog",
	Short: "データカタログを取得",
	Long:  "統計データのカタログ情報を取得します。",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat catalog は今後のバージョンで実装予定です")
	},
}

func init() {
	Cmd.Flags().String("survey", "", "調査名で絞り込み")
	Cmd.Flags().String("field", "", "統計分野コードで絞り込み")
	Cmd.Flags().String("dataset-type", "", "データセット種別（db または file）")
}
