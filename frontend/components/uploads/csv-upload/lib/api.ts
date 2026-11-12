const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type ValidationIssue = {
  type: string; // the kind of problem, e.g. "Duplicate asset IDs"
  line: number;
  message: string;
};

export type GeneratedFile = {
  name: string;
  content: string; // CSV text
};

export type ValidationResult = {
  totalRows: number;
  errors: ValidationIssue[];
  warnings: ValidationIssue[]; // do not fail validation
  validFile: GeneratedFile; // rows without errors
  invalidFile: GeneratedFile; // rows with errors
};

// The backend's error type for rows whose asset is already in the database.
const ALREADY_STORED = "Already in database";

// True if some row has an error the user can fix. Rows that are only already
// in the database need no fixing.
export function hasRowsToFix(result: ValidationResult): boolean {
  return result.errors.some((error) => error.type !== ALREADY_STORED);
}

// Throws if the backend cannot be reached or cannot process the file at all.
export async function validateOnBackend(file: File): Promise<ValidationResult> {
  const body = new FormData();
  body.append("file", file);

  let response: Response;
  try {
    response = await fetch(`${API_URL}/api/validate`, { method: "POST", body });
  } catch {
    throw new Error("Could not reach the backend. Is it running?");
  }

  const result = await response.json().catch(() => null);
  if (!response.ok || !result) {
    throw new Error(
      result?.error ?? `The backend returned an error (${response.status}).`,
    );
  }
  return {
    totalRows: result.total_rows,
    errors: result.errors,
    warnings: result.warnings,
    validFile: result.valid_file,
    invalidFile: result.invalid_file,
  };
}

export type CommitResult =
  | { committed: true; rowsCommitted: number }
  | { committed: false; message: string; errors: ValidationIssue[] };

// Sends the valid rows to be validated again and saved to the database.
export async function commitValidRows(file: GeneratedFile): Promise<CommitResult> {
  const body = new FormData();
  body.append("file", new Blob([file.content], { type: "text/csv" }), file.name);

  let response: Response;
  try {
    response = await fetch(`${API_URL}/api/imports`, { method: "POST", body });
  } catch {
    return failed("Could not reach the backend. Is it running?");
  }

  const result = await response.json().catch(() => null);
  if (!response.ok || !result?.committed) {
    return failed(
      result?.error ?? `The backend returned an error (${response.status}).`,
      result?.errors,
    );
  }
  return { committed: true, rowsCommitted: result.rows_committed };
}

function failed(message: string, errors: ValidationIssue[] = []): CommitResult {
  return { committed: false, message, errors };
}
