import { AssetExplorer } from "@/components/assets";

export default function AssetsPage() {
  return (
    <main className="page-column flex flex-1 flex-col py-10 font-sans">
      <div className="flex flex-col gap-8">
        <AssetExplorer />
      </div>
    </main>
  );
}
