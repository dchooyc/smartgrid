export type BackendStatus =
  | { state: "idle" }
  | { state: "loading" }
  | { state: "success"; rowCount: number }
  | { state: "error"; message: string };

type ValidationStatusProps = {
  status: BackendStatus;
};

export default function ValidationStatus({ status }: ValidationStatusProps) {
  // empty:hidden avoids a gap when there is nothing to show.
  return (
    <p className="text-center text-sm empty:hidden" aria-live="polite">
      {status.state === "loading" && (
        <span className="text-muted">
          Doing further validations on the backend...
        </span>
      )}
      {status.state === "success" && (
        <span className="font-medium text-accent">
          File has been validated successfully! {status.rowCount} rows read!
        </span>
      )}
      {status.state === "error" && (
        <span className="font-medium text-red-700">{status.message}</span>
      )}
    </p>
  );
}
