# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ni-idea** is a personal knowledge base CLI tool. It stores knowledge locally in markdown format and provides fast search via the `ni` command. Primary use case: Claude Code calling `ni` via bash to retrieve relevant context.

## Tech Stack

| Item        | Choice                                    |
| ----------- | ----------------------------------------- |
| Language    | Go (single binary, cross-platform)        |
| Note Format | Markdown + YAML frontmatter               |
| Search      | bleve (full-text search with fuzzy match) |
| Config      | `~/.ni-idea/config.yaml`                  |
| Index       | `~/.cache/ni-idea/index` (bleve)          |

## Project Structure

```
ni-idea/
├── cmd/ni/main.go           # CLI entrypoint
├── internal/
│   ├── cmd/                 # CLI commands
│   ├── index/               # bleve index wrapper
│   ├── store/               # Note read/write
│   ├── remote/              # API client for server
│   ├── formatter/           # stdout output formatting
│   └── config/              # Config loading
├── server/main.go           # REST API server
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
# Search
ni search "query"              # Search problems/decisions (default types)
ni search "query" --all        # Include knowledge/practice types
ni search "query" --tag infra  # Filter by tag
ni search "query" --fuzzy      # Fuzzy search (typo tolerance)

# Notes
ni get problems/nextjs-caching # Get full note content
ni list                        # List notes
ni list --tag infra
ni add                         # Add note (interactive)
ni add --title "Title" --type problem --tag tag1,tag2
ni tags                        # List all tags with counts

# Index
ni index status                # Check index status
ni index rebuild               # Rebuild search index

# Remote sync
ni remote add <name> <url>     # Add remote server
ni remote list                 # List remotes
ni remote remove <name>        # Remove remote
ni push --all                  # Push all notes
ni push problems/my-note.md    # Push specific note
ni pull                        # Pull from remote
ni pull --theirs               # Use remote version on conflict
```

## Note Types

| Type        | Purpose                     | Default Search |
| ----------- | --------------------------- | -------------- |
| `problem`   | Problem resolution records  | Yes            |
| `decision`  | Architecture/tech decisions | Yes            |
| `knowledge` | Concept/tech documentation  | No (`--all`)   |
| `practice`  | Practice notes              | No (`--all`)   |

## Architecture Notes

- Notes stored in `~/.ni-idea/notes/` with YAML frontmatter + markdown body
- Search uses bleve index stored in `~/.cache/ni-idea/index`
- `private: true` notes excluded from default search/list and push
- Output designed for Claude Code consumption: plain markdown to stdout, errors to stderr
- `--json` flag available for structured output
- Optional server (`server/main.go`) for remote sync across devices

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
