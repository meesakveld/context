package stats

import (
	"fmt"
	"io"

	"github.com/meesakveld/context/internal/scanner"
)

func Print(
	writer io.Writer,
	result scanner.ContextResult,
) {
	var characters int

	for _, file := range result.Files {
		characters += len([]rune(file.Content))
	}

	estimatedTokens := characters / 4

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Context statistics:")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "  Files:")
	fmt.Fprintf(
		writer,
		"    Included:       %d\n",
		len(result.Files),
	)
	fmt.Fprintf(
		writer,
		"    Skipped:        %d\n",
		result.SkippedFiles,
	)
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "  Size:")
	fmt.Fprintf(
		writer,
		"    Source files:   %s\n",
		formatByteCount(result.TotalBytes),
	)
	fmt.Fprintf(
		writer,
		"    Context:        %s\n",
		formatByteCount(int64(characters)),
	)
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "  AI estimate:")
	fmt.Fprintf(
		writer,
		"    Characters:     %d\n",
		characters,
	)
	fmt.Fprintf(
		writer,
		"    Tokens:         ~%d\n",
		estimatedTokens,
	)
}

func formatByteCount(value int64) string {
	const (
		kilobyte = 1024
		megabyte = 1024 * 1024
		gigabyte = 1024 * 1024 * 1024
	)

	switch {
	case value >= gigabyte:
		return fmt.Sprintf(
			"%.2f GB",
			float64(value)/gigabyte,
		)

	case value >= megabyte:
		return fmt.Sprintf(
			"%.2f MB",
			float64(value)/megabyte,
		)

	case value >= kilobyte:
		return fmt.Sprintf(
			"%.2f KB",
			float64(value)/kilobyte,
		)

	default:
		return fmt.Sprintf("%d B", value)
	}
}
