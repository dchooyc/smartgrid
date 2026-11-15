import { Suspense } from "react";
import { AssetExplorer } from "@/components/assets";

type AssetPageProps = {
  params: Promise<{ assetId: string }>;
};

export default function AssetPage({ params }: AssetPageProps) {
  return (
    <main className="page-column flex flex-1 flex-col py-10 font-sans">
      <div className="flex flex-col gap-8">
        {/* The asset ID is only known per request, so it is read in Suspense. */}
        <Suspense fallback={<p className="text-muted">Loading asset...</p>}>
          <SelectedAsset params={params} />
        </Suspense>
      </div>
    </main>
  );
}

async function SelectedAsset({ params }: AssetPageProps) {
  const { assetId } = await params;
  return <AssetExplorer assetId={decodeURIComponent(assetId)} />;
}
