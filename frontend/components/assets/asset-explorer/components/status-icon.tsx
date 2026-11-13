import type { IconType } from "react-icons";
import { FaCheck, FaWrench, FaXmark } from "react-icons/fa6";

const STATUS_ICONS: Record<string, { Icon: IconType; color: string }> = {
  IN_SERVICE: { Icon: FaCheck, color: "text-green-600" },
  MAINTENANCE: { Icon: FaWrench, color: "text-yellow-500" },
  OUT_OF_SERVICE: { Icon: FaXmark, color: "text-red-600" },
};

// The statuses that have an icon, from working to not.
export const STATUSES = Object.keys(STATUS_ICONS);

type StatusIconProps = {
  status: string;
  className?: string; // replaces the default size, which suits a line of text
};

// A small coloured icon for an operational status, for use beside or near
// the status text.
export default function StatusIcon({
  status,
  className = "inline size-4 align-[-2px]",
}: StatusIconProps) {
  const entry = STATUS_ICONS[status];
  if (!entry) return null;

  const { Icon, color } = entry;
  return <Icon aria-hidden className={`${className} ${color}`} />;
}
