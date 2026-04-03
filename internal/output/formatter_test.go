package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type testRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestTableFormatter(t *testing.T) {
	var buf bytes.Buffer
	rows := []testRow{
		{ID: "001", Name: "テスト1"},
		{ID: "002", Name: "テスト2"},
	}

	f := New("table")
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ID") {
		t.Errorf("table output missing header 'ID'")
	}
	if !strings.Contains(out, "001") {
		t.Errorf("table output missing value '001'")
	}
	if !strings.Contains(out, "テスト1") {
		t.Errorf("table output missing value 'テスト1'")
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	rows := []testRow{{ID: "001", Name: "テスト"}}

	f := New("json")
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var result []testRow
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if result[0].ID != "001" {
		t.Errorf("ID = %q, want %q", result[0].ID, "001")
	}
}

func TestCSVFormatter(t *testing.T) {
	var buf bytes.Buffer
	rows := []testRow{
		{ID: "001", Name: "テスト"},
	}

	f := New("csv")
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("CSV lines = %d, want 2 (header + 1 row)", len(lines))
	}
	if !strings.Contains(lines[0], "id") {
		t.Errorf("CSV header missing 'id': %q", lines[0])
	}
	if !strings.Contains(lines[1], "001") {
		t.Errorf("CSV row missing '001': %q", lines[1])
	}
}

func TestNew_default(t *testing.T) {
	f := New("unknown")
	if _, ok := f.(*TableFormatter); !ok {
		t.Errorf("New('unknown') should return TableFormatter")
	}
}

func TestTableFormatter_emptySlice(t *testing.T) {
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, []testRow{}); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("empty slice should produce no output, got %q", buf.String())
	}
}
