package diff

import (
	"fmt"
	"strings"
)

// OutputFormat represents the output format for diff results.
type OutputFormat int

const (
	FormatText OutputFormat = iota
	FormatJSON
	FormatCSV
)

// ParseFormat converts a string to an OutputFormat.
func ParseFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	default:
		return FormatText, fmt.Errorf("unknown format %q: must be one of text, json, csv", s)
	}
}

// FormatReport renders a Report in the requested OutputFormat and returns the result as a string.
func FormatReport(r Report, format OutputFormat) string {
	switch format {
	case FormatJSON:
		return formatJSON(r)
	case FormatCSV:
		return formatCSV(r)
	default:
		return formatText(r)
	}
}

func formatText(r Report) string {
	var sb strings.Builder
	for _, k := range r.MissingInSecond {
		fmt.Fprintf(&sb, "MISSING_IN_SECOND\t%s\n", k)
	}
	for _, k := range r.MissingInFirst {
		fmt.Fprintf(&sb, "MISSING_IN_FIRST\t%s\n", k)
	}
	for _, m := range r.Mismatched {
		fmt.Fprintf(&sb, "MISMATCH\t%s\t%s\t%s\n", m.Key, m.First, m.Second)
	}
	return sb.String()
}

func formatJSON(r Report) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString("  \"missing_in_second\": [")
	for i, k := range r.MissingInSecond {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "%q", k)
	}
	sb.WriteString("],\n")
	sb.WriteString("  \"missing_in_first\": [")
	for i, k := range r.MissingInFirst {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "%q", k)
	}
	sb.WriteString("],\n")
	sb.WriteString("  \"mismatched\": [")
	for i, m := range r.Mismatched {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "{\"key\": %q, \"first\": %q, \"second\": %q}", m.Key, m.First, m.Second)
	}
	sb.WriteString("]\n}")
	return sb.String()
}

func formatCSV(r Report) string {
	var sb strings.Builder
	sb.WriteString("type,key,first,second\n")
	for _, k := range r.MissingInSecond {
		fmt.Fprintf(&sb, "MISSING_IN_SECOND,%s,,\n", k)
	}
	for _, k := range r.MissingInFirst {
		fmt.Fprintf(&sb, "MISSING_IN_FIRST,%s,,\n", k)
	}
	for _, m := range r.Mismatched {
		fmt.Fprintf(&sb, "MISMATCH,%s,%s,%s\n", m.Key, m.First, m.Second)
	}
	return sb.String()
}
