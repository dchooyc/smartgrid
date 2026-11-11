const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function validateOnBackend(file: File): Promise<number> {
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
  return result.row_count;
}
