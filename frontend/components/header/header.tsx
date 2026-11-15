"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Suspense } from "react";
import {
  ALL_ASSETS,
  EXPLORER_PATH,
  RefreshData,
  TYPE_LISTINGS,
} from "@/components/assets";

const UPLOAD = { path: "/", title: "Upload CSV" };

// The pages that show stored assets: the explorer, then the lists of all
// assets and of each type.
const ASSET_PAGES = [
  { path: EXPLORER_PATH, title: "Explorer" },
  ALL_ASSETS,
  ...TYPE_LISTINGS,
];

// The bar at the top of every page, for moving between the pages.
export default function Header() {
  // The address is only known per request for some pages. Until it is, the
  // bar shows with no page marked as the current one.
  return (
    <Suspense fallback={<HeaderBar pathname="" />}>
      <CurrentHeader />
    </Suspense>
  );
}

function CurrentHeader() {
  return <HeaderBar pathname={usePathname()} />;
}

function HeaderBar({ pathname }: { pathname: string }) {
  return (
    <header className="border-b border-border bg-background font-sans">
      <div className="mx-auto flex w-full max-w-6xl flex-wrap items-center gap-x-6 gap-y-2 px-6 py-3 sm:px-16">
        <nav
          aria-label="Asset pages"
          className="flex flex-wrap items-center gap-x-5 gap-y-1"
        >
          {ASSET_PAGES.map((page) => (
            <HeaderLink key={page.path} page={page} pathname={pathname} />
          ))}
        </nav>
        <div className="ml-auto flex items-center gap-5">
          <HeaderLink page={UPLOAD} pathname={pathname} />
          <RefreshData />
        </div>
      </div>
    </header>
  );
}

type HeaderLinkProps = {
  page: { path: string; title: string };
  pathname: string;
};

function HeaderLink({ page, pathname }: HeaderLinkProps) {
  // A page is also the current one while a page beneath it is open, e.g. the
  // explorer while it shows one asset.
  const current =
    pathname === page.path ||
    (page.path !== "/" && pathname.startsWith(page.path + "/"));

  return (
    <Link
      href={page.path}
      aria-current={current ? "page" : undefined}
      className={
        current
          ? "font-semibold text-foreground"
          : "text-muted hover:text-accent"
      }
    >
      {page.title}
    </Link>
  );
}
