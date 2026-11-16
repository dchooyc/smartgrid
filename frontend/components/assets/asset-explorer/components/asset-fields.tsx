import type { ReactNode } from "react";
import type { Asset } from "../lib/api";
import { ASSET_FIELDS } from "../lib/fields";

type AssetFieldsProps = {
  asset: Asset;
  skipEmpty?: boolean; // leave out fields that have no value
  highlight?: string; // a search term to mark where the ID and name match it
};

// The fields that search looks in.
const SEARCHED: (keyof Asset)[] = ["asset_id", "asset_name"];

// The fields of an asset as plain text, with the required ones in bold.
export default function AssetFields({
  asset,
  skipEmpty,
  highlight,
}: AssetFieldsProps) {
  const fields = ASSET_FIELDS.filter(
    ({ key }) => !skipEmpty || asset[key] !== null,
  );

  return (
    <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1">
      {fields.map(({ key, label, required, format }) => {
        const value = asset[key];
        return (
          <div key={key} className="col-span-2 grid grid-cols-subgrid">
            <dt className="text-muted">{label}</dt>
            <dd className={valueStyle(key, required)}>
              {value === null ? (
                "—"
              ) : format ? (
                format(String(value))
              ) : highlight && SEARCHED.includes(key) ? (
                <Highlighted text={String(value)} term={highlight} />
              ) : (
                value
              )}
            </dd>
          </div>
        );
      })}
    </dl>
  );
}

// The text with every place it contains the term marked, ignoring case.
function Highlighted({ text, term }: { text: string; term: string }) {
  const lowerText = text.toLowerCase();
  const lowerTerm = term.toLowerCase();
  const parts: ReactNode[] = [];

  let from = 0;
  for (
    let at = lowerText.indexOf(lowerTerm);
    at !== -1;
    at = lowerText.indexOf(lowerTerm, from)
  ) {
    parts.push(text.slice(from, at));
    parts.push(
      <mark key={at} className="rounded-sm bg-primary/45 px-0.5 text-inherit">
        {text.slice(at, at + term.length)}
      </mark>,
    );
    from = at + term.length;
  }
  parts.push(text.slice(from));
  return <>{parts}</>;
}

// The type is also in its own colour, which only shows inside an AssetCard.
function valueStyle(key: keyof Asset, required?: boolean): string | undefined {
  if (!required) return undefined;
  return key === "asset_type"
    ? "font-semibold text-(--asset-text)"
    : "font-semibold";
}
