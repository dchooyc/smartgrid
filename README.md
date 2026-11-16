# Smart Grid

A full stack web app for importing electrical grid assets data from csv files and displaying the uploaded data.

Built with **Next.js, TypeScript, Go and PostgreSQL**.

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


## 2. Tech Stack

| Component | Technology |
|---|---|
| Frontend | Next.js, TypeScript, Tailwind CSS |
| Backend | Go (`net/http`) |
| Database | PostgreSQL |
| Testing | Go standard `testing` package |

## 3. Getting Started

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

## 4. Running Tests

```bash
docker compose run --rm test   # with Docker
./scripts/test.sh              # locally
```

Unit tests cover validation rules, cycle detection, hierarchy checks, CSV generation and import processing.

## 5. API Endpoints

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

## 6. Design Decisions and Trade-offs

- **Partial import:** Valid rows can be committed independently of rejected rows.
- **Validation before commit:** Users can review errors before saving data.
- **Revalidation:** Records are checked again against PostgreSQL before insertion.
- **Atomic transactions:** A failed insertion rolls back the entire batch.
- **Relational hierarchy:** A self-referencing foreign key models parent-child relationships.
- **Minimal dependencies:** Go standard library for HTTP and CSV processing; no ORM.

## 7. Asset Explorer and Search

**Asset explorer:** Basic asset listing and details implemented.

**Expandable hierarchy:** TODO

**Search and navigation:** TODO

## 8. Logging and Integration Tests

TODO

## 9. Known Limitations

- Existing assets cannot be updated through CSV imports.
- Validation currently loads existing assets into memory.
- No import-history audit table.
- No authentication, as it is outside the assignment scope.

## 10. Where to Find Each Requirement

| Criterion | In the app | In the code |
|---|---|---|
| CSV validation | Upload a file at `/` | [validate.go](backend/internal/asset/validate.go), [stored.go](backend/internal/asset/stored.go) |
| Agentic workflow | | [AI_JOURNAL.md](AI_JOURNAL.md) |
| Frontend | Hierarchy tree at `/explorer`, lists at `/assets` | [asset-explorer/](frontend/components/assets/asset-explorer/), [csv-upload/](frontend/components/uploads/csv-upload/) |
| PostgreSQL and hierarchy | | [001_create_assets.sql](backend/migrations/001_create_assets.sql), [store.go](backend/internal/store/store.go) |
| Backend and API | Endpoints in section 5 | [main.go](backend/main.go), [api/](backend/internal/api/) |
| Search | Search bar at `/explorer` | `SearchAssets` in [store.go](backend/internal/store/store.go), [asset-search.tsx](frontend/components/assets/asset-explorer/components/asset-search.tsx) |
| Logging | `docker compose logs backend` | [logging.go](backend/internal/api/logging.go), [api.go](backend/internal/api/api.go) |
| Unit tests | `docker compose run --rm test` | `*_test.go` files in [backend/](backend/) |
| Infrastructure | `docker compose up --build` | [compose.yaml](compose.yaml), [scripts/](scripts/) |
| Documentation | | This README |
