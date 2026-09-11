package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type IgnoreRules struct {
	Patterns []string
}

func loadIgnoreRules(noIgnore bool) (*IgnoreRules, error) {
	if noIgnore {
		return &IgnoreRules{}, nil
	}

	data, err := os.ReadFile(".contextignore")

	if os.IsNotExist(err) {
		return &IgnoreRules{}, nil
	}

	if err != nil {
		return nil, err
	}

	var patterns []string

	for _, line := range strings.Split(
		string(data),
		"\n",
	) {
		line = strings.TrimSpace(line)

		if line == "" ||
			strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}

	return &IgnoreRules{
		Patterns: patterns,
	}, nil
}

func (rules *IgnoreRules) Matches(
	path string,
	isDirectory bool,
) bool {
	if rules == nil {
		return false
	}

	path = filepath.ToSlash(path)

	matched := false

	for _, pattern := range rules.Patterns {
		negated := strings.HasPrefix(
			pattern,
			"!",
		)

		if negated {
			pattern = strings.TrimPrefix(
				pattern,
				"!",
			)
		}

		pattern = strings.TrimSpace(pattern)
		pattern = strings.TrimSuffix(pattern, "/")

		if pattern == "" {
			continue
		}

		if matchesPattern(
			path,
			pattern,
			isDirectory,
		) {
			matched = !negated
		}
	}

	return matched
}

func matchesPattern(
	path string,
	pattern string,
	isDirectory bool,
) bool {
	path = filepath.ToSlash(path)
	pattern = filepath.ToSlash(pattern)

	pattern = strings.TrimPrefix(
		pattern,
		"/",
	)

	if strings.Contains(pattern, "/") {
		if matched, _ := filepath.Match(
			pattern,
			path,
		); matched {
			return true
		}

		return strings.HasPrefix(
			path,
			pattern+"/",
		)
	}

	name := filepath.Base(path)

	if matched, _ := filepath.Match(
		pattern,
		name,
	); matched {
		return true
	}

	if isDirectory {
		return name == pattern
	}

	return false
}
