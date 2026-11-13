import Link from "next/link";
import { fetchAssets } from "../lib/api";
import { assetPath } from "../lib/fields";
import { ALL_ASSETS, type AssetListing } from "../lib/types";
import { useLoaded } from "../lib/use-loaded";
import AssetCard from "./asset-card";
import AssetFields from "./asset-fields";
import StatusCounts from "./status-counts";
import TypeCounts from "./type-counts";

export default function AssetList({ listing }: { listing: AssetListing }) {
  const assets = useLoaded(() => fetchAssets(listing.apiPath));

  if (assets.state === "loading") {
    return <p className="text-muted">Loading {listing.title}...</p>;
  }
  if (assets.state === "error") {
    return <p className="font-medium text-red-700">{assets.message}</p>;
  }
  if (assets.data.length === 0) {
    return <p className="text-muted">There are no {listing.title} yet.</p>;
  }

  return (
    <>
      <div className="flex flex-wrap items-center justify-between gap-x-6 gap-y-1 text-muted">
        <p>
          {assets.data.length} {listing.title}
        </p>
        {listing === ALL_ASSETS ? (
          <TypeCounts assets={assets.data} />
        ) : (
          <StatusCounts assets={assets.data} />
        )}
      </div>
      <ul className="flex flex-col gap-4">
        {assets.data.map((asset) => (
          <li key={asset.asset_id}>
            <AssetCard asset={asset}>
              <AssetFields asset={asset} skipEmpty />
              <Link
                href={assetPath(asset.asset_id)}
                className="self-start text-accent underline"
              >
                View details
              </Link>
            </AssetCard>
          </li>
        ))}
      </ul>
    </>
  );
}
