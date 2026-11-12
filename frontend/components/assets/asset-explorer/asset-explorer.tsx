"use client";

import Link from "next/link";
import { AssetDetails, AssetList } from "./components";

type AssetExplorerProps = {
  assetId?: string; // the selected asset; without one, all assets are listed
};

export default function AssetExplorer({ assetId }: AssetExplorerProps) {
  if (assetId === undefined) {
    return (
      <div className="flex flex-col gap-4 text-foreground">
        <h1 className="text-2xl font-semibold">Assets</h1>
        <AssetList />
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6 text-foreground">
      <Link href="/assets" className="self-start text-accent underline">
        All assets
      </Link>
      <h1 className="text-2xl font-semibold">Asset {assetId}</h1>
      {/* The key reloads the details when another asset is selected. */}
      <AssetDetails key={assetId} assetId={assetId} />
    </div>
  );
}
