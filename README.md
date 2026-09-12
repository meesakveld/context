<div align="center">

# context

**Turn your codebase into AI-ready context.**

A fast, lightweight CLI that turns your codebase into structured context for AI tools.

<br />

[Installation](#installation) · [Usage](#usage) · [Configuration](#configuration)

<br />

[![GitHub Stars](https://img.shields.io/github/stars/meesakveld/context?style=flat-square)](https://github.com/meesakveld/context/stargazers)
[![GitHub Downloads](https://img.shields.io/github/downloads/meesakveld/context/total?style=flat-square)](https://github.com/meesakveld/context/releases)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?style=flat-square\&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache--2.0-black?style=flat-square)](LICENSE)

</div>

---

## What is context?

`context` turns your project directory into a structured, AI-friendly representation.

Instead of manually copying files and explaining your project structure, run:

```bash
context
```

and get a single context file containing:

* Project directory structure
* Relevant source files
* File contents
* Project context notes
* Configurable exclusions
* File size limits

It is designed to make working with AI coding tools faster and more consistent.

---

## Quick start

Install `context`:

```bash
brew install meesakveld/tap/context
```

Navigate to your project:

```bash
cd your-project
```

Generate your project context and copy it directly to your clipboard:

```bash
context --clipboard
```

Then paste it into your AI assistant.

That's it.

---

## Why?

AI coding tools work much better when they understand the structure and relevant parts of a project.

But getting an entire codebase into an AI assistant isn't always convenient. Depending on the tool, uploading a project archive may not be supported, practical, or the best way to provide context.

`context` turns your codebase into a single, structured text representation that you can paste directly into AI tools such as **Claude, ChatGPT, Gemini, and other AI coding assistants**.

```bash
context --clipboard
```

Then paste the generated context directly into your AI assistant.

No manual file selection. No copying files one by one. No need to manually create or upload a project archive.

### Built for AI agents

`context` can also act as a simple bridge between your local codebase and AI agents that work primarily with text input.

```text
Your project
     │
     ▼
   context
     │
     ├── Remove irrelevant files
     ├── Apply .contextignore
     ├── Skip binaries
     ├── Respect file size limits
     └── Handle environment files safely
     │
     ▼
Structured project context
     │
     ├── ChatGPT
     ├── Claude
     ├── Gemini
     └── Other AI agents
```

The goal is simple:

**Less context preparation. More building.**

---

## Features

* 🚀 Fast project scanning
* 📁 Automatic directory tree generation
* 📄 Collect relevant text files
* 🚫 Built-in exclusions for dependencies and build output
* 🔍 `.contextignore` support
* ⚙️ YAML configuration
* 📋 Clipboard support
* 🖥️ Stdout support
* 📊 Context statistics
* 🌳 Tree-only mode
* 📄 Files-only mode
* 🔐 Safe handling of environment files
* 🎨 Text, Markdown and JSON output
* 💻 Cross-platform
* 🪶 Single binary with no runtime required

---

## Installation

### Homebrew

```bash
brew install meesakveld/tap/context
```

### From source

Make sure you have Go installed.

```bash
go install github.com/meesakveld/context/cmd/context@latest
```

### Build locally

```bash
git clone https://github.com/meesakveld/context.git
cd context

go build -o context ./cmd/context
```

---

## Usage

Run `context` from the root of your project:

```bash
context
```

This generates:

```text
context.txt
```

with your project structure and relevant file contents.

### Clipboard

Copy the generated context directly to your clipboard:

```bash
context --clipboard
```

This is useful when you want to paste your project context directly into an AI assistant.

### Standard output

Print the generated context directly to stdout:

```bash
context --stdout
```

This is useful for piping the output into other tools.

### Tree only

Only generate the project directory structure:

```bash
context --tree-only
```

### Files only

Only include file contents:

```bash
context --files-only
```

### Statistics

Show context statistics:

```bash
context --stats
```

### Custom output

Write the generated context to a custom file:

```bash
context -o project-context.txt
```

### Output formats

Generate plain text:

```bash
context --format txt
```

Generate Markdown:

```bash
context --format markdown
```

Generate JSON:

```bash
context --format json
```

### File size limit

Limit the maximum size of individual files:

```bash
context --max-file-size 500KB
```

---

## Configuration

Create the default configuration files with:

```bash
context --init
```

This creates:

```text
.context.yaml
.contextignore
```

If a `.gitignore` already exists, `context --init` also adds `context.txt` to it.

### `.context.yaml`

Example:

```yaml
output: context.txt

format: txt

max_file_size: 1MB

include_env: false

tree_only: false

files_only: false

include_notes: true

exclude: []

include: []
```

### `.contextignore`

Use `.contextignore` to exclude files or directories:

```gitignore
# Dependencies
node_modules/
venv/
.venv/

# Build output
.next/
dist/
build/

# Generated files
coverage/
.cache/
.tmp/

# Secrets
*.pem
*.key
*.p12
*.pfx
```

---

## Environment files

Environment files are handled carefully by default.

These example files are safe to include:

```text
.env.example
.env.template
.env.sample
```

Actual environment files such as:

```text
.env
.env.local
.env.production
```

are not included unless explicitly enabled:

```bash
context --include-env
```

This helps prevent accidentally exposing secrets when generating context for AI tools.

---

## Examples

### Typical project

Running:

```bash
context
```

might produce:

```text
PROJECT DIRECTORY STRUCTURE

├── app/
│   ├── components/
│   │   ├── Button.tsx
│   │   └── Header.tsx
│   ├── page.tsx
│   └── layout.tsx
├── package.json
├── README.md
└── tsconfig.json

FILE CONTENTS

--- START FILE: app/components/Button.tsx ---

export function Button() {
    ...
}

--- END FILE: app/components/Button.tsx ---
```

---

## How it works

`context` scans your project and builds a structured representation while filtering out files that are unlikely to be useful to an AI assistant.

```text
Project
   │
   ▼
Scan files
   │
   ├── Apply exclusions
   ├── Apply .contextignore
   ├── Check file types
   ├── Check file size
   └── Handle environment files
   │
   ▼
Build project tree
   │
   ▼
Collect file contents
   │
   ▼
Format output
   │
   ▼
context.txt
```

The generated context contains both the structure of your project and the contents of relevant files, making it easier for AI tools to understand how your project is organised.

---

## Command reference

| Command                     | Description                      |
| --------------------------- | -------------------------------- |
| `context`                   | Generate project context         |
| `context -o file.txt`       | Write to a custom file           |
| `context --clipboard`       | Copy context to clipboard        |
| `context --stdout`          | Print context to stdout          |
| `context --tree-only`       | Generate only the directory tree |
| `context --files-only`      | Generate only file contents      |
| `context --stats`           | Show context statistics          |
| `context --format txt`      | Generate plain text              |
| `context --format markdown` | Generate Markdown                |
| `context --format json`     | Generate JSON                    |
| `context --include-env`     | Include environment files        |
| `context --init`            | Create project configuration     |
| `context --version`         | Show version                     |

---

## Development

Clone the repository:

```bash
git clone https://github.com/meesakveld/context.git
cd context
```

Run the CLI:

```bash
go run ./cmd/context
```

Build:

```bash
go build ./...
```

Format the code:

```bash
gofmt -w cmd internal
```

---

## Contributing

Contributions, ideas and feedback are welcome.

If you find a bug or have an idea for a feature, open an issue or submit a pull request.

---

## License

Apache 2.0 © Mees Akveld
