import { useEffect, useRef } from "react";
import type { ColumnProblem } from "../lib/columns";

export type Message = string | ColumnProblem | null;

type MessageDialogProps = {
  message: Message;
  onClose: () => void;
};

// Used instead of alert(), which cannot format text. Opens when message is set.
export default function MessageDialog({ message, onClose }: MessageDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    if (message) dialogRef.current?.showModal();
  }, [message]);

  return (
    <dialog
      ref={dialogRef}
      onClose={onClose}
      className="m-auto w-full max-w-md rounded-xl border border-border bg-background p-6 text-foreground shadow-xl backdrop:bg-foreground/40"
    >
      <div className="flex flex-col gap-3 text-sm">
        {typeof message === "string" && <p>{message}</p>}
        {message && typeof message === "object" && (
          <>
            <p>{message.statement}</p>
            {message.columns.length > 0 && (
              <p className="break-words font-bold">
                {message.columns.join(", ")}
              </p>
            )}
            <p>{message.request}</p>
          </>
        )}
      </div>
      <form method="dialog" className="mt-6 flex justify-end">
        <button
          autoFocus
          className="cursor-pointer rounded-full bg-primary px-6 py-2 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          OK
        </button>
      </form>
    </dialog>
  );
}
