"use client";

import { useEffect, useRef, useState } from "react";
import { deleteAllAssets } from "./lib/api";

type Progress =
  | { state: "closed" }
  | { state: "confirming" }
  | { state: "deleting" }
  | { state: "done"; deleted: number }
  | { state: "error"; message: string };

const BUTTON =
  "cursor-pointer rounded-full px-5 py-2 text-sm font-semibold transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50";

// A button that deletes every stored asset, after a warning. It is for
// starting again with an empty database.
export default function DeleteAll() {
  const [progress, setProgress] = useState<Progress>({ state: "closed" });
  const dialogRef = useRef<HTMLDialogElement>(null);
  const cancelButtonRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (progress.state !== "confirming") return;
    dialogRef.current?.showModal();
    cancelButtonRef.current?.focus(); // the safe choice gets the focus
  }, [progress.state]);

  async function handleDelete() {
    setProgress({ state: "deleting" });
    try {
      setProgress({ state: "done", deleted: await deleteAllAssets() });
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setProgress({ state: "error", message });
    }
  }

  function handleClose() {
    // Once the assets are gone, anything this page still shows about them
    // is out of date, so it starts afresh.
    if (progress.state === "done") window.location.reload();
    else setProgress({ state: "closed" });
  }

  const asking = progress.state === "confirming" || progress.state === "deleting";

  return (
    <>
      <button
        type="button"
        onClick={() => setProgress({ state: "confirming" })}
        className={`${BUTTON} border border-red-700 text-red-700 hover:bg-red-700 hover:text-white`}
      >
        Delete all assets
      </button>

      <dialog
        ref={dialogRef}
        onClose={handleClose}
        // Closing with Escape is not offered while the delete is under way.
        onCancel={(event) => progress.state === "deleting" && event.preventDefault()}
        className="m-auto w-full max-w-md rounded-xl border border-border bg-background p-6 text-left text-foreground shadow-xl backdrop:bg-foreground/40"
      >
        {asking && (
          <>
            <p className="text-base font-semibold text-red-700">
              Delete all assets?
            </p>
            <p className="mt-3 text-sm">
              This permanently deletes every asset stored in the database.
              There is no undo and no backup is taken.
            </p>
            <div className="mt-6 flex flex-wrap justify-end gap-3">
              <button
                ref={cancelButtonRef}
                type="button"
                disabled={progress.state === "deleting"}
                onClick={() => dialogRef.current?.close()}
                className={`${BUTTON} border border-foreground text-foreground hover:bg-surface`}
              >
                Cancel
              </button>
              <button
                type="button"
                disabled={progress.state === "deleting"}
                onClick={handleDelete}
                className={`${BUTTON} bg-red-700 text-white hover:bg-red-800`}
              >
                {progress.state === "deleting" ? "Deleting..." : "Yes, delete all"}
              </button>
            </div>
          </>
        )}

        {(progress.state === "done" || progress.state === "error") && (
          <>
            {progress.state === "done" ? (
              <p className="text-base font-medium">
                {progress.deleted === 0
                  ? "There were no assets to delete."
                  : `Deleted ${progress.deleted} assets.`}
              </p>
            ) : (
              <p className="text-base font-medium text-red-700">
                {progress.message}
              </p>
            )}
            <div className="mt-6 flex justify-end">
              <button
                type="button"
                autoFocus
                onClick={() => dialogRef.current?.close()}
                className={`${BUTTON} bg-primary text-primary-foreground hover:bg-primary-hover`}
              >
                OK
              </button>
            </div>
          </>
        )}
      </dialog>
    </>
  );
}
