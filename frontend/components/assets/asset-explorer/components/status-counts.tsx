import { statusLabel } from "../lib/fields";
import StatusIcon, { STATUSES } from "./status-icon";

// How many of the assets are in each status, as an icon and a number.
export default function StatusCounts({
  counts,
}: {
  counts: Record<string, number>;
}) {
  return (
    <ul className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm tabular-nums text-muted/85">
      {STATUSES.map((status) => (
        <li
          key={status}
          title={statusLabel(status)}
          className="flex items-center gap-1"
        >
          <StatusIcon status={status} className="size-3.5 opacity-70" />
          <span className="sr-only">{statusLabel(status)}</span>
          {counts[status] ?? 0}
        </li>
      ))}
    </ul>
  );
}
