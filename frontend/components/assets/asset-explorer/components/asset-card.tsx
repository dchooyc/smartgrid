import type { ReactNode } from "react";

type AssetCardProps = {
  assetType: string;
  children: ReactNode;
};

// A card in the colour of its asset type. The colours are in globals.css.
export default function AssetCard({ assetType, children }: AssetCardProps) {
  return (
    <div
      data-asset-type={assetType}
      className="flex flex-col gap-3 rounded border-l-4 border-(--asset-accent) bg-(--asset-tint) px-5 py-4"
    >
      {children}
    </div>
  );
}
