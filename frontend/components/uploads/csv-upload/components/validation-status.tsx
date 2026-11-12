import type { ValidationResult } from "../lib/api";
import CommitRows from "./commit-rows";
import GeneratedFiles from "./generated-files";
import ValidationReport from "./validation-report";

export type BackendStatus =
  | { state: "idle" }
  | { state: "loading" }
  | { state: "done"; result: ValidationResult }
  | { state: "error"; message: string };

type ValidationStatusProps = {
  status: BackendStatus;
  onReset: () => void; // clears the upload so another file can be chosen
};

export default function ValidationStatus({
  status,
  onReset,
}: ValidationStatusProps) {
  // empty:hidden avoids a gap when there is nothing to show.
  return (
    <div
      className="flex flex-col gap-3 text-center text-sm empty:hidden"
      aria-live="polite"
    >
      {status.state === "loading" && (
        <p className="text-muted">Doing further validations on the backend...</p>
      )}
      {status.state === "done" && (
        <>
          <ValidationReport result={status.result} />
          <GeneratedFiles
            validFile={status.result.validFile}
            invalidFile={
              status.result.errors.length > 0
                ? status.result.invalidFile
                : undefined
            }
          />
          <CommitRows result={status.result} onReset={onReset} />
        </>
      )}
      {status.state === "error" && (
        <p className="font-medium text-red-700">{status.message}</p>
      )}
    </div>
  );
}
