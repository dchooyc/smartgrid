# AI usage journal

## Tool(s) and model(s) used.
- Claude Code extension in vscode, Fable 5.1
- ChatGPT to validate ideas, GPT-6

## Work decomposition and sequencing

1. Frontend component to upload .csv files
2. Backend to receive .csv from frontend
3. Backend to validate data according to listed rules
4. Frontend to receive validation results from backend
5. Frontend to display results and option to commit valid rows
6. Commit valid rows functionality writes to DB
7. Frontend explorer component to display data from backend API
8. Frontend search functionality
9. Review unit tests for import and hierachy validation
10. Containerize full app
11. Add logging for errors
12. Review and refactor code for clarity
13. Write all documentation
14. Host app

## Three representative prompts
```
# Frontend component to upload csv

Build a frontend file upload component on the Next.js main page.

Requirements:
1. Accept only `.csv` files.
2. Store the selected file in React state without processing or uploading it.
3. Display the selected filename.
4. Add a button labeled **Validate** directly below the file input.
5. Disable the Validate button when no valid file is selected.
6. Keep the component simple, readable, and reusable, using TypeScript and Tailwind CSS.

Do not implement validation logic, a backend, or any additional features. Do not modify unrelated files.
```

```
# Backend to receive .csv from frontend

Create a minimal Go backend (as a folder named backend outside of frontend) and connect it to the existing Next.js CSV upload component.

### Backend

1. Use Go's standard library (`net/http` and `encoding/csv`). Keep the code extremely simple, clean, and readable.
2. Create a `POST /api/validate` endpoint that receives the CSV file using `multipart/form-data`.
3. Read the CSV in memory without saving it.
4. Verify that the file contains exactly these columns, with no missing or additional columns:
   - `asset_id`
   - `parent_asset_id`
   - `asset_type`
   - `asset_name`
   - `operational_status`
   - `voltage_kv`
   - `rating_kva`
   - `manufacturer`
   - `model`
   - `serial_number`
   - `commissioned_date`
5. Count the number of data rows, excluding the header.
6. Return a JSON response containing the row count and whether the columns are valid. Return a meaningful error if parsing or column validation fails.

### Frontend

1. When **Validate** is clicked, send the selected CSV file to the Go backend.
2. While waiting for the backend response, display directly below the button: "Doing further validations on the backend..."
3. On success, replace the loading message with: "File has been validated successfully! X rows read!"
4. On failure, display the error returned by the backend.
5. Prevent repeated submissions while loading.

### Constraints

- No database, Docker, or file persistence.
- No asset-level or hierarchy validation yet.
- Avoid unnecessary dependencies, abstractions, and files.
- Keep the existing frontend design unchanged apart from the requested functionality.
- Provide simple instructions to run the Go backend locally.

```

## Rejected suggestions and why
1. In the csv upload component it's check that the file is a csv only checked that it ends with .csv. I prompted for a better check as files can be mistakenly named. While these "simpler" validations can be made on the backend, my intention is to catch errors early and reduce unneccessary calls to the backend. I also prompted for missing/ duplicate columns checks.

## Reviewing and verification of generated code

## What the agent did poorly

## Approximate human working time

| Day | Work done | Human working time |
|-----|-------------------------|--------------------|
| Day 0 | Read requirements, planning, setup | 2 hours |
| Day 1 | TBD | TBD |
| Day 2 | TBD | TBD |
| Day 3 |TBD | TBD |
| Day 4 | TBD | TBD |
| Day 5 | TBD | TBD |
| **Total** | | **TBD** |