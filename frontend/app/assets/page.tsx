import Link from "next/link";
import { AssetExplorer } from "@/components/assets";

export default function AssetsPage() {
  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-surface font-sans">
      <main className="flex flex-1 w-full max-w-3xl flex-col gap-8 py-16 px-6 bg-background sm:px-16">
        <Link href="/" className="self-start text-accent underline">
          Upload a CSV
        </Link>
        <AssetExplorer />
      </main>
    </div>
  );
}
