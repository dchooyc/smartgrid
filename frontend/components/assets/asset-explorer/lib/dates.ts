const MONTHS = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

// Writes a year-month-day date as e.g. "20 Jun 2011".
export function formatDate(value: string): string {
  const [year, month, day] = value.split("-").map(Number);
  if (!year || !MONTHS[month - 1] || !day) return value;
  return `${day} ${MONTHS[month - 1]} ${year}`;
}

// Singapore time on a 12-hour clock, whatever the viewer's own time zone.
const singaporeTime = new Intl.DateTimeFormat("en-SG", {
  timeZone: "Asia/Singapore",
  day: "numeric",
  month: "short",
  year: "numeric",
  hour: "numeric",
  minute: "2-digit",
  second: "2-digit",
  hour12: true,
});

// Writes a timestamp as e.g. "12 Nov 2026, 2:58:53 pm".
export function formatTimestamp(value: string): string {
  const timestamp = new Date(value);
  if (Number.isNaN(timestamp.getTime())) return value;

  const parts: Record<string, string> = {};
  for (const { type, value } of singaporeTime.formatToParts(timestamp)) {
    parts[type] = value;
  }
  const { day, month, year, hour, minute, second, dayPeriod } = parts;
  return `${day} ${month} ${year}, ${hour}:${minute}:${second} ${dayPeriod.toLowerCase()}`;
}
