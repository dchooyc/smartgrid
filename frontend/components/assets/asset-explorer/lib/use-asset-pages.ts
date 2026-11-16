import { useEffect, useRef, useState } from "react";
import { fetchAssetPage, type Asset, type AssetCounts } from "./api";

type AssetPages = {
  assets: Asset[]; // every asset fetched so far
  counts: AssetCounts | null; // null until the first page arrives
  nextAfter: string | null; // null once the whole list is fetched
  loading: boolean;
  error: string | null; // why the latest page could not be fetched
};

// Fetches a list of assets a page at a time. The first page is fetched when
// the component appears and each call to loadMore adds the next one.
export function useAssetPages(assetType?: string) {
  const [pages, setPages] = useState<AssetPages>({
    assets: [],
    counts: null,
    nextAfter: null,
    loading: true,
    error: null,
  });

  // Where the list has got to. These are refs, not state, so that every
  // caller sees the latest values at once, even one created before the last
  // page arrived. That is what stops a page being fetched or added twice.
  const fetching = useRef(false);
  const next = useRef<string | null | undefined>(undefined); // undefined: nothing fetched yet

  function loadMore() {
    if (fetching.current || next.current === null) return;

    const after = next.current;
    fetching.current = true;
    if (after !== undefined) {
      setPages((previous) => ({ ...previous, loading: true, error: null }));
    }

    fetchAssetPage(assetType, after).then(
      (page) => {
        fetching.current = false;
        next.current = page.nextAfter;
        setPages((previous) => ({
          assets: [
            ...(after === undefined ? [] : previous.assets),
            ...page.assets,
          ],
          counts: page.counts,
          nextAfter: page.nextAfter,
          loading: false,
          error: null,
        }));
      },
      (error) => {
        fetching.current = false;
        setPages((previous) => ({
          ...previous,
          loading: false,
          error: error.message,
        }));
      },
    );
  }

  // Fetches the first page. The effect can run again while the list is kept,
  // e.g. on returning to this page, and must not start the list over then.
  useEffect(() => {
    if (next.current === undefined) loadMore();
    // eslint-disable-next-line react-hooks/exhaustive-deps -- load once
  }, []);

  return { ...pages, loadMore };
}
