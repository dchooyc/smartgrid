import { useSearchParams } from "next/navigation";
import { AssetCard, AssetFields } from "@/components/assets/asset-card";
import { fetchAssetDetails } from "../lib/api";
import { SEARCH_PARAM } from "../lib/fields";
import { useLoaded } from "../lib/use-loaded";
import { Descendants } from "./asset-children";
import AssetParent from "./asset-parent";
import DescendantCounts from "./descendant-counts";

// The parent card is indented from the asset's own card, by the same space
// that opened children and parents get.
const RELATED = "ml-[42px]";

export default function AssetDetails({ assetId }: { assetId: string }) {
  const details = useLoaded(() => fetchAssetDetails(assetId));
  // The search this asset was found with, if it was opened from one.
  const searchTerm = useSearchParams().get(SEARCH_PARAM)?.trim();

  if (details.state === "loading") {
    return (
      <p className="appear-late min-h-screen text-muted">
        Loading asset {assetId}...
      </p>
    );
  }
  if (details.state === "error") {
    return <p className="font-medium text-red-700">{details.message}</p>;
  }
  const { asset, parent, children, descendant_counts } = details.data;

  return (
    <>
      <section className="fade-in flex flex-col gap-2">
        <h2 className="text-lg font-semibold">Details</h2>
        <AssetCard asset={asset}>
          <AssetFields asset={asset} highlight={searchTerm} />
        </AssetCard>
      </section>

      <section className="fade-in flex flex-col gap-2">
        <h2 className="text-lg font-semibold">Parent</h2>
        {parent ? (
          <div className={RELATED}>
            <AssetParent asset={parent} />
          </div>
        ) : (
          <p className="text-muted">This asset has no parent.</p>
        )}
      </section>

      <section className="fade-in flex flex-col gap-2">
        <h2 className="text-lg font-semibold">Children ({children.length})</h2>
        {children.length === 0 ? (
          <p className="text-muted">This asset has no children.</p>
        ) : (
          <>
            <DescendantCounts counts={descendant_counts} />
            <Descendants asset={asset} loaded={children} />
          </>
        )}
      </section>
    </>
  );
}
