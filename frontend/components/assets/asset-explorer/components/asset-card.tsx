import type { ReactNode } from "react";
import type { IconType } from "react-icons";
import { LuInspectionPanel, LuPlug, LuUtilityPole } from "react-icons/lu";
import { TbCircuitInductor, TbCircuitSwitchOpen } from "react-icons/tb";
import type { Asset } from "../lib/api";
import { statusLabel } from "../lib/fields";
import StatusIcon from "./status-icon";

// A faint icon for each asset type, shown in the corner of its cards. The
// shapes do not all fill their box, so each has a nudge that rests its
// bottom edge on the card's last line of text.
const TYPE_ICONS: Record<string, { Icon: IconType; nudge: string }> = {
  // where the power lines arrive
  SUBSTATION: { Icon: LuUtilityPole, nudge: "" },
  // the coil symbol from circuit diagrams
  TRANSFORMER: { Icon: TbCircuitInductor, nudge: "translate-y-[11px]" },
  // low voltage, ready to use
  LV_BOARD: { Icon: LuPlug, nudge: "" },
  // the switch symbol from circuit diagrams
  SWITCHBOARD: { Icon: TbCircuitSwitchOpen, nudge: "translate-y-[11px]" },
  // one panel of a switchboard
  SWITCHBOARD_PANEL: { Icon: LuInspectionPanel, nudge: "translate-y-px" },
};

// Likewise, a nudge that puts the top of each status icon level with the top
// of the card's first line of text.
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
