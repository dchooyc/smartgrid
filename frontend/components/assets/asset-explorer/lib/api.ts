const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// An asset as the backend stores it. Fields with no stored value are null.
export type Asset = {
  asset_id: string;
  parent_asset_id: string | null;
  asset_type: string;
  asset_name: string;
  operational_status: string;
  voltage_kv: number | null;
  rating_kva: number | null;
  manufacturer: string | null;
  model: string | null;
  serial_number: string | null;
  commissioned_date: string | null; // year-month-day
  created_at: string;
};

export type AssetDetails = {
  asset: Asset;
  parent: Asset | null; // null for substations
  children: Asset[]; // immediate children only
};

// Fetches a list of assets, e.g. from "/api/assets" or "/api/substations".
export async function fetchAssets(apiPath: string): Promise<Asset[]> {
  return (await get(apiPath)).assets;
}

export async function fetchAssetDetails(assetId: string): Promise<AssetDetails> {
  return get(`/api/assets/${encodeURIComponent(assetId)}`);
}

// Throws with a readable message if the backend cannot answer.
async function get(path: string) {
  let response: Response;
  try {
    response = await fetch(`${API_URL}${path}`);
  } catch {
    throw new Error("Could not reach the backend. Is it running?");
  }

  const result = await response.json().catch(() => null);
  if (!response.ok || !result) {
    throw new Error(
      result?.error ?? `The backend returned an error (${response.status}).`,
    );
  }
  return result;
}
