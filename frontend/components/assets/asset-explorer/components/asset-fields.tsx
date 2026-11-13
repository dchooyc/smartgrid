import type { Asset } from "../lib/api";
import { ASSET_FIELDS } from "../lib/fields";

type AssetFieldsProps = {
  asset: Asset;
  skipEmpty?: boolean; // leave out fields that have no value
};

// The fields of an asset as plain text, with the required ones in bold.
export default function AssetFields({ asset, skipEmpty }: AssetFieldsProps) {
  const fields = ASSET_FIELDS.filter(
    ({ key }) => !skipEmpty || asset[key] !== null,
  );

  return (
    <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1">
      {fields.map(({ key, label, required }) => (
        <div key={key} className="col-span-2 grid grid-cols-subgrid">
          <dt className="text-muted">{label}</dt>
          <dd className={valueStyle(key, required)}>{asset[key] ?? "—"}</dd>
        </div>
      ))}
    </dl>
  );
}

// The type is also in its own colour, which only shows inside an AssetCard.
function valueStyle(key: keyof Asset, required?: boolean): string | undefined {
  if (!required) return undefined;
  return key === "asset_type"
    ? "font-semibold text-(--asset-text)"
    : "font-semibold";
}
