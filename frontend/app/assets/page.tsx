import { AssetExplorer } from "@/components/assets";

export default function AssetsPage() {
  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-surface font-sans">
      <main className="flex flex-1 w-full max-w-3xl flex-col gap-8 py-10 px-6 bg-background sm:px-16">
        <AssetExplorer />
      </main>
    </div>
  );
}
