import Link from "next/link";
import { fetchAssets } from "../lib/api";
import { assetPath } from "../lib/fields";
import { useLoaded } from "../lib/use-loaded";
import AssetFields from "./asset-fields";

export default function AssetList() {
  const assets = useLoaded(fetchAssets);

  if (assets.state === "loading") {
    return <p className="text-muted">Loading assets...</p>;
  }
  if (assets.state === "error") {
    return <p className="font-medium text-red-700">{assets.message}</p>;
  }
  if (assets.data.length === 0) {
    return <p className="text-muted">There are no assets yet.</p>;
  }

  return (
    <>
      <p className="text-muted">{assets.data.length} assets</p>
      <ul className="flex flex-col divide-y divide-border">
        {assets.data.map((asset) => (
          <li key={asset.asset_id} className="flex flex-col gap-2 py-4">
            <AssetFields asset={asset} skipEmpty />
            <Link
              href={assetPath(asset.asset_id)}
              className="self-start text-accent underline"
            >
              View details
            </Link>
          </li>
        ))}
      </ul>
    </>
  );
}
