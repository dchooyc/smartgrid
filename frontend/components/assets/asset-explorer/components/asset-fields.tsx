import type { Asset } from "../lib/api";
import { ASSET_FIELDS } from "../lib/fields";

type AssetFieldsProps = {
  asset: Asset;
  skipEmpty?: boolean; // leave out fields that have no value
};

// Every field of an asset as plain text.
export default function AssetFields({ asset, skipEmpty }: AssetFieldsProps) {
  const fields = ASSET_FIELDS.filter(
    ({ key }) => !skipEmpty || asset[key] !== null,
  );

  return (
    <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1">
      {fields.map(({ key, label }) => (
        <div key={key} className="col-span-2 grid grid-cols-subgrid">
          <dt className="text-muted">{label}</dt>
          <dd>{asset[key] ?? "—"}</dd>
        </div>
      ))}
    </dl>
  );
}
