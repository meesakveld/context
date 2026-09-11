package initializer

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const defaultConfigTemplate = `# Context configuration

output: context.txt

format: txt

max_file_size: 1MB

include_env: false

tree_only: false
files_only: false

include_notes: true

exclude: []

include: []
`

const defaultIgnoreTemplate = `# Dependencies
node_modules/
venv/
.venv/

# Build output
.next/
dist/
build/

# Python cache
__pycache__/

# IDE
.idea/
.vscode/

# Generated files
coverage/
.cache/
.tmp/

# Secrets
*.pem
*.key
*.p12
*.pfx
`

func Initialize() error {
	if _, err := os.Stat(".context.yaml"); err == nil {
		return errors.New(".context.yaml already exists")
	}

	if _, err := os.Stat(".contextignore"); err == nil {
		return errors.New(".contextignore already exists")
	}

	if err := os.WriteFile(
		".context.yaml",
		[]byte(defaultConfigTemplate),
		0644,
	); err != nil {
		return err
	}

	if err := os.WriteFile(
		".contextignore",
		[]byte(defaultIgnoreTemplate),
		0644,
	); err != nil {
		return err
	}

	fmt.Println("✓ Created .context.yaml")
	fmt.Println("✓ Created .contextignore")

	updated, err := updateGitignore()
	if err != nil {
		return err
	}

	if updated {
		fmt.Println("✓ Updated .gitignore")
	}

	return nil
}

func updateGitignore() (bool, error) {
	const entry = "# Context output\ncontext.txt\n"

	data, err := os.ReadFile(".gitignore")

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	content := string(data)

	if strings.Contains(content, "context.txt") {
		return false, nil
	}

	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	if content != "" {
		content += "\n"
	}

	content += entry

	if err := os.WriteFile(
		".gitignore",
		[]byte(content),
		0644,
	); err != nil {
		return false, err
	}

	return true, nil
}
