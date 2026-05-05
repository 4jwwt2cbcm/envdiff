package diff

import (
	"strings"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected OutputFormat
	}{
		{"text", FormatText},
		{"TEXT", FormatText},
		{"", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"csv", FormatCSV},
		{"CSV", FormatCSV},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("got %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func sampleReport() Report {
	return Report{
		MissingInSecond: []string{"DB_HOST"},
		MissingInFirst:  []string{"NEW_KEY"},
		Mismatched: []MismatchedKey{
			{Key: "APP_ENV", First: "dev", Second: "prod"},
		},
	}
}

func TestFormatReport_Text(t *testing.T) {
	out := FormatReport(sampleReport(), FormatText)
	if !strings.Contains(out, "MISSING_IN_SECOND\tDB_HOST") {
		t.Errorf("expected MISSING_IN_SECOND line, got:\n%s", out)
	}
	if !strings.Contains(out, "MISSING_IN_FIRST\tNEW_KEY") {
		t.Errorf("expected MISSING_IN_FIRST line, got:\n%s", out)
	}
	if !strings.Contains(out, "MISMATCH\tAPP_ENV") {
		t.Errorf("expected MISMATCH line, got:\n%s", out)
	}
}

func TestFormatReport_JSON(t *testing.T) {
	out := FormatReport(sampleReport(), FormatJSON)
	if !strings.Contains(out, "\"missing_in_second\"") {
		t.Errorf("expected JSON key missing_in_second, got:\n%s", out)
	}
	if !strings.Contains(out, "\"DB_HOST\"") {
		t.Errorf("expected DB_HOST in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, "\"APP_ENV\"") {
		t.Errorf("expected APP_ENV in JSON output, got:\n%s", out)
	}
}

func TestFormatReport_CSV(t *testing.T) {
	out := FormatReport(sampleReport(), FormatCSV)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if lines[0] != "type,key,first,second" {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
	if len(lines) != 4 { // header + 3 entries
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestFormatReport_EmptyReport(t *testing.T) {
	r := Report{}
	if out := FormatReport(r, FormatText); out != "" {
		t.Errorf("expected empty text output, got: %q", out)
	}
	if out := FormatReport(r, FormatCSV); !strings.HasPrefix(out, "type,key") {
		t.Errorf("expected CSV header only, got: %q", out)
	}
}
