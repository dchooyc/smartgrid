import { useState } from "react";
import { FaChevronDown, FaChevronRight } from "react-icons/fa6";
import { fetchAssetDetails, type Asset } from "../lib/api";
import { useLoaded } from "../lib/use-loaded";
import AssetSummary from "./asset-summary";

// The children of an asset as cards. A child that has children of its own
// can be opened to show them beneath it, indented, and so on down the
// hierarchy.
export default function AssetChildren({ assets }: { assets: Asset[] }) {
  return (
    <ul className="flex flex-col gap-4">
      {assets.map((asset) => (
        <li key={asset.asset_id} className="flex flex-col gap-2">
          <AssetSummary asset={asset} />
          {asset.child_count > 0 && <Descendants asset={asset} />}
        </li>
      ))}
    </ul>
  );
}

// The button that opens and closes an asset's children.
function Descendants({ asset }: { asset: Asset }) {
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
          <LoadedChildren assetId={asset.asset_id} />
        </div>
      )}
    </>
  );
}

// Fetches an asset's children when they are first shown.
function LoadedChildren({ assetId }: { assetId: string }) {
  const details = useLoaded(() => fetchAssetDetails(assetId));

  if (details.state === "loading") {
    return <p className="text-sm text-muted">Loading children...</p>;
  }
  if (details.state === "error") {
    return (
      <p className="text-sm font-medium text-red-700">{details.message}</p>
    );
  }
  return <AssetChildren assets={details.data.children} />;
}
