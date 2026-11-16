# AI usage journal

## Tool(s) and model(s) used
- Claude Code extension in VS Code, Fable 5.2
- ChatGPT to validate ideas, GPT-6

## Work decomposition and sequencing

1. Frontend component to upload .csv files ✅
2. Backend to receive .csv from frontend ✅
3. Backend to validate data according to listed rules ✅
4. Frontend to receive validation results from backend ✅
5. Frontend to display results and option to commit valid rows ✅
6. Commit valid rows functionality writes to DB ✅
7. Frontend explorer component to display data from backend API ✅
8. Containerize full app with Docker ✅
9. Update README with setup and testing ✅
10. Frontend component to display assets hierarchy ✅
11. Frontend search functionality with tree expansion and highlighting ✅
12. Add logging for errors, requests and imports
13. Review testing comprehensively
14. Review and refactor code for clarity
15. Write architecture, trade-off and limitations
16. Complete all documentation
17. Test fresh setup
18. Host app

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

```
# Backend to validate data according to listed rules

Extend the existing Go backend to validate CSV asset records. Keep the code simple, readable, and organized into separate validation functions executed in this order:

1. **Required fields:** `asset_id`, `asset_type`, and `asset_name` cannot be empty.
2. **Unique IDs:** Detect duplicate `asset_id` values and report every row containing a duplicate.
3. **Supported values:** Validate asset types and operational statuses.
4. **Dates and numbers:** Validate dates and numeric fields. `rating_kva` cannot be negative.
5. **Parent requirements:** Substations cannot have parents; all other assets must have one.
6. **Cycle detection:** Build a map of assets by ID. Detect self-parenting and circular parent-child relationships using a separate DFS function.
7. **Hierarchy validation:** Verify that every parent exists in the CSV, each relationship follows the permitted hierarchy, and every asset ultimately resolves to a substation root.

### Permitted hierarchy

| Asset type | Permitted parent |
|---|---|
| SUBSTATION | None |
| TRANSFORMER | SUBSTATION |
| LV_BOARD | SUBSTATION |
| SWITCHBOARD | SUBSTATION |
| SWITCHBOARD_PANEL | SWITCHBOARD |

Supported statuses: `IN_SERVICE`, `MAINTENANCE`, `OUT_OF_SERVICE`.

### Backend response

- Run all validation functions and collect errors rather than stopping at the first error.
- Return a JSON list of errors containing the original CSV line number, asset ID where available, and meaningful error message.
- Include the total rows processed and whether validation passed.
- No database or persistence.

### Frontend

- Display all validation errors as a numbered list, e.g. "Line 12: Asset TR-001 has an invalid parent."
- If no errors exist, display the existing success message with the row count.
- Preserve the existing loading behavior and design.

Keep each validation function independent and easy to test. Avoid unnecessary dependencies, abstractions, or unrelated changes.
```

## Rejected suggestions and why
1. In the CSV upload component the check that the file is a CSV only checked that it ends with .csv. I prompted for a better check as files can be mistakenly named. While these "simpler" validations can be made on the backend, my intention is to catch errors early and reduce unnecessary calls to the backend. I also prompted for missing/duplicate columns checks.

2. I prompted for paths for each asset type and the agent produced a new page on the frontend for each asset type /substations, /transformers and so on. I then prompted for dynamic paths under /asset-type/[assetType] instead as the pages contained the same code.

## Reviewing and verification of generated code
I did not allow the agent (Claude Code) to run on its own with permission to commit and push code directly. Every commit was made by me after I had run the code to verify the change I had intended with the prompt. I did this to ensure I could meet the deadline and avoid having to revert and redo poorly executed prompts.

**How I worked**
1. I wrote a description of the functionality/change I intended to make and gave it to ChatGPT to write a more precise prompt for Claude Code.
2. I kept each change request small enough for me to be able to review the code change but also large enough to make reasonable progress.
3. When I wasn't sure of how a change would be implemented, I would prompt the agent for implementation options. I would choose an implementation option that solves the problem and requires minimal code changes if possible.

**Verifying a change**
1. I would run the app to test the change.
2. I then read the code changes for each file and if I did not understand something, I would clarify it with the agent.
3. If I thought there were simple improvements that could be made, I would keep prompting for the small changes.
4. If I knew my idea for an improvement would be a bigger code change, I would commit the current changes first.


## What the agent did poorly
1. The agent doesn't refactor code for human readability very well. It tends to overdo the refactor. Of course refactoring is also a matter of personal taste.

## Approximate human working time

| Day | Work done | Human working time |
|-----|-------------------------|--------------------|
| Day 0 | Read requirements, planning, setup | 2 hours |
| Day 1 | Upload CSV, backend | 3 hours |
| Day 2 | Validation, commit valid rows, display data | 4 hours |
| Day 3 | Asset type pages, asset card design, lazy loading, CSV generator | 6 hours |
| Day 4 | Rest | - |
| Day 5 | Docker, setup script, asset hierarchy, layout | 5 hours |
| Day 6 | Search, highlighting | TBD |
| Day 7 | TBD | TBD |
| **Total** | | **TBD** |