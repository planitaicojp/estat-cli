package api

import (
	"fmt"

	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
	"github.com/planitaicojp/estat-cli/internal/model"
)

// StatsListResult is the unwrapped response for command usage.
type StatsListResult struct {
	Result      model.ResultInf
	DatalistInf model.DatalistInf
}

// GetStatsList calls the getStatsList API endpoint.
func GetStatsList(c *Client, params map[string]string) (*StatsListResult, error) {
	var resp model.StatsListResponse
	if err := c.Get("/json/getStatsList", params, &resp); err != nil {
		return nil, err
	}

	inner := resp.GetStatsList
	if inner.Result.Status != 0 {
		return nil, &cerrors.APIError{
			StatusCode: inner.Result.Status,
			Code:       fmt.Sprintf("%d", inner.Result.Status),
			Message:    inner.Result.ErrorMsg,
		}
	}

	return &StatsListResult{
		Result:      inner.Result,
		DatalistInf: inner.DatalistInf,
	}, nil
}
