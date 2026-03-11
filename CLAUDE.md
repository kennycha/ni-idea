# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ni-idea** is a personal knowledge base CLI tool. It stores knowledge locally in markdown format and provides fast search via the `ni` command. Primary use case: Claude Code calling `ni` via bash to retrieve relevant context.

## Tech Stack

| Item        | Choice                                     |
| ----------- | ------------------------------------------ |
| Language    | Go (single binary, cross-platform)         |
| Note Format | Markdown + YAML frontmatter                |
| Search      | File system grep (initial), bleve (future) |
| Config      | `~/.config/ni-idea/config.yaml`            |

## Project Structure

```
ni-idea/
├── cmd/ni/main.go           # CLI entrypoint
├── internal/
│   ├── search/              # Search logic
│   ├── store/               # Note read/write
│   ├── formatter/           # stdout output formatting
│   └── config/              # Config loading
├── go.mod
├── go.sum
└── Makefile
```

## Commands

```bash
# Build
go build -o ni ./cmd/ni

# Run tests
go test ./...

# Run single test
go test ./internal/search -run TestSearchByTag

# Install locally
go install ./cmd/ni
```

## CLI Usage

```bash
ni search "query"              # Search problems/decisions (default types)
ni search "query" --all        # Include knowledge/practice types
ni search "query" --tag infra  # Filter by tag
ni search "query" --domain company-a
ni get problems/frontend/nextjs-caching   # Get full note content
ni list                        # List notes
ni list --tag infra
ni add                         # Add note (interactive)
ni add --title "Title" --type problem --tag tag1,tag2
ni tags                        # List all tags with counts
```

## Note Types

| Type        | Purpose                     | Default Search |
| ----------- | --------------------------- | -------------- |
| `problem`   | Problem resolution records  | Yes            |
| `decision`  | Architecture/tech decisions | Yes            |
| `knowledge` | Concept/tech documentation  | No (`--all`)   |
| `practice`  | Practice notes              | No (`--all`)   |

## Architecture Notes

- Notes stored in `~/notes/` (configurable) with YAML frontmatter + markdown body
- `private: true` notes excluded from default search/list
- Output designed for Claude Code consumption: plain markdown to stdout, errors to stderr
- `--json` flag available for structured output

## Note Frontmatter Schema

```yaml
title: ""
type: problem | decision | knowledge | practice
tags: []
domain: general | company-a | ...
private: false
created: YYYY-MM-DD
updated: YYYY-MM-DD
```

See `docs/project-spec.md` for detailed note templates and full specification.
