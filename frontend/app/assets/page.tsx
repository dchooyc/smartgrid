import Link from "next/link";
import { AssetExplorer } from "@/components/assets";

export default function AssetsPage() {
  return (
    <>
      <Link href="/" className="self-start text-accent underline">
        Upload a CSV
      </Link>
      <AssetExplorer />
    </>
  );
}
