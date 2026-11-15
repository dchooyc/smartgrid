import { Suspense } from "react";
import { AssetExplorer } from "@/components/assets";

type AssetTypePageProps = {
  params: Promise<{ assetType: string }>;
};

export default function AssetTypePage({ params }: AssetTypePageProps) {
  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-surface font-sans">
      <main className="flex flex-1 w-full max-w-3xl flex-col gap-8 py-10 px-6 bg-background sm:px-16">
        {/* The asset type is only known per request, so it is read in Suspense. */}
        <Suspense fallback={<p className="text-muted">Loading assets...</p>}>
          <AssetsOfType params={params} />
        </Suspense>
      </main>
    </div>
  );
}

async function AssetsOfType({ params }: AssetTypePageProps) {
  const { assetType } = await params;
  return <AssetExplorer assetType={decodeURIComponent(assetType)} />;
}
