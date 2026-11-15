import { useState } from "react";
import { FaChevronDown, FaChevronRight } from "react-icons/fa6";
import { fetchAssetDetails, type Asset } from "../lib/api";
import { useLoaded } from "../lib/use-loaded";
import AssetSummary from "./asset-summary";

// The parent of an asset as a card. A parent that has a parent of its own
// can be opened to show it beneath, indented, and so on up the hierarchy.
export default function AssetParent({ asset }: { asset: Asset }) {
  return (
    <div className="flex flex-col gap-2">
      <AssetSummary asset={asset} />
      {asset.parent_asset_id && <Ancestor asset={asset} />}
    </div>
  );
}

// The button that opens and closes an asset's parent. An asset has only one
// parent, so there is no count.
function Ancestor({ asset }: { asset: Asset }) {
  const [open, setOpen] = useState(false);
  const Chevron = open ? FaChevronDown : FaChevronRight;

  return (
    <>
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen(!open)}
        className="flex cursor-pointer items-center gap-2 self-start text-sm text-muted hover:text-accent"
      >
        <Chevron aria-hidden className="size-3" />
        {open ? "Hide" : "Show"} parent
      </button>
      {open && (
        // The space and line on the left show this belongs to the card above.
        <div className="ml-4 border-l-2 border-border pl-6">
          <LoadedParent assetId={asset.asset_id} />
        </div>
      )}
    </>
  );
}

// Fetches an asset's parent when it is first shown.
function LoadedParent({ assetId }: { assetId: string }) {
  const details = useLoaded(() => fetchAssetDetails(assetId));

  if (details.state === "loading") {
    return <p className="text-sm text-muted">Loading parent...</p>;
  }
  if (details.state === "error") {
    return (
      <p className="text-sm font-medium text-red-700">{details.message}</p>
    );
  }
  if (!details.data.parent) {
    return <p className="text-sm text-muted">This asset has no parent.</p>;
  }
  return <AssetParent asset={details.data.parent} />;
}
