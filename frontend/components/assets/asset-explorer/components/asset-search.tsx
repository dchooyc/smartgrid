"use client";

import { useRouter } from "next/navigation";
import { useEffect, useId, useRef, useState, type KeyboardEvent } from "react";
import { FaMagnifyingGlass } from "react-icons/fa6";
import { searchAssets, type Asset } from "../lib/api";
import { explorerPath } from "../lib/fields";
import { TYPE_ICONS } from "../lib/type-icons";
import StatusIcon from "./status-icon";

// How long typing must pause before the backend is asked.
const TYPING_PAUSE = 200;

// What the backend found for a term, or why it could not search.
type Found = {
  term: string;
  assets: Asset[];
  more: boolean;
  error: string | null;
};

// The explorer's search bar. Typing lists the assets whose ID or name
// matches, and choosing one opens it in the hierarchy.
export default function AssetSearch() {
  const router = useRouter();
  const input = useRef<HTMLInputElement>(null);
  const listId = useId();
  const [query, setQuery] = useState("");
  const [found, setFound] = useState<Found | null>(null);
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState(0); // the match that Enter opens

  const term = query.trim();
  // The last matches stay up while the next ones are fetched.
  const shown = term && open ? found : null;

  // Searches once typing pauses.
  useEffect(() => {
    if (!term) return;

    let cancelled = false;
    const timer = setTimeout(() => {
      searchAssets(term).then(
        (result) => {
          if (cancelled) return;
          setFound({ term, ...result, error: null });
          setActive(0);
        },
        (error) =>
          !cancelled &&
          setFound({ term, assets: [], more: false, error: error.message }),
      );
    }, TYPING_PAUSE);
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [term]);

  // Pressing "/" anywhere on the page moves to the search bar, unless the
  // key is being typed into something.
  useEffect(() => {
    function handleKey(event: globalThis.KeyboardEvent) {
      if (event.key !== "/" || event.ctrlKey || event.metaKey || event.altKey) {
        return;
      }
      const target = event.target as HTMLElement;
      if (target.closest("input, textarea, select, [contenteditable]")) return;

      event.preventDefault();
      input.current?.focus();
      input.current?.select();
    }

    document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  }, []);

  function openAsset(asset: Asset) {
    setOpen(false);
    input.current?.blur();
    router.push(explorerPath(asset.asset_id));
  }

  function handleKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    const assets = shown?.assets ?? [];

    if (event.key === "Escape") {
      // First closes the matches, then leaves the search bar.
      if (shown) setOpen(false);
      else input.current?.blur();
    } else if (event.key === "ArrowDown" && assets.length > 0) {
      event.preventDefault();
      setActive((active + 1) % assets.length);
    } else if (event.key === "ArrowUp" && assets.length > 0) {
      event.preventDefault();
      setActive((active - 1 + assets.length) % assets.length);
    } else if (event.key === "Enter" && assets[active]) {
      openAsset(assets[active]);
    } else if (event.key === "ArrowDown") {
      setOpen(true);
    }
  }

  return (
    <div
      role="search"
      className="flex flex-col gap-2 rounded-xl border border-border bg-surface p-4"
    >
      <div className="relative">
        <FaMagnifyingGlass
          aria-hidden
          className="pointer-events-none absolute left-5 top-1/2 size-4 -translate-y-1/2 text-accent"
        />
        <input
          ref={input}
          type="text"
          role="combobox"
          aria-label="Search assets by ID or name"
          aria-expanded={shown !== null}
          aria-controls={listId}
          aria-activedescendant={
            shown?.assets[active] ? `${listId}-${active}` : undefined
          }
          autoComplete="off"
          spellCheck={false}
          enterKeyHint="search"
          placeholder="Search assets by ID or name"
          value={query}
          onChange={(event) => {
            setQuery(event.target.value);
            setOpen(true);
          }}
          onFocus={() => setOpen(true)}
          onBlur={() => setOpen(false)}
          onKeyDown={handleKeyDown}
          className="w-full rounded-full border-2 border-border bg-background py-2.5 pl-12 pr-5 text-base text-foreground transition-colors placeholder:text-muted/70 hover:border-primary-hover focus:border-primary-hover focus:outline-2 focus:outline-offset-2 focus:outline-accent"
        />

        {shown && (
          <div
            // Keeps the search bar focused while a match is clicked.
            onMouseDown={(event) => event.preventDefault()}
            className="fade-in absolute inset-x-0 top-full z-10 mt-2 overflow-hidden rounded-xl border border-border bg-background shadow-lg"
          >
            <ul
              id={listId}
              role="listbox"
              aria-label="Matching assets"
              className="max-h-96 overflow-y-auto py-2 [scrollbar-width:thin]"
            >
              {shown.assets.map((asset, index) => (
                <Match
                  key={asset.asset_id}
                  id={`${listId}-${index}`}
                  asset={asset}
                  active={index === active}
                  onHover={() => setActive(index)}
                  onChoose={() => openAsset(asset)}
                />
              ))}
            </ul>
            {shown.error ? (
              <p className="px-5 pb-4 pt-2 text-sm font-medium text-red-700">
                {shown.error}
              </p>
            ) : shown.assets.length === 0 ? (
              <p className="px-5 pb-4 pt-2 text-sm text-muted">
                No assets match &ldquo;{shown.term}&rdquo;.
              </p>
            ) : (
              shown.more && (
                <p className="border-t border-border px-5 py-3 text-sm text-muted">
                  Showing the first {shown.assets.length} matches. Keep typing
                  to narrow them down.
                </p>
              )
            )}
          </div>
        )}
      </div>

      <p className="pl-5 text-xs text-muted/85">
        Tip: press{" "}
        <kbd className="rounded border border-border bg-background px-1.5 py-0.5 font-mono">
          /
        </kbd>{" "}
        anywhere on this page to search.
      </p>
    </div>
  );
}

type MatchProps = {
  id: string;
  asset: Asset;
  active: boolean;
  onHover: () => void;
  onChoose: () => void;
};

// One matching asset: its type, ID, name and status.
function Match({ id, asset, active, onHover, onChoose }: MatchProps) {
  const TypeIcon = TYPE_ICONS[asset.asset_type]?.Icon;
  const row = useRef<HTMLLIElement>(null);

  // Keeps the match that the arrow keys moved to in view.
  useEffect(() => {
    if (active) row.current?.scrollIntoView({ block: "nearest" });
  }, [active]);

  return (
    <li
      ref={row}
      id={id}
      role="option"
      aria-selected={active}
      data-asset-type={asset.asset_type}
      onMouseMove={onHover}
      onClick={onChoose}
      className={`flex cursor-pointer items-center gap-3 px-5 py-2 ${active ? "bg-surface" : ""}`}
    >
      {TypeIcon && (
        <TypeIcon
          aria-hidden
          className="size-4 shrink-0 text-(--asset-text) opacity-70"
        />
      )}
      <span className="whitespace-nowrap font-semibold">{asset.asset_id}</span>
      <span className="min-w-0 flex-1 truncate text-muted">
        {asset.asset_name}
      </span>
      <StatusIcon
        status={asset.operational_status}
        className="size-4 shrink-0"
      />
    </li>
  );
}
