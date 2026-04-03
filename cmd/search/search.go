package search

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/planitaicojp/estat-cli/cmd/cmdutil"
	"github.com/planitaicojp/estat-cli/internal/api"
	"github.com/planitaicojp/estat-cli/internal/model"
	"github.com/planitaicojp/estat-cli/internal/output"
)

var Cmd = &cobra.Command{
	Use:   "search [キーワード]",
	Short: "統計表を検索",
	Long: `e-Stat APIで統計表を検索します。

使用例:
  estat search 人口
  estat search --survey "国勢調査"
  estat search --field "02" --open-year 2020
  estat search 人口 --limit 5 --format json`,
	RunE: runSearch,
}

func init() {
	Cmd.Flags().String("survey", "", "調査名で絞り込み")
	Cmd.Flags().String("field", "", "統計分野コードで絞り込み")
	Cmd.Flags().String("open-year", "", "公開年で絞り込み")
	Cmd.Flags().String("stats-code", "", "政府統計コードで絞り込み")
	Cmd.Flags().Int("limit", 0, "取得件数の上限")
	Cmd.Flags().Int("start", 0, "取得開始位置")
	Cmd.Flags().String("base-url", "", "APIベースURL（テスト用）")
	_ = Cmd.Flags().MarkHidden("base-url")
}

func runSearch(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient(cmd)
	if err != nil {
		return err
	}

	if baseURL, _ := cmd.Flags().GetString("base-url"); baseURL != "" {
		client.BaseURL = baseURL
	}

	params := buildParams(cmd, args)

	result, err := api.GetStatsList(client, params)
	if err != nil {
		return err
	}

	rows := model.ToTableRows(result.DatalistInf.TableInf)

	w := cmd.OutOrStdout()
	if w == os.Stdout {
		w = os.Stdout
	}
	format := cmdutil.GetFormat(cmd)
	return output.New(format).Format(w, rows)
}

func buildParams(cmd *cobra.Command, args []string) map[string]string {
	params := make(map[string]string)

	if len(args) > 0 {
		params["searchWord"] = args[0]
	}
	if v, _ := cmd.Flags().GetString("survey"); v != "" {
		params["surveyYears"] = v
	}
	if v, _ := cmd.Flags().GetString("field"); v != "" {
		params["statsField"] = v
	}
	if v, _ := cmd.Flags().GetString("open-year"); v != "" {
		params["openYears"] = v
	}
	if v, _ := cmd.Flags().GetString("stats-code"); v != "" {
		params["statsCode"] = v
	}
	if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
		params["limit"] = fmt.Sprintf("%d", v)
	}
	if v, _ := cmd.Flags().GetInt("start"); v > 0 {
		params["startPosition"] = fmt.Sprintf("%d", v)
	}

	return params
}
