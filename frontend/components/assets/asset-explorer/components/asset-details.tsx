import Link from "next/link";
import { fetchAssetDetails, type Asset } from "../lib/api";
import { assetPath } from "../lib/fields";
import { useLoaded } from "../lib/use-loaded";
import AssetCard from "./asset-card";
import AssetChildren from "./asset-children";
import AssetFields from "./asset-fields";
import StatusIcon from "./status-icon";

export default function AssetDetails({ assetId }: { assetId: string }) {
  const details = useLoaded(() => fetchAssetDetails(assetId));

  if (details.state === "loading") {
    return <p className="text-muted">Loading asset {assetId}...</p>;
  }
  if (details.state === "error") {
    return <p className="font-medium text-red-700">{details.message}</p>;
  }
  const { asset, parent, children } = details.data;

  return (
    <>
      <section className="flex flex-col gap-2">
        <h2 className="text-lg font-semibold">Details</h2>
        <AssetCard asset={asset}>
          <AssetFields asset={asset} />
        </AssetCard>
      </section>

      <section className="flex flex-col gap-2">
        <h2 className="text-lg font-semibold">Parent</h2>
        {parent ? (
          <RelatedAsset asset={parent} />
        ) : (
          <p className="text-muted">This asset has no parent.</p>
        )}
      </section>

      <section className="flex flex-col gap-2">
        <h2 className="text-lg font-semibold">Children ({children.length})</h2>
        {children.length === 0 ? (
          <p className="text-muted">This asset has no children.</p>
        ) : (
          <AssetChildren assets={children} />
        )}
      </section>
    </>
  );
}

// The parent: its ID as a link, then a few fields as plain text.
function RelatedAsset({ asset }: { asset: Asset }) {
  return (
    <p>
      <Link href={assetPath(asset.asset_id)} className="text-accent underline">
        {asset.asset_id}
      </Link>
      {` · ${asset.asset_type} · ${asset.asset_name} · ${asset.operational_status} `}
      <StatusIcon status={asset.operational_status} />
    </p>
  );
}
