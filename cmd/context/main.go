package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/meesakveld/context/internal/config"
	"github.com/meesakveld/context/internal/formatter"
	"github.com/meesakveld/context/internal/initializer"
	"github.com/meesakveld/context/internal/output"
	"github.com/meesakveld/context/internal/progress"
	"github.com/meesakveld/context/internal/scanner"
	"github.com/meesakveld/context/internal/stats"
	"github.com/meesakveld/context/internal/version"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	flag.Usage = printHelp

	outputPath := flag.String(
		"o",
		"",
		"output file path",
	)

	outputPathLong := flag.String(
		"output",
		"",
		"output file path",
	)

	clipboard := flag.Bool(
		"clipboard",
		false,
		"copy context to clipboard",
	)

	clipboardShort := flag.Bool(
		"c",
		false,
		"copy context to clipboard",
	)

	stdout := flag.Bool(
		"stdout",
		false,
		"write context to stdout",
	)

	includeEnv := flag.Bool(
		"include-env",
		cfg.IncludeEnv,
		"include environment file contents",
	)

	treeOnly := flag.Bool(
		"tree-only",
		cfg.TreeOnly,
		"only generate the directory tree",
	)

	filesOnly := flag.Bool(
		"files-only",
		cfg.FilesOnly,
		"only generate file contents",
	)

	showStats := flag.Bool(
		"stats",
		false,
		"show context statistics",
	)

	noNotes := flag.Bool(
		"no-notes",
		!cfg.IncludeNotes,
		"exclude project context notes",
	)

	formatShort := flag.String(
		"f",
		cfg.Format,
		"output format: txt, markdown, md, json, zip",
	)

	formatLong := flag.String(
		"format",
		cfg.Format,
		"output format: txt, markdown, md, json, zip",
	)

	maxFileSize := flag.String(
		"max-file-size",
		cfg.MaxFileSize,
		"maximum file size to include",
	)

	exclude := flag.String(
		"exclude",
		strings.Join(cfg.Exclude, ","),
		"additional files or directories to exclude",
	)

	include := flag.String(
		"include",
		strings.Join(cfg.Include, ","),
		"files or directories to explicitly include",
	)

	noIgnore := flag.Bool(
		"no-ignore",
		false,
		"ignore the .contextignore file",
	)

	initProject := flag.Bool(
		"init",
		false,
		"create context configuration files",
	)

	initShort := flag.Bool(
		"i",
		false,
		"create context configuration files",
	)

	showVersion := flag.Bool(
		"version",
		false,
		"show version",
	)

	versionShort := flag.Bool(
		"v",
		false,
		"show version",
	)

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	shouldClipboard := *clipboard || *clipboardShort
	shouldInit := *initProject || *initShort
	shouldShowVersion := *showVersion || *versionShort

	if shouldShowVersion {
		fmt.Println(version.String())
		return
	}

	if shouldInit {
		if err := initializer.Initialize(); err != nil {
			fmt.Println("Error initializing project:", err)
			os.Exit(1)
		}
		return
	}

	formatValue := cfg.Format
	if *formatShort != cfg.Format {
		formatValue = *formatShort
	}
	if *formatLong != cfg.Format {
		formatValue = *formatLong
	}

	// Bepaal output pad: expliciet opgegeven via -o/--output, anders automatisch op basis van format
	outputPathValue := cfg.Output
	if *outputPath != "" {
		outputPathValue = *outputPath
	} else if *outputPathLong != "" {
		outputPathValue = *outputPathLong
	} else {
		if cfg.Output == "context.txt" || cfg.Output == "" {
			outputPathValue = config.DefaultOutputForFormat(formatValue)
		}
	}

	if *treeOnly && *filesOnly {
		fmt.Println("Error: --tree-only and --files-only cannot be used together")
		os.Exit(1)
	}

	if *stdout && shouldClipboard {
		fmt.Println("Error: --stdout and --clipboard cannot be used together")
		os.Exit(1)
	}

	if strings.ToLower(formatValue) == "zip" && (*stdout || shouldClipboard) {
		fmt.Println("Error: --stdout and --clipboard are not supported for zip format")
		os.Exit(1)
	}

	size, err := config.ParseBytes(*maxFileSize)
	if err != nil {
		fmt.Println("Error parsing --max-file-size:", err)
		os.Exit(1)
	}

	options := scanner.Options{
		IncludeEnv:   *includeEnv,
		NoIgnore:     *noIgnore,
		Excludes:     splitPatterns(*exclude),
		Includes:     splitPatterns(*include),
		MaxFileSize:  size,
		TreeOnly:     *treeOnly,
		FilesOnly:    *filesOnly,
		IncludeNotes: !*noNotes,
	}

	var spinner *progress.Spinner

	if !*stdout && !shouldClipboard && progress.IsTerminal() {
		spinner = progress.New(os.Stdout)
		spinner.Start("Scanning files")
	}

	result, err := scanner.Generate(options)

	if spinner != nil {
		spinner.Stop()
	}

	if err != nil {
		fmt.Println("Error generating context:", err)
		os.Exit(1)
	}

	if strings.ToLower(formatValue) == "zip" {
		if err := output.WriteZip(outputPathValue, result); err != nil {
			fmt.Println("Error writing zip archive:", err)
			os.Exit(1)
		}
	} else {
		data, err := formatter.Format(result, formatValue)
		if err != nil {
			fmt.Println("Error formatting context:", err)
			os.Exit(1)
		}

		switch {
		case *stdout:
			fmt.Print(string(data))

		case shouldClipboard:
			if err := output.CopyToClipboard(data); err != nil {
				fmt.Println("Error copying context to clipboard:", err)
				os.Exit(1)
			}

		default:
			if err := output.WriteFile(outputPathValue, data); err != nil {
				fmt.Println("Error writing output file:", err)
				os.Exit(1)
			}
		}
	}

	if *showStats {
		stats.Print(os.Stderr, result)
	}

	if !*stdout && !shouldClipboard {
		fmt.Printf("✓ %d files collected\n", len(result.Files))
		if strings.ToLower(formatValue) == "zip" {
			fmt.Printf("✓ Context archive written to %s\n", outputPathValue)
		} else {
			fmt.Printf("✓ Context written to %s\n", outputPathValue)
		}
	}

	if shouldClipboard {
		fmt.Println("✓ Context copied to clipboard")
	}
}

func printHelp() {
	fmt.Println(`context — Turn your codebase into AI-ready context.

Usage:
  context [options]

Options:
  -h, --help               Show this help message
  -v, --version            Show version
  -o, --output <file>      Output file path
  -c, --clipboard          Copy context to clipboard
      --stdout             Write context to stdout
  -i, --init               Create context configuration files
      --tree-only          Only generate the directory tree
      --files-only         Only generate file contents
      --stats              Show context statistics
  -f, --format <format>    Output format: txt, markdown, md, json, zip
      --max-file-size <s>  Maximum file size to include
      --include-env        Include environment file contents
      --exclude <patterns> Additional files or directories to exclude
      --include <patterns> Files or directories to explicitly include
      --no-ignore          Ignore the .contextignore file
      --no-notes           Exclude project context notes

Examples:
  context
  context -c
  context -o project.txt
  context -f markdown
  context -f json
  context -f zip
  context --stdout
  context --tree-only
  context -i`)
}

func splitPatterns(value string) []string {
	if value == "" {
		return nil
	}

	var patterns []string

	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)

		if part != "" {
			patterns = append(patterns, part)
		}
	}

	return patterns
}