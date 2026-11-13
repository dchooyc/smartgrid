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

// How many assets a list has, however many of them have been fetched.
export type AssetCounts = {
  total: number;
  byType: Record<string, number>;
  byStatus: Record<string, number>;
};

export type AssetPage = {
  assets: Asset[];
  counts: AssetCounts;
  nextAfter: string | null; // fetches the next page; null on the last one
};

// About seven screens of cards: enough that the next page is fetched long
// before the reader gets there, and small enough to arrive at once.
export const PAGE_SIZE = 20;

// Fetches one page of a list of assets, e.g. from "/api/assets" or
// "/api/substations", starting after the previous page's nextAfter.
export async function fetchAssetPage(
  apiPath: string,
  after?: string,
): Promise<AssetPage> {
  const query = new URLSearchParams({ limit: String(PAGE_SIZE) });
  if (after) query.set("after", after);

  const result = await get(`${apiPath}?${query}`);
  return {
    assets: result.assets,
    counts: {
      total: result.counts.total,
      byType: result.counts.by_type,
      byStatus: result.counts.by_status,
    },
    nextAfter: result.next_after,
  };
}

export async function fetchAssetDetails(
  assetId: string,
): Promise<AssetDetails> {
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
