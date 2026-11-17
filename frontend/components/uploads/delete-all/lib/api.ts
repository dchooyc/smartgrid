import { apiErrorMessage } from "@/lib/api-error";
import { bumpDataVersion } from "@/lib/data-version";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// Deletes every stored asset and returns how many there were. Throws with a
// readable message if the backend cannot do it.
export async function deleteAllAssets(): Promise<number> {
  let response: Response;
  try {
    response = await fetch(`${API_URL}/api/assets`, { method: "DELETE" });
  } catch {
    throw new Error("Could not reach the backend. Is it running?");
  }

  const result = await response.json().catch(() => null);
  if (!response.ok || !result) {
    throw new Error(apiErrorMessage(result, response.status));
  }
  bumpDataVersion(); // the asset pages must not show what the browser kept
  return result.deleted;
}
