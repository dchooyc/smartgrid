// The name of each asset type in the plural, for headings and counts.
export const TYPE_NAMES: Record<string, string> = {
  SUBSTATION: "Substations",
  TRANSFORMER: "Transformers",
  LV_BOARD: "LV Boards",
  SWITCHBOARD: "Switchboards",
  SWITCHBOARD_PANEL: "Switchboard Panels",
};

// A list of assets: the page that shows it and the type of asset it lists.
// Without a type it lists every asset.
export type AssetListing = { path: string; assetType?: string; title: string };

export const ALL_ASSETS: AssetListing = { path: "/assets", title: "Assets" };

// The asset types, by the name used for them in page URLs.
const TYPES_BY_URL_NAME: Record<string, string> = {
  substations: "SUBSTATION",
  transformers: "TRANSFORMER",
  "lv-boards": "LV_BOARD",
  switchboards: "SWITCHBOARD",
  "switchboard-panels": "SWITCHBOARD_PANEL",
};

export const TYPE_LISTINGS: AssetListing[] = Object.entries(
  TYPES_BY_URL_NAME,
).map(([urlName, assetType]) => ({
  path: `/asset-type/${urlName}`,
  assetType,
  title: TYPE_NAMES[assetType],
}));

// Finds the listing for an asset type from the URL, e.g. "lv-boards".
export function findTypeListing(assetType: string): AssetListing | undefined {
  return TYPE_LISTINGS.find(({ path }) => path === `/asset-type/${assetType}`);
}
