"use client";

import Link from "next/link";
import { AssetDetails, AssetList } from "./components";
import { ALL_ASSETS, TYPE_LISTINGS, findTypeListing } from "./lib/types";

type AssetExplorerProps = {
  assetId?: string; // the selected asset; without one, assets are listed
  assetType?: string; // lists only this type of asset, e.g. "substations"
};

export default function AssetExplorer({ assetId, assetType }: AssetExplorerProps) {
  if (assetId === undefined) {
    const listing = assetType ? findTypeListing(assetType) : ALL_ASSETS;
    if (!listing) {
      return (
        <p className="font-medium text-red-700">
          There is no asset type called {assetType}.
        </p>
      );
    }

    return (
      <div className="flex flex-col gap-4 text-foreground">
        <h1 className="text-2xl font-semibold">{listing.title}</h1>
        <nav className="flex flex-wrap gap-x-4 gap-y-1">
          {[ALL_ASSETS, ...TYPE_LISTINGS].map(({ path, title }) =>
            path === listing.path ? (
              <span key={path} className="font-semibold">
                {title}
              </span>
            ) : (
              <Link key={path} href={path} className="text-accent underline">
                {title}
              </Link>
            ),
          )}
        </nav>
        <AssetList key={listing.path} listing={listing} />
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
