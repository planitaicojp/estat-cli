package model

// StatsListResponse is the top-level response from getStatsList.
type StatsListResponse struct {
	GetStatsList struct {
		Result      ResultInf   `json:"RESULT"`
		Parameter   any         `json:"PARAMETER"`
		DatalistInf DatalistInf `json:"DATALIST_INF"`
	} `json:"GET_STATS_LIST"`
}

// ResultInf is the common result section in API responses.
type ResultInf struct {
	Status   int    `json:"STATUS"`
	ErrorMsg string `json:"ERROR_MSG"`
	Date     string `json:"DATE"`
}

// DatalistInf contains the list of statistical tables.
type DatalistInf struct {
	Number    int         `json:"NUMBER"`
	ResultInf PageInfo    `json:"RESULT_INF"`
	TableInf  []TableInfo `json:"TABLE_INF"`
}

// PageInfo contains pagination information.
type PageInfo struct {
	FromNumber int `json:"FROM_NUMBER"`
	ToNumber   int `json:"TO_NUMBER"`
	NextKey    int `json:"NEXT_KEY,omitempty"`
}

// TableInfo represents a single statistical table entry.
type TableInfo struct {
	ID                 string       `json:"@id"`
	StatName           CodeNamePair `json:"STAT_NAME"`
	GovOrg             CodeNamePair `json:"GOV_ORG"`
	StatisticsName     string       `json:"STATISTICS_NAME"`
	Title              NoNamePair   `json:"TITLE"`
	SurveyDate         string       `json:"SURVEY_DATE"`
	OpenDate           string       `json:"OPEN_DATE"`
	OverallTotalNumber int          `json:"OVERALL_TOTAL_NUMBER"`
	UpdatedDate        string       `json:"UPDATED_DATE"`
}

// CodeNamePair represents a JSON object with @code and $ fields.
type CodeNamePair struct {
	Code string `json:"@code"`
	Name string `json:"$"`
}

// NoNamePair represents a JSON object with @no and $ fields.
type NoNamePair struct {
	No   string `json:"@no"`
	Name string `json:"$"`
}

// TableRow is a flattened row for output formatting.
type TableRow struct {
	ID         string `json:"id"`
	StatName   string `json:"stat_name"`
	Title      string `json:"title"`
	SurveyDate string `json:"survey_date"`
	OpenDate   string `json:"open_date"`
}

// ToTableRows converts a list of TableInfo to output-ready rows.
func ToTableRows(tables []TableInfo) []TableRow {
	rows := make([]TableRow, len(tables))
	for i, t := range tables {
		rows[i] = TableRow{
			ID:         t.ID,
			StatName:   t.StatName.Name,
			Title:      t.Title.Name,
			SurveyDate: t.SurveyDate,
			OpenDate:   t.OpenDate,
		}
	}
	return rows
}
