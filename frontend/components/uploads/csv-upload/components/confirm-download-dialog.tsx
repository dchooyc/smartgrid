import { useEffect, useRef } from "react";

type ConfirmDownloadDialogProps = {
  open: boolean;
  onDownload: () => void; // "No, Download": the dialog stays open afterwards
  onConfirm: () => void; // "Yes, Upload New CSV"
  onClose: () => void;
};

// Asks whether the invalid rows were downloaded before they are cleared.
export default function ConfirmDownloadDialog({
  open,
  onDownload,
  onConfirm,
  onClose,
}: ConfirmDownloadDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const downloadButtonRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (!open) return;
    dialogRef.current?.showModal();
    downloadButtonRef.current?.focus(); // the safe choice gets the focus
  }, [open]);

  return (
    <dialog
      ref={dialogRef}
      onClose={onClose}
      className="m-auto w-full max-w-md rounded-xl border border-border bg-background p-6 text-left text-foreground shadow-xl backdrop:bg-foreground/40"
    >
      <p className="text-base font-medium">
        Have you downloaded the invalid CSV rows?
      </p>
      <div className="mt-6 flex flex-wrap justify-end gap-3">
        <button
          ref={downloadButtonRef}
          type="button"
          onClick={onDownload}
          className="cursor-pointer rounded-full bg-amber-400 px-5 py-2 text-sm font-semibold text-foreground transition-colors hover:bg-amber-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          No, Download
        </button>
        <button
          type="button"
          onClick={() => {
            dialogRef.current?.close();
            onConfirm();
          }}
          className="cursor-pointer rounded-full bg-green-700 px-5 py-2 text-sm font-semibold text-white transition-colors hover:bg-green-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          Yes, Upload New CSV
        </button>
      </div>
    </dialog>
  );
}
