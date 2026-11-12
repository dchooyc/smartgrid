import { useState } from "react";
import {
  commitValidRows,
  type CommitResult,
  type ValidationResult,
} from "../lib/api";
import { download } from "../lib/download";
import ConfirmDownloadDialog from "./confirm-download-dialog";

type CommitState =
  | { state: "idle" }
  | { state: "submitting" }
  | { state: "done"; result: CommitResult };

type CommitRowsProps = {
  result: ValidationResult;
  onReset: () => void; // clears the upload so another file can be chosen
};

export default function CommitRows({ result, onReset }: CommitRowsProps) {
  const [commit, setCommit] = useState<CommitState>({ state: "idle" });
  const [confirming, setConfirming] = useState(false);

  const outcome = commit.state === "done" ? commit.result : null;
  const committed = outcome?.committed === true;
  const hasInvalidRows = result.errors.length > 0;
  const linesWithErrors = new Set(result.errors.map((error) => error.line)).size;
  const hasValidRows = result.totalRows > linesWithErrors;

  async function handleCommit() {
    setCommit({ state: "submitting" });
    setCommit({ state: "done", result: await commitValidRows(result.validFile) });
  }

  function handleUploadAnother() {
    if (hasInvalidRows) setConfirming(true);
    else onReset();
  }

  return (
    <>
      <button
        type="button"
        // Stays disabled after a successful commit; a failed one can be retried.
        disabled={!hasValidRows || commit.state === "submitting" || committed}
        onClick={handleCommit}
        className="rounded-full bg-primary px-6 py-3 text-base font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-primary"
      >
        Commit Valid Rows
      </button>

      {commit.state === "submitting" && (
        <p className="text-muted">Submitting valid rows...</p>
      )}

      {outcome?.committed === false && (
        <div className="flex flex-col gap-2 text-red-700">
          <p className="font-medium">{outcome.message}</p>
          {outcome.errors.length > 0 && (
            <ol className="max-h-60 list-decimal space-y-1 overflow-y-auto rounded-lg border border-border bg-background py-3 pl-10 pr-4 text-left">
              {outcome.errors.map((error, index) => (
                <li key={index}>
                  Line {error.line}: {error.message}
                </li>
              ))}
            </ol>
          )}
        </div>
      )}

      {outcome?.committed === true && (
        <>
          <p className="font-medium text-accent">
            Process completed successfully! {outcome.rowsCommitted} rows
            committed.
          </p>
          <button
            type="button"
            onClick={handleUploadAnother}
            className="cursor-pointer rounded-full border border-foreground bg-background px-6 py-3 text-base font-semibold text-foreground transition-colors hover:bg-surface focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            Upload Another CSV
          </button>
        </>
      )}

      <ConfirmDownloadDialog
        open={confirming}
        onDownload={() => download(result.invalidFile)}
        onConfirm={onReset}
        onClose={() => setConfirming(false)}
      />
    </>
  );
}
