package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/meesakveld/context/internal/scanner"
)

func Format(
	result scanner.ContextResult,
	format string,
) ([]byte, error) {
	switch strings.ToLower(format) {
	case "txt", "text":
		return formatText(result), nil

	case "markdown", "md":
		return formatMarkdown(result), nil

	case "json":
		return formatJSON(result)

	default:
		return nil, fmt.Errorf(
			"unsupported format %q; supported formats are txt, markdown, json",
			format,
		)
	}
}

func formatText(result scanner.ContextResult) []byte {
	var output strings.Builder

	if len(result.Tree) > 0 {
		fmt.Fprintln(
			&output,
			"PROJECT DIRECTORY STRUCTURE",
		)
		fmt.Fprintln(
			&output,
			"==========================",
		)
		fmt.Fprintln(&output)

		for _, line := range result.Tree {
			fmt.Fprintln(&output, line)
		}

		fmt.Fprintln(&output)
	}

	if len(result.Notes) > 0 {
		fmt.Fprintln(
			&output,
			"PROJECT CONTEXT NOTES",
		)
		fmt.Fprintln(
			&output,
			"=====================",
		)
		fmt.Fprintln(&output)

		for _, note := range result.Notes {
			fmt.Fprintf(
				&output,
				"- %s — %s\n",
				note.Path,
				note.Description,
			)
		}

		fmt.Fprintln(&output)
	}

	if len(result.Files) > 0 {
		fmt.Fprintln(
			&output,
			"FILE CONTENTS",
		)
		fmt.Fprintln(
			&output,
			"=============",
		)
		fmt.Fprintln(&output)

		for _, file := range result.Files {
			fmt.Fprintf(
				&output,
				"--- START FILE: %s ---\n",
				file.Path,
			)

			fmt.Fprintln(&output, file.Content)

			fmt.Fprintf(
				&output,
				"--- END FILE: %s ---\n\n",
				file.Path,
			)
		}
	}

	return []byte(output.String())
}

func formatMarkdown(result scanner.ContextResult) []byte {
	var output strings.Builder

	if len(result.Tree) > 0 {
		fmt.Fprintln(
			&output,
			"# Project Directory Structure",
		)
		fmt.Fprintln(&output)
		fmt.Fprintln(&output, "```text")

		for _, line := range result.Tree {
			fmt.Fprintln(&output, line)
		}

		fmt.Fprintln(&output, "```")
		fmt.Fprintln(&output)
	}

	if len(result.Notes) > 0 {
		fmt.Fprintln(
			&output,
			"# Project Context Notes",
		)
		fmt.Fprintln(&output)

		for _, note := range result.Notes {
			fmt.Fprintf(
				&output,
				"- `%s` — %s\n",
				note.Path,
				note.Description,
			)
		}

		fmt.Fprintln(&output)
	}

	if len(result.Files) > 0 {
		fmt.Fprintln(
			&output,
			"# File Contents",
		)
		fmt.Fprintln(&output)

		for _, file := range result.Files {
			fmt.Fprintf(
				&output,
				"## `%s`\n\n",
				file.Path,
			)

			fmt.Fprintln(&output, "```")
			fmt.Fprintln(&output, file.Content)
			fmt.Fprintln(&output, "```")
			fmt.Fprintln(&output)
		}
	}

	return []byte(output.String())
}

func formatJSON(result scanner.ContextResult) ([]byte, error) {
	type jsonResult struct {
		Tree         []string              `json:"tree,omitempty"`
		Notes        []scanner.ContextNote `json:"notes,omitempty"`
		Files        []scanner.ContextFile `json:"files,omitempty"`
		TotalBytes   int64                 `json:"total_bytes"`
		SkippedFiles int                   `json:"skipped_files"`
	}

	data := jsonResult{
		Tree:         result.Tree,
		Notes:        result.Notes,
		Files:        result.Files,
		TotalBytes:   result.TotalBytes,
		SkippedFiles: result.SkippedFiles,
	}

	return json.MarshalIndent(
		data,
		"",
		"  ",
	)
}
