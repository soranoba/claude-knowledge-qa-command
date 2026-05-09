---
description: Index PDF/Excel/CSV files in a directory and answer questions from their content
---

Arguments: $ARGUMENTS

Format: `<directory-or-file> "<question>"`

Follow these steps:

1. **Parse the arguments**
   - First token = path to a directory or a single document file
   - Remaining text (strip surrounding quotes) = question

2. **Run the command**
   ```bash
   ~/.claude/commands/bin/knowledge-qa <path> <question>
   ```
   - stderr may contain indexing progress messages (normal)
   - stdout returns a JSON array of relevant chunks

3. **Answer based on the JSON output**
   - Each element has `source` (filename), `location` (page/sheet/row), and `text` (content)
   - Answer the question in the same language the user used
   - End the response with a **References** section listing the source files and locations used

4. **Error cases**
   - Empty array `[]`: inform the user that no relevant documents were found
   - `command not found`: instruct the user to run `make install` in the project directory
