import { Suspense, type ReactNode } from "react";
import { AssetSearch, AssetTree } from "@/components/assets";

// The explorer pages share the search bar and the hierarchy tree. They live
// here, not in the pages, so that they stay as they are while moving from one
// asset to another.
export default function ExplorerLayout({ children }: { children: ReactNode }) {
  return (
    <div className="page-column flex flex-1 flex-col gap-8 py-10 font-sans">
      <AssetSearch />
      <div className="flex flex-1 flex-col gap-10 lg:flex-row">
        <aside className="max-h-96 shrink-0 overflow-y-auto [scrollbar-width:thin] lg:sticky lg:top-8 lg:max-h-[calc(100vh-4rem)] lg:w-88 lg:self-start">
          {/* The tree reads the selected asset from the URL, which is only
                known per request. */}
          <Suspense
            fallback={
              <p className="appear-late text-sm text-muted">
                Loading hierarchy...
              </p>
            }
          >
            <AssetTree />
          </Suspense>
        </aside>
        <main className="flex min-w-0 flex-1 flex-col gap-8">{children}</main>
      </div>
    </div>
  );
}
