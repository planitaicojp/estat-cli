package cmdutil

import (
	"github.com/spf13/cobra"
	"github.com/planitaicojp/estat-cli/internal/api"
	"github.com/planitaicojp/estat-cli/internal/config"
	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
)

func NewClient(cmd *cobra.Command) (*api.Client, error) {
	appID, _ := cmd.Flags().GetString("app-id")
	if appID == "" {
		appID = config.EnvOr(config.EnvAppID, "")
	}
	if appID == "" {
		cfg, err := config.Load()
		if err != nil {
			return nil, err
		}
		appID = cfg.AppID
	}
	if appID == "" {
		return nil, &cerrors.ConfigError{
			Message: "アプリケーションIDが設定されていません。--app-id フラグ、ESTAT_APP_ID 環境変数、または ~/.config/estat/config.yaml で設定してください",
		}
	}

	lang, _ := cmd.Flags().GetString("lang")
	if lang == "" {
		lang = config.EnvOr(config.EnvLang, "")
	}
	if lang == "" {
		cfg, err := config.Load()
		if err != nil {
			lang = config.DefaultLang
		} else {
			lang = cfg.Lang
		}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	client := api.NewClient("https://api.e-stat.go.jp/rest/3.0/app", appID, lang)
	client.Verbose = verbose
	return client, nil
}
