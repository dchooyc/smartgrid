import { TYPE_ICONS } from "../lib/type-icons";
import { TYPE_NAMES } from "../lib/types";

// How many assets there are beneath an asset, at any depth, of each type
// that it has.
export default function DescendantCounts({
  counts,
}: {
  counts: Record<string, number>;
}) {
  const types = Object.keys(TYPE_ICONS).filter(
    (assetType) => counts[assetType] > 0,
  );

  return (
    <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-muted">
      <span>Beneath this asset:</span>
      <ul className="flex flex-wrap items-center gap-x-4 gap-y-1 tabular-nums">
        {types.map((assetType) => {
          const { Icon } = TYPE_ICONS[assetType];
          return (
            <li key={assetType} className="flex items-center gap-1">
              <Icon aria-hidden className="size-4 opacity-70" />
              {TYPE_NAMES[assetType]} {counts[assetType]}
            </li>
          );
        })}
      </ul>
    </div>
  );
}
