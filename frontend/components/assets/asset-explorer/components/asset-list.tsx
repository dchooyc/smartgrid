import Link from "next/link";
import { fetchAssets } from "../lib/api";
import { assetPath } from "../lib/fields";
import type { AssetListing } from "../lib/types";
import { useLoaded } from "../lib/use-loaded";
import AssetCard from "./asset-card";
import AssetFields from "./asset-fields";

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
      <p className="text-muted">{assets.data.length} {listing.title}</p>
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
