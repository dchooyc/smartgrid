import { Suspense } from "react";
import { AssetExplorer } from "@/components/assets";

type AssetTypePageProps = {
  params: Promise<{ assetType: string }>;
};

export default function AssetTypePage({ params }: AssetTypePageProps) {
  return (
    <main className="page-column flex flex-1 flex-col py-10 font-sans">
      <div className="flex flex-col gap-8">
        {/* The asset type is only known per request, so it is read in Suspense. */}
        <Suspense
          fallback={
            <p className="appear-late min-h-screen text-muted">
              Loading assets...
            </p>
          }
        >
          <AssetsOfType params={params} />
        </Suspense>
      </div>
    </main>
  );
}

async function AssetsOfType({ params }: AssetTypePageProps) {
  const { assetType } = await params;
  return <AssetExplorer assetType={decodeURIComponent(assetType)} />;
}
