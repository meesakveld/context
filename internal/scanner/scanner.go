package scanner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

type Options struct {
	IncludeEnv   bool
	NoIgnore     bool
	Excludes     []string
	Includes     []string
	MaxFileSize  int64
	TreeOnly     bool
	FilesOnly    bool
	IncludeNotes bool
}

type ContextResult struct {
	Tree         []string
	Files        []ContextFile
	Notes        []ContextNote
	TotalBytes   int64
	SkippedFiles int
}

type ContextFile struct {
	Path    string
	Content string
	Size    int64
}

type ContextNote struct {
	Path        string
	Description string
}

type treeEntry struct {
	Entry     os.DirEntry
	Collapsed bool
}

func Generate(options Options) (ContextResult, error) {
	result := ContextResult{}

	if !options.FilesOnly {
		tree, notes, err := generateTree(options)
		if err != nil {
			return result, err
		}

		result.Tree = tree
		result.Notes = notes
	}

	if !options.TreeOnly {
		files, notes, totalBytes, skippedFiles, err :=
			collectFiles(options)

		if err != nil {
			return result, err
		}

		result.Files = files
		result.TotalBytes = totalBytes
		result.SkippedFiles = skippedFiles

		if options.IncludeNotes {
			result.Notes = append(
				result.Notes,
				notes...,
			)
		}
	}

	return result, nil
}

func generateTree(
	options Options,
) ([]string, []ContextNote, error) {
	ignoreRules, err := loadIgnoreRules(
		options.NoIgnore,
	)
	if err != nil {
		return nil, nil, err
	}

	var output strings.Builder
	var notes []ContextNote

	err = printTreeRecursive(
		".",
		"",
		"",
		&output,
		options,
		ignoreRules,
		&notes,
	)
	if err != nil {
		return nil, nil, err
	}

	content := strings.TrimSuffix(
		output.String(),
		"\n",
	)

	if content == "" {
		return nil, notes, nil
	}

	return strings.Split(content, "\n"), notes, nil
}

func collectFiles(
	options Options,
) ([]ContextFile, []ContextNote, int64, int, error) {
	var files []ContextFile
	var notes []ContextNote
	var totalBytes int64
	var skippedFiles int

	ignoreRules, err := loadIgnoreRules(
		options.NoIgnore,
	)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	err = collectFilesRecursive(
		".",
		"",
		&files,
		&notes,
		&totalBytes,
		&skippedFiles,
		options,
		ignoreRules,
	)
	if err != nil {
		return nil, nil, 0, 0, err
	}

	sort.Slice(
		files,
		func(i, j int) bool {
			return strings.ToLower(
				files[i].Path,
			) < strings.ToLower(
				files[j].Path,
			)
		},
	)

	return files, notes, totalBytes, skippedFiles, nil
}

func collectFilesRecursive(
	path string,
	relativePath string,
	files *[]ContextFile,
	notes *[]ContextNote,
	totalBytes *int64,
	skippedFiles *int,
	options Options,
	ignoreRules *IgnoreRules,
) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(
			path,
			entry.Name(),
		)

		entryRelativePath := filepath.Join(
			relativePath,
			entry.Name(),
		)

		if shouldExclude(
			entry.Name(),
			entryRelativePath,
			entry.IsDir(),
			options,
			ignoreRules,
		) {
			continue
		}

		if entry.IsDir() {
			if err := collectFilesRecursive(
				entryPath,
				entryRelativePath,
				files,
				notes,
				totalBytes,
				skippedFiles,
				options,
				ignoreRules,
			); err != nil {
				return err
			}

			continue
		}

		if isEnvironmentFile(entry.Name()) &&
			!isSafeEnvironmentFile(entry.Name()) &&
			!options.IncludeEnv {
			if options.IncludeNotes {
				*notes = append(
					*notes,
					ContextNote{
						Path:        entryRelativePath,
						Description: "exists, but its contents are intentionally excluded because it may contain sensitive environment variables",
					},
				)
			}

			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if options.MaxFileSize > 0 &&
			info.Size() > options.MaxFileSize {
			*skippedFiles++

			if options.IncludeNotes {
				*notes = append(
					*notes,
					ContextNote{
						Path:        entryRelativePath,
						Description: "exists, but its contents are excluded because the file exceeds the configured maximum file size",
					},
				)
			}

			continue
		}

		if isBinaryFile(entryPath) ||
			!isTextFile(entryPath) {
			continue
		}

		data, err := os.ReadFile(entryPath)
		if err != nil {
			return err
		}

		normalizedPath := filepath.ToSlash(
			entryRelativePath,
		)

		*files = append(
			*files,
			ContextFile{
				Path:    normalizedPath,
				Content: string(data),
				Size:    info.Size(),
			},
		)

		*totalBytes += info.Size()
	}

	return nil
}

func printTreeRecursive(
	path string,
	relativePath string,
	prefix string,
	output io.Writer,
	options Options,
	ignoreRules *IgnoreRules,
	notes *[]ContextNote,
) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var filtered []treeEntry

	for _, entry := range entries {
		entryPath := filepath.Join(
			path,
			entry.Name(),
		)

		entryRelativePath := filepath.Join(
			relativePath,
			entry.Name(),
		)

		if isEnvironmentFile(entry.Name()) &&
			!isSafeEnvironmentFile(entry.Name()) &&
			!options.IncludeEnv {
			filtered = append(
				filtered,
				treeEntry{
					Entry: entry,
				},
			)

			if options.IncludeNotes {
				*notes = append(
					*notes,
					ContextNote{
						Path:        entryRelativePath,
						Description: "exists, but its contents are intentionally excluded because it may contain sensitive environment variables",
					},
				)
			}

			continue
		}

		if isExcludedDirectory(
			entry.Name(),
			entryRelativePath,
			options,
			ignoreRules,
		) {
			if entry.IsDir() {
				filtered = append(
					filtered,
					treeEntry{
						Entry:     entry,
						Collapsed: true,
					},
				)

				if options.IncludeNotes {
					*notes = append(
						*notes,
						ContextNote{
							Path:        entryRelativePath + "/",
							Description: exclusionDescription(),
						},
					)
				}
			}

			continue
		}

		if shouldExclude(
			entry.Name(),
			entryRelativePath,
			entry.IsDir(),
			options,
			ignoreRules,
		) {
			continue
		}

		if !entry.IsDir() &&
			(isBinaryFile(entryPath) ||
				!isTextFile(entryPath)) {
			continue
		}

		filtered = append(
			filtered,
			treeEntry{
				Entry: entry,
			},
		)
	}

	sortTreeEntries(filtered)

	for i, item := range filtered {
		entry := item.Entry
		isLast := i == len(filtered)-1

		branch := "├── "
		nextPrefix := prefix + "│   "

		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		name := entry.Name()

		if item.Collapsed {
			name += "/..."
		}

		fmt.Fprintln(
			output,
			prefix+branch+name,
		)

		if entry.IsDir() && !item.Collapsed {
			if err := printTreeRecursive(
				filepath.Join(path, entry.Name()),
				filepath.Join(relativePath, entry.Name()),
				nextPrefix,
				output,
				options,
				ignoreRules,
				notes,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func sortTreeEntries(entries []treeEntry) {
	sort.SliceStable(
		entries,
		func(i, j int) bool {
			left := entries[i].Entry
			right := entries[j].Entry

			if left.IsDir() != right.IsDir() {
				return left.IsDir()
			}

			return strings.ToLower(
				left.Name(),
			) < strings.ToLower(
				right.Name(),
			)
		},
	)
}

func isExcludedDirectory(
	name string,
	relativePath string,
	options Options,
	ignoreRules *IgnoreRules,
) bool {
	if matchesAnyPattern(
		name,
		relativePath,
		options.Includes,
	) {
		return false
	}

	if matchesAnyPattern(
		name,
		relativePath,
		options.Excludes,
	) {
		return true
	}

	if ignoreRules.Matches(
		relativePath,
		true,
	) {
		return true
	}

	return excludedDirectories[name]
}

func shouldExclude(
	name string,
	relativePath string,
	isDirectory bool,
	options Options,
	ignoreRules *IgnoreRules,
) bool {
	if isEnvironmentFile(name) {
		return !options.IncludeEnv &&
			!isSafeEnvironmentFile(name)
	}

	if matchesAnyPattern(
		name,
		relativePath,
		options.Includes,
	) {
		return false
	}

	if matchesAnyPattern(
		name,
		relativePath,
		options.Excludes,
	) {
		return true
	}

	if ignoreRules.Matches(
		relativePath,
		isDirectory,
	) {
		return true
	}

	if isDirectory {
		return excludedDirectories[name]
	}

	if excludedFiles[name] {
		return true
	}

	if isBinaryFile(relativePath) {
		return true
	}

	return false
}

func matchesAnyPattern(
	name string,
	path string,
	patterns []string,
) bool {
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(
			strings.TrimSpace(pattern),
		)

		if pattern == "" {
			continue
		}

		if matched, _ := filepath.Match(
			pattern,
			name,
		); matched {
			return true
		}

		if matched, _ := filepath.Match(
			pattern,
			filepath.ToSlash(path),
		); matched {
			return true
		}

		cleanPattern := strings.TrimSuffix(
			pattern,
			"/",
		)

		if strings.HasPrefix(
			filepath.ToSlash(path),
			cleanPattern+"/",
		) {
			return true
		}
	}

	return false
}

func exclusionDescription() string {
	return "exists, but its contents are excluded from the project context"
}

func isTextFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return utf8.Valid(data)
}
