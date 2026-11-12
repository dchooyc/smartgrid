import type { Asset } from "./api";

// The asset fields, in the order they are displayed.
export const ASSET_FIELDS: { key: keyof Asset; label: string }[] = [
  { key: "asset_id", label: "Asset ID" },
  { key: "parent_asset_id", label: "Parent asset ID" },
  { key: "asset_type", label: "Type" },
  { key: "asset_name", label: "Name" },
  { key: "operational_status", label: "Status" },
  { key: "voltage_kv", label: "Voltage (kV)" },
  { key: "rating_kva", label: "Rating (kVA)" },
  { key: "manufacturer", label: "Manufacturer" },
  { key: "model", label: "Model" },
  { key: "serial_number", label: "Serial number" },
  { key: "commissioned_date", label: "Commissioned date" },
  { key: "created_at", label: "Created at" },
];

export function assetPath(assetId: string): string {
  return `/assets/${encodeURIComponent(assetId)}`;
}
