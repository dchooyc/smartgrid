import type { ReactNode } from "react";
import type { IconType } from "react-icons";
import { LuInspectionPanel, LuPlug, LuUtilityPole } from "react-icons/lu";
import { TbCircuitInductor, TbCircuitSwitchOpen } from "react-icons/tb";

// A faint icon for each asset type, shown in the corner of its cards.
const TYPE_ICONS: Record<string, IconType> = {
  SUBSTATION: LuUtilityPole, // where the power lines arrive
  TRANSFORMER: TbCircuitInductor, // the coil symbol from circuit diagrams
  LV_BOARD: LuPlug, // low voltage, ready to use
  SWITCHBOARD: TbCircuitSwitchOpen, // the switch symbol from circuit diagrams
  SWITCHBOARD_PANEL: LuInspectionPanel, // one panel of a switchboard
};

type AssetCardProps = {
  assetType: string;
  children: ReactNode;
};

// A card in the colour of its asset type. The colours are in globals.css.
export default function AssetCard({ assetType, children }: AssetCardProps) {
  const Icon = TYPE_ICONS[assetType];

  return (
    <div
      data-asset-type={assetType}
      className="relative flex flex-col gap-3 rounded border-l-4 border-(--asset-accent) bg-(--asset-tint) px-5 py-4"
    >
      {Icon && (
        <Icon
          aria-hidden
          className="pointer-events-none absolute right-4 top-4 size-10 text-(--asset-accent) opacity-60"
        />
      )}
      {children}
    </div>
  );
}
