# Smart Grid

A full stack web app for importing electrical grid assets data from csv files and displaying the uploaded data.

Built with **Next.js, TypeScript, Go and PostgreSQL**.

![The Explorer: the hierarchy tree, search and the selected asset's details](images/explorer.png)

![The upload page, where a CSV file is chosen and validated](images/upload.png)

## 1. Features

**Import**
- Upload a CSV file and see every error and warning by line.
- Checks fields, duplicate IDs, parent-child rules and cycles.
- Download the valid and invalid rows as separate files.
- Save the valid rows to PostgreSQL in one transaction: all rows or none.

**Explore**
- Expandable hierarchy tree, from substations down to switchboard panels.
- Asset details with parent, children grouped by type and counts of everything beneath.
- Asset lists, all together or by type.

**Search**
- Find assets by part of their ID or name.
- Choosing a result opens the asset, expands the tree to it and highlights it.
- Every asset has its own address, which survives a refresh.

**Logging**
- The backend logs every request it receives: what was asked for, the result and how long it took.
- When something fails on the server, the log says why.
- Every validation and import is logged with its file name and row counts.
- An error shown in the app ends with a reference, which finds that failure in the logs.


## 2. Where to Find Each Requirement

| Criterion | In the app | In the code |
|---|---|---|
| CSV validation | Upload a file at `/` | [validate.go](backend/internal/asset/validate.go), [stored.go](backend/internal/asset/stored.go) |
| Agentic workflow | | [AI_JOURNAL.md](AI_JOURNAL.md) |
| Frontend | Hierarchy tree at `/explorer`, lists at `/assets` | [asset-explorer/](frontend/components/assets/asset-explorer/), [csv-upload/](frontend/components/uploads/csv-upload/) |
| PostgreSQL and hierarchy | | [001_create_assets.sql](backend/migrations/001_create_assets.sql), [store.go](backend/internal/store/store.go) |
| Backend and API | Endpoints in section 6 | [main.go](backend/main.go), [api/](backend/internal/api/) |
| Search | Search bar at `/explorer` | `SearchAssets` in [store.go](backend/internal/store/store.go), [asset-search.tsx](frontend/components/assets/asset-explorer/components/asset-search.tsx) |
| Logging | `docker compose logs backend` | [logging.go](backend/internal/api/logging.go), [api.go](backend/internal/api/api.go) |
| Unit tests | `docker compose run --rm test` | `*_test.go` files in [backend/](backend/) |
| Infrastructure | `docker compose up --build` | [compose.yaml](compose.yaml), [scripts/](scripts/) |
| Documentation | | This README |

## 3. Tech Stack

| Component | Technology |
|---|---|
| Frontend | Next.js, TypeScript, Tailwind CSS |
| Backend | Go (`net/http`) |
| Database | PostgreSQL |
| Testing | Go standard `testing` package |

## 4. Getting Started

**Run with Docker** (requires only Docker)

```bash
docker compose up --build
```

Open http://localhost:3000. Reset the data with `docker compose down -v`.

If port 5432 is already in use: `DB_PORT=5433 docker compose up --build`

**Develop locally** (requires Docker, Go and Node.js)

```bash
./scripts/setup.sh   # once: starts the database, applies migrations, installs dependencies
./scripts/dev.sh     # frontend on :3000, backend on :8080
```

## 5. Running Tests

```bash
docker compose run --rm test   # with Docker
./scripts/test.sh              # locally
```

Unit tests cover validation rules, cycle detection, hierarchy checks, CSV generation and import processing.

## 6. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/validate` | Validate CSV and generate results |
| POST | `/api/imports` | Commit validated records |
| GET | `/api/assets` | List assets, optionally of one type and a page at a time |
| GET | `/api/assets/search?q=...` | Find assets whose ID or name contains the text |
| GET | `/api/assets/{assetId}` | Asset details, parent, immediate children, ancestors and descendant counts |

**Listing assets**

`GET /api/assets` takes three optional query parameters, which can be combined:

| Parameter | Meaning | Example |
|---|---|---|
| `type` | Only assets of this type: `SUBSTATION`, `TRANSFORMER`, `LV_BOARD`, `SWITCHBOARD` or `SWITCHBOARD_PANEL` | `type=TRANSFORMER` |
| `limit` | At most this many assets, from 1 to 100 | `limit=20` |
| `after` | Continue after this asset ID, taken from the previous response's `next_after` | `after=TX-001-2` |

```bash
# Every asset
curl "http://localhost:8080/api/assets"

# Only the substations
curl "http://localhost:8080/api/assets?type=SUBSTATION"

# The first 20 transformers
curl "http://localhost:8080/api/assets?type=TRANSFORMER&limit=20"

# The next 20, where the previous response ended with "next_after": "TX-001-2"
curl "http://localhost:8080/api/assets?type=TRANSFORMER&limit=20&after=TX-001-2"
```

Assets are ordered by ID. `next_after` is `null` on the last page, and an unsupported `type` returns a 400 that lists the supported ones.

**Searching assets**

`GET /api/assets/search?q=...` returns the assets whose ID or name contains the text in `q`. Upper and lower case are treated the same.

```bash
# By name: finds "Aurora Panel A1"
curl "http://localhost:8080/api/assets/search?q=panel"

# By ID: finds "PNL-001-1-01"
curl "http://localhost:8080/api/assets/search?q=pnl-001"
```

It returns the first 20 matches, ordered by ID. `more` is `true` when more assets match than were returned.

## 7. Design Decisions and Trade-offs

**Import**
- **Validate first, commit second.** The user sees every error before anything is saved. The cost is two uploads of the same data.
- **Partial import.** The valid rows can be saved without waiting for the rejected ones to be fixed. An asset is still never stored without its parent: the commit checks this again and the database enforces it.
- **Validated again on commit.** The import does not trust the earlier check: it validates every row again against the database, inside the transaction that saves them.
- **All rows or none.** One transaction per import. A failure part-way leaves the database as it was.

**Database**
- **One table.** Each asset stores the ID of its parent, which is how the hierarchy is kept.
- **The database enforces the rules too.** It refuses a row that breaks them, such as an asset without a parent.
- **One query for a whole branch.** An asset's ancestors, and the counts beneath it, each take a single query.
- **Indexes** on the parent and the serial number, to find children and duplicates quickly.
- **Search reads every row.** That is fast at this size. A large table would need a search index.

**API**
- **One call for an asset's details.** It returns the asset with its parent, children, ancestors and counts, so the page needs one request instead of four. There are no separate `/children` and `/ancestors` endpoints.
- **Paging by "after this ID"**, not by page number, so a page stays correct when assets are added.
- **Standard library only** for HTTP, CSV and logging. No framework and no ORM, so the SQL is in plain view.

**Frontend**
- **The tree loads one level at a time.** Opening an asset fetches its children. The first load is fast however large the hierarchy is, at the cost of a short wait on first opening an asset.
- **The address is the selection.** The selected asset is in the URL, so a link, a refresh and the back button all work.

## 8. Asset Explorer and Search

**Explorer** (`/explorer`)
- The tree on the left starts with the substations. Opening an asset shows the assets beneath it.
- Selecting an asset shows its details on the right: its fields, its parent, its children grouped by type, and how many assets of each type are beneath it.
- Each asset has its own address, `/explorer/{asset_id}`. Opening that address directly, or refreshing, opens the tree down to the asset and highlights it.

**Search**
- The search bar is above the tree. Press `/` to jump to it.
- It matches any part of an asset's ID or name, ignoring case, and lists the first 20 matches as you type.
- Choosing a result opens the asset's page, opens the tree down to it, highlights it in the tree and shows its details, with the matched text marked.

**Asset lists** (`/assets` and one page per type)
- Cards for every asset, with counts by type and status. More load as you scroll.

Every view shows a message while loading, when there is nothing to show and when the backend cannot answer.

## 9. Logging

The backend writes its logs to standard output, one line per event. See them with `docker compose logs backend`, or in the terminal that runs `./scripts/dev.sh`.

**What is logged**

| Event | Level | Contains |
|---|---|---|
| Server start | INFO | The database in use, with the password hidden |
| Every request | INFO, or WARN for a 4xx, or ERROR for a 5xx | Method, path, query, status and duration |
| A server failure | ERROR | What could not be done and the cause |
| A validation | INFO | File name, rows, errors and warnings |
| An import | INFO if committed, WARN if rejected | File name, rows and errors |
| A crash in a handler | ERROR | The cause and where it happened. The server keeps running and answers with a 500. |

Every line for a request carries the same `request_id`. The contents of uploaded files are never logged, only the file name and counts.

## 10. Known Limitations

- Existing assets cannot be updated through CSV imports.
- Validation currently loads existing assets into memory.
- No import-history audit table.
- No authentication, as it is outside the assignment scope.
