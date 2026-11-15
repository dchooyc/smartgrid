// The name of each asset type in the plural, for headings and counts.
export const TYPE_NAMES: Record<string, string> = {
  SUBSTATION: "Substations",
  TRANSFORMER: "Transformers",
  LV_BOARD: "LV Boards",
  SWITCHBOARD: "Switchboards",
  SWITCHBOARD_PANEL: "Switchboard Panels",
};

// A list of assets: the page that shows it and the backend route it is
// fetched from.
export type AssetListing = { path: string; apiPath: string; title: string };

export const ALL_ASSETS: AssetListing = {
  path: "/assets",
  apiPath: "/api/assets",
  title: "Assets",
};

// The asset types, by the name used for them in URLs.
const TYPE_TITLES: Record<string, string> = {
  substations: "Substations",
  transformers: "Transformers",
  "lv-boards": "LV Boards",
  switchboards: "Switchboards",
  "switchboard-panels": "Switchboard Panels",
};

export const TYPE_LISTINGS: AssetListing[] = Object.keys(TYPE_TITLES).map(
  (assetType) => ({
    path: `/asset-type/${assetType}`,
    apiPath: `/api/${assetType}`,
    title: TYPE_TITLES[assetType],
  }),
);

// Finds the listing for an asset type from the URL, e.g. "lv-boards".
export function findTypeListing(assetType: string): AssetListing | undefined {
  return TYPE_LISTINGS.find(({ path }) => path === `/asset-type/${assetType}`);
}
