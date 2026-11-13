import type { IconType } from "react-icons";
import { LuInspectionPanel, LuPlug, LuUtilityPole } from "react-icons/lu";
import { TbCircuitInductor, TbCircuitSwitchOpen } from "react-icons/tb";

// An icon for each asset type, in the order of the hierarchy. The shapes do
// not all fill their box, so each has a nudge that rests its bottom edge on
// the last line of text when it is shown in the corner of a card.
export const TYPE_ICONS: Record<string, { Icon: IconType; nudge: string }> = {
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
