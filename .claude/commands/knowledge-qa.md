---
description: Index PDF/Excel/CSV files in a directory and answer questions from their content
---

Arguments: $ARGUMENTS

Format: `<directory-or-file> "<question>"`

Follow these steps:

1. **Parse the arguments**
   - First token = path to a directory or a single document file
   - Remaining text (strip surrounding quotes) = question

2. **Expand the query**
   - Translate the question into both Japanese and English
   - Extract key terms from both translations (nouns and domain-specific words only, strip particles and stop words)
   - Combine all terms into a single space-separated search query

3. **Run the command**
   ```bash
   ~/.claude/commands/bin/knowledge-qa <path> "<expanded query>"
   ```
   - stderr may contain indexing progress messages (normal)
   - stdout returns a JSON array of relevant chunks

4. **Answer based on the JSON output**
   - Each element has `source` (filename), `location` (page/sheet/row), and `text` (content)
   - Answer the question in the same language the user used
   - End the response with a **References** section listing the source files and locations used

5. **Error cases**
   - Empty array `[]`: retry with the original question as-is, and if still empty inform the user that no relevant documents were found
   - `command not found`: instruct the user to run `make install` in the project directory
