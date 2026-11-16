import { useEffect, useRef } from "react";
import { PAGE_SIZE } from "../lib/api";
import { ALL_ASSETS, type AssetListing } from "../lib/types";
import { useAssetPages } from "../lib/use-asset-pages";
import AssetSummary from "./asset-summary";
import StatusCounts from "./status-counts";
import TypeCounts from "./type-counts";

// The next page is fetched when the reader comes within about two screens of
// the end of the list, so that it is in place before they reach it.
const LOAD_AHEAD = "1500px";

export default function AssetList({ listing }: { listing: AssetListing }) {
  const { assets, counts, nextAfter, loading, error, loadMore } = useAssetPages(
    listing.assetType,
  );
  const end = useRef<HTMLDivElement>(null);

  // Watches for the end of the list coming near. It is set up afresh after
  // every page, which also covers a page too short to push the end away.
  useEffect(() => {
    if (!end.current) return;

    const observer = new IntersectionObserver(
      ([entry]) => entry.isIntersecting && loadMore(),
      { rootMargin: `${LOAD_AHEAD} 0px` },
    );
    observer.observe(end.current);
    return () => observer.disconnect();
  });

  if (!counts) {
    return error ? (
      <p className="font-medium text-red-700">{error}</p>
    ) : (
      <p className="appear-late min-h-screen text-muted">
        Loading {listing.title}...
      </p>
    );
  }
  if (counts.total === 0) {
    return <p className="text-muted">There are no {listing.title} yet.</p>;
  }

  return (
    <>
      <div className="fade-in flex flex-wrap items-center justify-between gap-x-6 gap-y-1 text-muted">
        <p>
          {counts.total} {listing.title}
        </p>
        {listing === ALL_ASSETS ? (
          <TypeCounts counts={counts.byType} />
        ) : (
          <StatusCounts counts={counts.byStatus} />
        )}
      </div>
      {/* Two cards side by side where the page is wide enough. */}
      <ul className="fade-in grid gap-4 lg:grid-cols-2">
        {assets.map((asset) => (
          <li key={asset.asset_id}>
            <AssetSummary asset={asset} />
          </li>
        ))}
      </ul>

      {error ? (
        <p className="text-center text-sm text-red-700">
          {error}{" "}
          <button
            type="button"
            onClick={loadMore}
            className="cursor-pointer text-accent underline"
          >
            Try again
          </button>
        </p>
      ) : nextAfter !== null ? (
        <div ref={end} className="h-5 text-center text-sm text-muted/85">
          {loading && `Loading more ${listing.title}...`}
        </div>
      ) : (
        assets.length > PAGE_SIZE && (
          <p className="text-center text-sm text-muted/85">
            All {counts.total} {listing.title} shown
          </p>
        )
      )}
    </>
  );
}
