import { Suspense, type ReactNode } from "react";
import { AssetTree } from "@/components/assets";

// The asset pages share the hierarchy tree. It lives here, not in the pages,
// so that it stays as it is while moving from one asset to another.
export default function AssetsLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex flex-1 justify-center bg-surface font-sans">
      <div className="flex w-full max-w-6xl flex-col gap-10 bg-background px-6 py-16 sm:px-16 lg:flex-row">
        <aside className="max-h-96 shrink-0 overflow-y-auto lg:sticky lg:top-8 lg:max-h-[calc(100vh-4rem)] lg:w-88 lg:self-start">
          {/* The tree reads the selected asset from the URL, which is only
              known per request. */}
          <Suspense
            fallback={
              <p className="text-sm text-muted">Loading hierarchy...</p>
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
