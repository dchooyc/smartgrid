import { TYPE_ICONS } from "../lib/type-icons";

// How many of the assets there are of each type, as an icon and a number.
export default function TypeCounts({
  counts,
}: {
  counts: Record<string, number>;
}) {
  return (
    <ul className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm tabular-nums text-muted/85">
      {Object.entries(TYPE_ICONS).map(([assetType, { Icon }]) => (
        <li
          key={assetType}
          title={assetType}
          className="flex items-center gap-1"
        >
          <Icon aria-hidden className="size-4 opacity-70" />
          <span className="sr-only">{assetType}</span>
          {counts[assetType] ?? 0}
        </li>
      ))}
    </ul>
  );
}
