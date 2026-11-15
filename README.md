# Smart Grid

A full stack web app for importing electrical grid assets data from csv files and displaying the uploaded data.

Built with **Next.js, TypeScript, Go and PostgreSQL**.

## 1. Features

- CSV upload and validation with row-level error reporting.
- Validation of asset fields, duplicate IDs and parent-child relationships.
- Cycle detection and hierarchy validation.
- Downloadable CSV files containing valid and invalid records.
- Duplicate serial number warnings.
- PostgreSQL persistence with atomic transactions.
- Read-only REST APIs for retrieving assets by ID or type.
- Basic asset explorer displaying asset details, direct parents and children.

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
// To be improved..
| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/validate` | Validate CSV and generate results |
| POST | `/api/imports` | Commit validated records |
| GET | `/api/assets` | List all assets |
| GET | `/api/assets?type=...` | Filter assets by type |
| GET | `/api/assets/{assetId}` | Asset details, parent and immediate children |
| GET | `/api/substations` | List substations |
| GET | `/api/transformers` | List transformers |
| GET | `/api/lv-boards` | List LV boards |
| GET | `/api/switchboards` | List switchboards |
| GET | `/api/switchboard-panels` | List switchboard panels |

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
