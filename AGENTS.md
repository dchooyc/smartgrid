# Agent instructions

Smart Grid imports electrical asset data from CSV files, validates it, stores it in PostgreSQL and shows the asset hierarchy. The frontend is Next.js with TypeScript (`frontend/`), the backend is Go (`backend/`) and the helper scripts are in `scripts/`.

Read `README.md` for what the app does. Read `frontend/AGENTS.md` before changing frontend code.

## Architecture conventions

**Backend** (`backend/internal/`)
- `api` is the only package with HTTP: handlers, request and response shapes, logging.
- `store` is the only package with SQL.
- `asset` is pure logic: the asset model, the validation rules and the hierarchy rules. It has no HTTP and no SQL.
- `csvfile` reads uploaded CSV files and builds the result files.
- Use the Go standard library. Do not add a framework, an ORM or any other dependency without asking first.
- A failure on the server's side goes through `fail`, which logs the cause. A mistake in the request goes through `writeError` and is not logged as a failure.
- An import validates and saves inside one transaction: all rows or none.
- A schema change is a new numbered file in `backend/migrations/`. Never edit a migration that has been applied.

**Frontend** (`frontend/`)
- Pages are in `app/` and stay thin. Components are in `components/`, one folder per feature, each with an `index.ts` that exports what other folders may use.
- Import with `./` inside a folder and with `@/components/...` across folders.
- All calls to the backend go through a `lib/api.ts` file. Components do not call `fetch` themselves.
- The selected asset is kept in the URL, not in component state.
- The theme colours are defined in `app/globals.css`. Use them before adding a new colour.

**Both**
- Keep the code simple and readable. Prefer the smallest change that solves the problem.
- Comments say why, in plain words.

## Development commands

```bash
./scripts/setup.sh                # once: database, migrations, dependencies
./scripts/dev.sh                  # frontend on :3000, backend on :8080
docker compose up --build         # the whole app in containers
./scripts/migrate.sh              # apply new migrations
./scripts/auto-generate-csvs.sh   # five sample CSV files in generated-csv/
./scripts/generate-csv.sh         # one sample CSV file, with chosen problems
./scripts/delete-assets.sh        # empty the development database
```

## Test commands

```bash
./scripts/test.sh                          # backend tests, locally
docker compose run --rm test               # backend tests, in Docker
cd backend && gofmt -l . && go vet ./...   # formatting and static checks
cd frontend && npx tsc --noEmit            # type check
cd frontend && npm run lint                # lint
```

## Constraints

- Never commit, push or stage. The human reviews and commits every change.
- Make one small change per request. Do not modify unrelated files.
- When asked not to write code, answer only. When the approach is unclear, give options before implementing.
- Never write to the development database to test something. Use read-only requests, or create a throwaway database and drop it afterwards.
- Do not run `next build` while the dev server is running.
- Start a server only when its port is free, and stop only the processes you started.
- Never log the contents of an uploaded file. Log the file name and counts only.
- Do not run `scripts/delete-assets.sh` unless asked.

## Definition of done

A change is done when all of these are true:

1. The backend builds, and `gofmt`, `go vet` and the tests pass.
2. The frontend type check and lint pass.
3. The change was exercised in the running app, not only compiled.
4. A change to a validation rule comes with a unit test.
5. `README.md` is updated if behaviour, setup or the API changed.
6. The reply says what changed and how to test it.
7. Nothing is committed or staged.
