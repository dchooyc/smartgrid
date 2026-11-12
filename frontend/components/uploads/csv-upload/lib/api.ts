const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type ValidationIssue = {
  type: string; // the kind of problem, e.g. "Duplicate asset IDs"
  line: number;
  message: string;
};

export type ValidationResult = {
  totalRows: number;
  errors: ValidationIssue[];
  warnings: ValidationIssue[]; // do not fail validation
};

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
  };
}
