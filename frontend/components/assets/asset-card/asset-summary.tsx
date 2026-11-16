import Link from "next/link";
import { usePathname } from "next/navigation";
import type { Asset } from "../lib/api";
import { EXPLORER_PATH, assetPath, explorerPath } from "../lib/fields";
import AssetCard from "./asset-card";
import AssetFields from "./asset-fields";

// An asset as it appears in a list: its card, with a link to its own page.
export default function AssetSummary({ asset }: { asset: Asset }) {
  // Inside the explorer the link stays in the explorer.
  const inExplorer = usePathname().startsWith(EXPLORER_PATH);
  const path = inExplorer ? explorerPath : assetPath;

  return (
    <AssetCard asset={asset}>
      <AssetFields asset={asset} skipEmpty />
      <Link
        href={path(asset.asset_id)}
        className="self-start text-accent underline"
      >
        View details
      </Link>
    </AssetCard>
  );
}
