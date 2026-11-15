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

// The page of an asset in the explorer, beside the hierarchy tree.
export function explorerPath(assetId: string): string {
  return `${EXPLORER_PATH}/${encodeURIComponent(assetId)}`;
}

// Writes a status as ordinary words, e.g. "OUT_OF_SERVICE" as "Out of service".
export function statusLabel(status: string): string {
  const words = status.replaceAll("_", " ").toLowerCase();
  return words.charAt(0).toUpperCase() + words.slice(1);
}
