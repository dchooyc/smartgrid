import { Suspense } from "react";
import { AssetExplorer } from "@/components/assets";

type AssetPageProps = {
  params: Promise<{ assetId: string }>;
};

export default function AssetPage({ params }: AssetPageProps) {
  // The asset ID is only known per request, so it is read in Suspense.
  return (
    <Suspense fallback={<p className="text-muted">Loading asset...</p>}>
      <SelectedAsset params={params} />
    </Suspense>
  );
}

async function SelectedAsset({ params }: AssetPageProps) {
  const { assetId } = await params;
  return <AssetExplorer assetId={decodeURIComponent(assetId)} />;
}
