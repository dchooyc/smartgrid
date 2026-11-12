import type { ValidationError } from "../lib/api";

export type BackendStatus =
  | { state: "idle" }
  | { state: "loading" }
  | { state: "success"; rowCount: number }
  | { state: "invalid"; errors: ValidationError[] }
  | { state: "error"; message: string };

type ValidationStatusProps = {
  status: BackendStatus;
};

export default function ValidationStatus({ status }: ValidationStatusProps) {
  // empty:hidden avoids a gap when there is nothing to show.
  return (
    <div className="text-center text-sm empty:hidden" aria-live="polite">
      {status.state === "loading" && (
        <p className="text-muted">Doing further validations on the backend...</p>
      )}
      {status.state === "success" && (
        <p className="font-medium text-accent">
          File has been validated successfully! {status.rowCount} rows read!
        </p>
      )}
      {status.state === "invalid" && (
        <ol className="max-h-80 list-decimal space-y-1 overflow-y-auto pl-8 text-left text-red-700">
          {status.errors.map((error, index) => (
            <li key={index}>
              Line {error.line}: {error.message}
            </li>
          ))}
        </ol>
      )}
      {status.state === "error" && (
        <p className="font-medium text-red-700">{status.message}</p>
      )}
    </div>
  );
}
