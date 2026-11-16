import { useState } from "react";
import { FaChevronDown, FaChevronRight } from "react-icons/fa6";
import { AssetSummary } from "@/components/assets/asset-card";
import { fetchAssetDetails, type Asset } from "../lib/api";
import { TYPE_ICONS } from "../lib/type-icons";
import { TYPE_NAMES } from "../lib/types";
import { useLoaded } from "../lib/use-loaded";

// The children of an asset as cards, grouped by type in the order of the
// hierarchy. A child that has children of its own can be opened to show them
// beneath it, indented, and so on down the hierarchy.
export default function AssetChildren({ assets }: { assets: Asset[] }) {
  const groups = Object.keys(TYPE_ICONS)
    .map((assetType) => ({
      assetType,
      assets: assets.filter((asset) => asset.asset_type === assetType),
    }))
    .filter((group) => group.assets.length > 0);

  return (
    <div className="flex flex-col gap-6">
      {groups.map((group) => (
        <section key={group.assetType} className="flex flex-col gap-2">
          <h3 className="text-sm font-semibold text-muted">
            {TYPE_NAMES[group.assetType]} ({group.assets.length})
          </h3>
          <ul className="flex flex-col gap-4">
            {group.assets.map((asset) => (
              <li key={asset.asset_id} className="flex flex-col gap-2">
                <AssetSummary asset={asset} />
                {asset.child_count > 0 && <Descendants asset={asset} />}
              </li>
            ))}
          </ul>
        </section>
      ))}
    </div>
  );
}

type DescendantsProps = {
  asset: Asset;
  loaded?: Asset[]; // the asset's children, if they are already fetched
};

// The button that opens and closes an asset's children.
export function Descendants({ asset, loaded }: DescendantsProps) {
  const [open, setOpen] = useState(false);
  const Chevron = open ? FaChevronDown : FaChevronRight;
  const children = asset.child_count === 1 ? "child" : "children";

  return (
    <>
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen(!open)}
        className="flex cursor-pointer items-center gap-2 self-start text-sm text-muted hover:text-accent"
      >
        <Chevron aria-hidden className="size-3" />
        {open ? "Hide" : "Show"} {asset.child_count} {children}
      </button>
      {open && (
        // The space and line on the left show these belong to the card above.
        <div className="ml-4 border-l-2 border-border pl-6">
          {loaded ? (
            <AssetChildren assets={loaded} />
          ) : (
            <LoadedChildren assetId={asset.asset_id} />
          )}
        </div>
      )}
    </>
  );
}

// Fetches an asset's children when they are first shown.
function LoadedChildren({ assetId }: { assetId: string }) {
  const details = useLoaded(() => fetchAssetDetails(assetId));

  if (details.state === "loading") {
    return (
      <p className="appear-late text-sm text-muted">Loading children...</p>
    );
  }
  if (details.state === "error") {
    return (
      <p className="text-sm font-medium text-red-700">{details.message}</p>
    );
  }
  return <AssetChildren assets={details.data.children} />;
}
