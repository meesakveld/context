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

	outputPath := flag.String(
		"o",
		cfg.Output,
		"output file path",
	)

	clipboard := flag.Bool(
		"clipboard",
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

	format := flag.String(
		"format",
		cfg.Format,
		"output format: txt, markdown, json",
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

	showVersion := flag.Bool(
		"version",
		false,
		"show version",
	)

	flag.Parse()

	if *showVersion {
		fmt.Println(version.String())
		return
	}

	if *initProject {
		if err := initializer.Initialize(); err != nil {
			fmt.Println("Error initializing project:", err)
			os.Exit(1)
		}

		return
	}

	if *treeOnly && *filesOnly {
		fmt.Println("Error: --tree-only and --files-only cannot be used together")
		os.Exit(1)
	}

	if *stdout && *clipboard {
		fmt.Println("Error: --stdout and --clipboard cannot be used together")
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

	if !*stdout && !*clipboard && progress.IsTerminal() {
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

	data, err := formatter.Format(result, *format)
	if err != nil {
		fmt.Println("Error formatting context:", err)
		os.Exit(1)
	}

	switch {
	case *stdout:
		fmt.Print(string(data))

	case *clipboard:
		if err := output.CopyToClipboard(data); err != nil {
			fmt.Println("Error copying to clipboard:", err)
			os.Exit(1)
		}

	default:
		if err := output.WriteFile(*outputPath, data); err != nil {
			fmt.Println("Error writing output file:", err)
			os.Exit(1)
		}
	}

	if *showStats {
		stats.Print(os.Stderr, result)
	}

	if !*stdout && !*clipboard {
		fmt.Printf("✓ %d files collected\n", len(result.Files))
		fmt.Printf("✓ Context written to %s\n", *outputPath)
	}

	if *clipboard {
		fmt.Println("✓ Context copied to clipboard")
	}
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
