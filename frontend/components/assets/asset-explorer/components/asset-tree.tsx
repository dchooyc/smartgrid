"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { FaChevronDown, FaChevronRight } from "react-icons/fa6";
import { fetchAssetDetails, fetchSubstations, type Asset } from "../lib/api";
import { explorerPath } from "../lib/fields";
import { TYPE_ICONS } from "../lib/type-icons";
import { useLoaded } from "../lib/use-loaded";
import StatusIcon from "./status-icon";

const NO_ANCESTORS: string[] = [];

// The whole hierarchy as a tree, starting from the substations. An asset
// with children can be opened to show them. The asset whose page is open is
// highlighted, with everything above it opened.
export default function AssetTree() {
  const { assetId } = useParams<{ assetId?: string }>();
  const selectedId = assetId && decodeURIComponent(assetId);
  const substations = useLoaded(fetchSubstations);
  const ancestorIds = useAncestorIds(selectedId);

  return (
    <nav aria-label="Asset hierarchy" className="flex flex-col gap-2">
      <h2 className="text-lg font-semibold">Hierarchy</h2>
      {substations.state === "loading" ? (
        <p className="text-sm text-muted">Loading hierarchy...</p>
      ) : substations.state === "error" ? (
        <p className="text-sm font-medium text-red-700">
          {substations.message}
        </p>
      ) : substations.data.length === 0 ? (
        <p className="text-sm text-muted">There are no assets yet.</p>
      ) : (
        <TreeLevel
          assets={substations.data}
          selectedId={selectedId}
          ancestorIds={ancestorIds}
        />
      )}
    </nav>
  );
}

// Fetches the IDs of the assets above the selected one. Until they arrive,
// or with nothing selected, there are none.
function useAncestorIds(assetId: string | undefined): string[] {
  const [found, setFound] = useState<{ assetId: string; ids: string[] }>();

  useEffect(() => {
    if (!assetId) return;

    let cancelled = false;
    fetchAssetDetails(assetId).then(
      (details) =>
        !cancelled && setFound({ assetId, ids: details.ancestor_ids }),
      () => {}, // the details area reports the error
    );
    return () => {
      cancelled = true;
    };
  }, [assetId]);

  return found && found.assetId === assetId ? found.ids : NO_ANCESTORS;
}

type TreeProps = {
  selectedId: string | undefined;
  ancestorIds: string[];
};

function TreeLevel({ assets, ...tree }: TreeProps & { assets: Asset[] }) {
  return (
    <ul className="flex flex-col">
      {assets.map((asset) => (
        <TreeNode key={asset.asset_id} asset={asset} {...tree} />
      ))}
    </ul>
  );
}

// One asset in the tree: a button that opens its children, if it has any,
// and a link to its page.
function TreeNode({ asset, ...tree }: TreeProps & { asset: Asset }) {
  const selected = asset.asset_id === tree.selectedId;
  const isAncestor = tree.ancestorIds.includes(asset.asset_id);
  const [open, setOpen] = useState(isAncestor);

  // Opens the asset when an asset beneath it becomes the selected one.
  const [wasAncestor, setWasAncestor] = useState(isAncestor);
  if (isAncestor !== wasAncestor) {
    setWasAncestor(isAncestor);
    if (isAncestor) setOpen(true);
  }

  // Brings the selected asset into view, e.g. on opening its page directly.
  const row = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (selected) row.current?.scrollIntoView({ block: "nearest" });
  }, [selected]);

  const Chevron = open ? FaChevronDown : FaChevronRight;
  const TypeIcon = TYPE_ICONS[asset.asset_type]?.Icon;

  return (
    <li>
      <div
        ref={row}
        data-asset-type={asset.asset_type}
        className={`flex items-center rounded pr-2 ${
          selected
            ? "bg-(--asset-tint) ring-1 ring-(--asset-accent) ring-inset"
            : "hover:bg-surface"
        }`}
      >
        {asset.child_count > 0 ? (
          <button
            type="button"
            aria-expanded={open}
            aria-label={`${open ? "Hide" : "Show"} the children of ${asset.asset_id}`}
            onClick={() => setOpen(!open)}
            className="cursor-pointer p-2 text-muted hover:text-accent"
          >
            <Chevron aria-hidden className="size-3" />
          </button>
        ) : (
          <span className="w-7 shrink-0" />
        )}
        <Link
          href={explorerPath(asset.asset_id)}
          aria-current={selected ? "page" : undefined}
          title={`${asset.asset_id} · ${asset.asset_name}`}
          className="flex min-w-0 flex-1 items-center gap-2 py-1 text-sm"
        >
          {TypeIcon && (
            <TypeIcon
              aria-hidden
              className="size-4 shrink-0 text-(--asset-text) opacity-70"
            />
          )}
          <span
            className={`whitespace-nowrap ${selected ? "font-semibold" : "font-medium"}`}
          >
            {asset.asset_id}
          </span>
          <span className="truncate text-muted">{asset.asset_name}</span>
        </Link>
        <StatusIcon
          status={asset.operational_status}
          className="size-3.5 shrink-0"
        />
      </div>
      {open && (
        <div className="ml-3.5 border-l border-border pl-1">
          <TreeChildren assetId={asset.asset_id} {...tree} />
        </div>
      )}
    </li>
  );
}

// Fetches an asset's children when it is first opened.
function TreeChildren({ assetId, ...tree }: TreeProps & { assetId: string }) {
  const details = useLoaded(() => fetchAssetDetails(assetId));

  if (details.state === "loading") {
    return <p className="py-1 pl-7 text-sm text-muted">Loading...</p>;
  }
  if (details.state === "error") {
    return (
      <p className="py-1 pl-7 text-sm font-medium text-red-700">
        {details.message}
      </p>
    );
  }
  return <TreeLevel assets={details.data.children} {...tree} />;
}
