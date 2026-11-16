import type { Asset } from "./api";
import { formatDate, formatTimestamp } from "./dates";

// The asset fields, in the order they are displayed.
export const ASSET_FIELDS: {
  key: keyof Asset;
  label: string;
  required?: boolean; // every asset has a value for it
  format?: (value: string) => string; // how its value is displayed
}[] = [
  { key: "asset_id", label: "Asset ID", required: true },
  { key: "parent_asset_id", label: "Parent asset ID" },
  { key: "asset_type", label: "Type", required: true },
  { key: "asset_name", label: "Name", required: true },
  { key: "operational_status", label: "Status" },
  { key: "voltage_kv", label: "Voltage (kV)" },
  { key: "rating_kva", label: "Rating (kVA)" },
  { key: "manufacturer", label: "Manufacturer" },
  { key: "model", label: "Model" },
  { key: "serial_number", label: "Serial number" },
  { key: "commissioned_date", label: "Commissioned date", format: formatDate },
  { key: "created_at", label: "Created at", format: formatTimestamp },
];

// The page of an asset on its own.
export function assetPath(assetId: string): string {
  return `/assets/${encodeURIComponent(assetId)}`;
}

export const EXPLORER_PATH = "/explorer";

// The name of the URL parameter that carries the search an asset was found
// with, so that its page can highlight the matching text.
export const SEARCH_PARAM = "q";

// The page of an asset in the explorer, beside the hierarchy tree. With the
// search term it was found with, the page highlights the matching text.
export function explorerPath(assetId: string, searchTerm?: string): string {
  const path = `${EXPLORER_PATH}/${encodeURIComponent(assetId)}`;
  if (!searchTerm) return path;
  return `${path}?${new URLSearchParams({ [SEARCH_PARAM]: searchTerm })}`;
}

// Writes a status as ordinary words, e.g. "OUT_OF_SERVICE" as "Out of service".
export function statusLabel(status: string): string {
  const words = status.replaceAll("_", " ").toLowerCase();
  return words.charAt(0).toUpperCase() + words.slice(1);
}
