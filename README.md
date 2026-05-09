# claude-knowledge-qa-command

A Claude Code custom command that indexes PDF, Excel, and CSV files in a directory and answers questions based on their content.

## Overview

Type `/knowledge-qa <directory-or-file> "<question>"` to search across your local documents and get an answer.

- **Supported formats**: PDF / Excel (.xlsx, .xls) / CSV
- **Fully local**: No external API required. Files are indexed locally and relevant chunks are retrieved via bigram search.
- **Bilingual**: The command expands queries into both Japanese and English internally, so questions in either language work.

## Installation

```bash
git clone https://github.com/soranoba/claude-knowledge-qa-command.git
cd claude-knowledge-qa-command
make install
```

`make install` does the following:

- Builds the binary and places it at `~/.claude/commands/bin/knowledge-qa`
- Installs the command definition at `~/.claude/commands/knowledge-qa.md`

## allowList Configuration (Recommended)

Without an allowList entry, Claude Code will prompt for permission every time the command runs the binary. Add the following to `~/.claude/settings.json` to allow it to run without prompts:

```json
{
  "permissions": {
    "allow": [
      "Bash(~/.claude/commands/bin/knowledge-qa *)"
    ]
  }
}
```

## Usage

In the Claude Code chat, type:

```
/knowledge-qa <directory-or-file-path> "<question>"
```

**Examples:**

```
/knowledge-qa ~/Documents/reports "What was the revenue in FY2024?"
/knowledge-qa ~/Documents/manual.pdf "How do I fix error code E001?"
/knowledge-qa ./data "What is the average age of users?"
```

### How it works

1. Scans the specified path for PDF, Excel, and CSV files. Only files that have changed are re-indexed (indexes are stored in a `.index/` directory alongside each file).
2. Expands the question into Japanese and English, then scores chunks using bigram matching.
3. Claude generates an answer from the top-ranked chunks and cites the source files and locations.
