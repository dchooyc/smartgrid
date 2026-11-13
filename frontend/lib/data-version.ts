// The backend lets the browser keep asset responses for a while. The data
// version is added to every asset request, so changing it makes the browser
// ask the backend afresh instead of reusing what it kept.

const STORAGE_KEY = "asset-data-version";

// Used when the browser has no storage: a new version on every page load,
// which is never out of date.
const fallback = String(Date.now());

export function dataVersion(): string {
  try {
    return localStorage.getItem(STORAGE_KEY) ?? "0";
  } catch {
    return fallback;
  }
}

// Call after the stored assets change, e.g. when rows are committed.
export function bumpDataVersion() {
  try {
    localStorage.setItem(STORAGE_KEY, String(Date.now()));
  } catch {
    // Without storage, dataVersion already changes on every page load.
  }
}
