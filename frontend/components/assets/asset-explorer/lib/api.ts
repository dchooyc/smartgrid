import { apiErrorMessage } from "@/lib/api-error";
import { dataVersion } from "@/lib/data-version";

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
  child_count: number; // assets directly beneath it
};

export type AssetDetails = {
  asset: Asset;
  parent: Asset | null; // null for substations
  children: Asset[]; // immediate children only
  ancestor_ids: string[]; // from its substation down to its parent
  descendant_counts: Record<string, number>; // all assets beneath it, by type
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

// Fetches one page of the assets, or of the assets of one type such as
// "SUBSTATION", starting after the previous page's nextAfter.
export async function fetchAssetPage(
  assetType?: string,
  after?: string,
): Promise<AssetPage> {
  const query = new URLSearchParams({ limit: String(PAGE_SIZE) });
  if (assetType) query.set("type", assetType);
  if (after) query.set("after", after);

  const result = await get(`/api/assets?${query}`);
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

export type SearchResult = {
  assets: Asset[]; // the first matches, ordered by ID
  more: boolean; // more assets match than are listed
};

// Finds the assets whose ID or name contains the term.
export async function searchAssets(term: string): Promise<SearchResult> {
  return get(`/api/assets/search?${new URLSearchParams({ q: term })}`);
}

// Fetches every substation: the top level of the hierarchy.
export async function fetchSubstations(): Promise<Asset[]> {
  return (await get("/api/assets?type=SUBSTATION")).assets;
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
    // The data version keeps the browser from reusing a response it kept
    // from before the stored assets last changed.
    const url = new URL(API_URL + path);
    url.searchParams.set("v", dataVersion());
    response = await fetch(url);
  } catch {
    throw new Error("Could not reach the backend. Is it running?");
  }

  const result = await response.json().catch(() => null);
  if (!response.ok || !result) {
    throw new Error(apiErrorMessage(result, response.status));
  }
  return result;
}
