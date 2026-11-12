import { Suspense } from "react";
import { AssetExplorer } from "@/components/assets";

type AssetPageProps = {
  params: Promise<{ assetId: string }>;
};

export default function AssetPage({ params }: AssetPageProps) {
  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-surface font-sans">
      <main className="flex flex-1 w-full max-w-3xl flex-col gap-8 py-16 px-6 bg-background sm:px-16">
        {/* The asset ID is only known per request, so it is read in Suspense. */}
        <Suspense fallback={<p className="text-muted">Loading asset...</p>}>
          <SelectedAsset params={params} />
        </Suspense>
      </main>
    </div>
  );
}

async function SelectedAsset({ params }: AssetPageProps) {
  const { assetId } = await params;
  return <AssetExplorer assetId={decodeURIComponent(assetId)} />;
}
