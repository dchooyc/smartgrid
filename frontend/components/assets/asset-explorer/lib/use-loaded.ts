import { useEffect, useState } from "react";

export type Loaded<T> =
  | { state: "loading" }
  | { state: "done"; data: T }
  | { state: "error"; message: string };

// Runs load once, when the component first appears, and tracks its progress.
export function useLoaded<T>(load: () => Promise<T>): Loaded<T> {
  const [loaded, setLoaded] = useState<Loaded<T>>({ state: "loading" });

  useEffect(() => {
    let cancelled = false;
    load().then(
      (data) => !cancelled && setLoaded({ state: "done", data }),
      (error) =>
        !cancelled && setLoaded({ state: "error", message: error.message }),
    );
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps -- load once
  }, []);

  return loaded;
}
