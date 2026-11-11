"use client";

import { useState, type ChangeEvent } from "react";
import { FaCloudArrowUp, FaFileCsv, FaXmark } from "react-icons/fa6";

type CsvUploadProps = {
  /** Called with the selected file when Validate is clicked. */
  onValidate?: (file: File) => void;
};

function isCsv(file: File): boolean {
  return file.name.toLowerCase().endsWith(".csv");
}

export default function CsvUpload({ onValidate }: CsvUploadProps) {
  // Only one file at a time: selecting is only possible while this is null.
  const [file, setFile] = useState<File | null>(null);

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const selected = event.target.files?.[0] ?? null;

    // The `accept` attribute is only a hint to the file picker, so check here too.
    if (selected && !isCsv(selected)) {
      event.target.value = "";
      return;
    }

    setFile(selected);
  }

  return (
    <div className="flex w-full flex-1 flex-col gap-4 rounded-xl border border-border bg-surface p-6 sm:p-10">
      {file ? (
        <div className="flex flex-1 items-center justify-center rounded-lg border-2 border-border bg-background p-6">
          <div className="relative flex max-w-full flex-col items-center gap-3 px-8 pt-6">
            <button
              type="button"
              onClick={() => setFile(null)}
              aria-label={`Remove ${file.name}`}
              className="absolute right-0 top-0 flex size-7 cursor-pointer items-center justify-center rounded-full bg-foreground text-background transition-colors hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            >
              <FaXmark className="size-4" aria-hidden />
            </button>
            <FaFileCsv className="size-20 text-accent" aria-hidden />
            <p className="max-w-full break-all text-center text-sm font-medium text-foreground">
              {file.name}
            </p>
          </div>
        </div>
      ) : (
        <label className="flex flex-1 cursor-pointer flex-col items-center justify-center gap-3 rounded-lg border-2 border-dashed border-border bg-background p-6 text-center transition-colors hover:border-primary-hover focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-accent">
          <FaCloudArrowUp className="size-12 text-accent" aria-hidden />
          <span className="text-base font-medium text-foreground">
            Click to choose a CSV file
          </span>
          <span className="text-sm text-muted">.csv only, one file at a time</span>
          <input
            type="file"
            accept=".csv"
            onChange={handleChange}
            className="sr-only"
          />
        </label>
      )}

      <button
        type="button"
        disabled={!file}
        onClick={() => file && onValidate?.(file)}
        className="rounded-full bg-primary px-6 py-4 text-lg font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-primary"
      >
        Validate
      </button>
    </div>
  );
}
