import type { ReactNode } from "react";
import type { Asset } from "../lib/api";
import { statusLabel } from "../lib/fields";
import { TYPE_ICONS } from "../lib/type-icons";
import StatusIcon from "./status-icon";

// A nudge for each status icon whose shape does not fill its box, putting its
// top level with the top of the card's first line of text.
const STATUS_NUDGES: Record<string, string> = {
  IN_SERVICE: "-translate-y-[3px]",
  OUT_OF_SERVICE: "-translate-y-[3px]",
};

type AssetCardProps = {
  asset: Asset;
  children: ReactNode;
};

// A card in the colour of its asset's type, with the status in the top
// corner and the type's icon in the bottom one. The colours are in
// globals.css.
export default function AssetCard({ asset, children }: AssetCardProps) {
  const typeIcon = TYPE_ICONS[asset.asset_type];
  const statusNudge = STATUS_NUDGES[asset.operational_status] ?? "";

  return (
    <div
      data-asset-type={asset.asset_type}
      className="relative flex flex-col gap-3 rounded border-l-4 border-(--asset-accent) bg-(--asset-tint) px-5 py-4"
    >
      {/*
        Both corner icons share one 32px column. The status icon's top is level
        with the top of the first line of text and the type icon's bottom
        rests on the last line. The status is a glanceable copy of the Status
        field; hovering shows its name.
      */}
      <span
        aria-hidden
        className="group absolute right-4 top-[21px] flex w-8 items-center justify-center"
      >
        <StatusIcon
          status={asset.operational_status}
          className={`size-5 ${statusNudge}`}
        />
        <span className="pointer-events-none absolute right-full mr-1 whitespace-nowrap rounded border border-border bg-background px-2 py-0.5 text-xs text-muted opacity-0 transition-opacity delay-150 group-hover:opacity-100">
          {statusLabel(asset.operational_status)}
        </span>
      </span>
      {typeIcon && (
        <typeIcon.Icon
          aria-hidden
          className={`pointer-events-none absolute bottom-[21px] right-4 size-8 text-(--asset-accent) opacity-60 ${typeIcon.nudge}`}
        />
      )}
      {children}
    </div>
  );
}
