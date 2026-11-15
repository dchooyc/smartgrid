import Link from "next/link";
import type { Asset } from "../lib/api";
import { assetPath } from "../lib/fields";
import AssetCard from "./asset-card";
import AssetFields from "./asset-fields";

// An asset as it appears in a list: its card, with a link to its own page.
export default function AssetSummary({ asset }: { asset: Asset }) {
  return (
    <AssetCard asset={asset}>
      <AssetFields asset={asset} skipEmpty />
      <Link
        href={assetPath(asset.asset_id)}
        className="self-start text-accent underline"
      >
        View details
      </Link>
    </AssetCard>
  );
}
