package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatsList(t *testing.T) {
	responseJSON := `{
		"GET_STATS_LIST": {
			"RESULT": {
				"STATUS": 0,
				"ERROR_MSG": "正常に終了しました。",
				"DATE": "2026-04-03T10:00:00.000+09:00"
			},
			"PARAMETER": {
				"LANG": "J",
				"SEARCH_WORD": "人口"
			},
			"DATALIST_INF": {
				"NUMBER": 1,
				"RESULT_INF": {
					"FROM_NUMBER": 1,
					"TO_NUMBER": 1
				},
				"TABLE_INF": [
					{
						"@id": "0003410379",
						"STAT_NAME": {"@code": "00200521", "$": "国勢調査"},
						"GOV_ORG": {"@code": "00200", "$": "総務省"},
						"STATISTICS_NAME": "国勢調査 人口等基本集計",
						"TITLE": {"@no": "001", "$": "男女別人口－全国，都道府県"},
						"SURVEY_DATE": "202001",
						"OPEN_DATE": "2021-11-30",
						"OVERALL_TOTAL_NUMBER": 12345,
						"UPDATED_DATE": "2022-01-15"
					}
				]
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-id", "J")
	resp, err := GetStatsList(client, map[string]string{"searchWord": "人口"})
	if err != nil {
		t.Fatalf("GetStatsList() error: %v", err)
	}

	if resp.Result.Status != 0 {
		t.Errorf("Status = %d, want 0", resp.Result.Status)
	}

	tables := resp.DatalistInf.TableInf
	if len(tables) != 1 {
		t.Fatalf("len(TableInf) = %d, want 1", len(tables))
	}

	table := tables[0]
	if table.ID != "0003410379" {
		t.Errorf("ID = %q, want %q", table.ID, "0003410379")
	}
	if table.StatName.Name != "国勢調査" {
		t.Errorf("StatName = %q, want %q", table.StatName.Name, "国勢調査")
	}
	if table.Title.Name != "男女別人口－全国，都道府県" {
		t.Errorf("Title = %q, want %q", table.Title.Name, "男女別人口－全国，都道府県")
	}
}

func TestGetStatsList_apiError(t *testing.T) {
	responseJSON := `{
		"GET_STATS_LIST": {
			"RESULT": {
				"STATUS": 100,
				"ERROR_MSG": "アプリケーションIDが不正です。",
				"DATE": "2026-04-03T10:00:00.000+09:00"
			},
			"PARAMETER": {}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-id", "J")
	_, err := GetStatsList(client, nil)
	if err == nil {
		t.Fatal("GetStatsList() should return error for STATUS != 0")
	}
}
