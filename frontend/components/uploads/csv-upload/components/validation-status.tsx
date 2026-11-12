import type { ValidationResult } from "../lib/api";

export type BackendStatus =
  | { state: "idle" }
  | { state: "loading" }
  | { state: "done"; result: ValidationResult }
  | { state: "error"; message: string };

type ValidationStatusProps = {
  status: BackendStatus;
};

export default function ValidationStatus({ status }: ValidationStatusProps) {
  // empty:hidden avoids a gap when there is nothing to show.
  return (
    <div
      className="flex flex-col gap-3 text-center text-sm empty:hidden"
      aria-live="polite"
    >
      {status.state === "loading" && (
        <p className="text-muted">Doing further validations on the backend...</p>
      )}
      {status.state === "done" && <ValidationReport result={status.result} />}
      {status.state === "error" && (
        <p className="font-medium text-red-700">{status.message}</p>
      )}
    </div>
  );
}

function ValidationReport({ result }: { result: ValidationResult }) {
  const { totalRows, errors, warnings } = result;

  return (
    <>
      {errors.length === 0 && (
        <p className="font-medium text-accent">
          File has been validated successfully! {totalRows} rows read!
        </p>
      )}
      {errors.length + warnings.length > 0 && (
        <ol className="max-h-80 list-decimal space-y-1 overflow-y-auto pl-8 text-left">
          {errors.map((error, index) => (
            <li key={`error-${index}`} className="text-red-700">
              Line {error.line}: {error.message}
            </li>
          ))}
          {warnings.map((warning, index) => (
            <li key={`warning-${index}`} className="text-amber-600">
              Line {warning.line}: Warning: {warning.message}
            </li>
          ))}
        </ol>
      )}
    </>
  );
}
