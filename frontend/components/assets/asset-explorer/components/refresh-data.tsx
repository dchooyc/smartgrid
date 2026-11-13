import { bumpDataVersion } from "@/lib/data-version";

// Fetches everything afresh, in case the browser is showing assets it kept
// from before they changed.
export default function RefreshData() {
  function handleRefresh() {
    bumpDataVersion();
    window.location.reload();
  }

  return (
    <button
      type="button"
      onClick={handleRefresh}
      title="Fetch the latest assets from the server"
      className="cursor-pointer text-sm text-muted/85 underline hover:text-accent"
    >
      Refresh data
    </button>
  );
}
