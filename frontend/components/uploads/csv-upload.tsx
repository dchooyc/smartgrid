"use client";

import { useState, type ChangeEvent } from "react";

type CsvUploadProps = {
  /** Called with the selected file when Validate is clicked. */
  onValidate?: (file: File) => void;
};

function isCsv(file: File): boolean {
  return file.name.toLowerCase().endsWith(".csv");
}

export default function CsvUpload({ onValidate }: CsvUploadProps) {
  const [file, setFile] = useState<File | null>(null);

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const selected = event.target.files?.[0] ?? null;

    // The `accept` attribute is only a hint to the file picker, so check here too.
    if (selected && !isCsv(selected)) {
      event.target.value = "";
      setFile(null);
      return;
    }

    setFile(selected);
  }

  return (
    <div className="flex w-full max-w-sm flex-col gap-3">
      <label
        htmlFor="csv-upload"
        className="text-sm font-medium text-foreground"
      >
        Upload a CSV file
      </label>

      <input
        id="csv-upload"
        type="file"
        accept=".csv"
        onChange={handleChange}
        className="block w-full cursor-pointer rounded-md border border-border text-sm text-muted file:mr-3 file:cursor-pointer file:border-0 file:bg-surface file:px-3 file:py-2 file:text-sm file:font-medium file:text-accent hover:file:bg-border focus-visible:outline-2 focus-visible:outline-accent"
      />

      <button
        type="button"
        disabled={!file}
        onClick={() => file && onValidate?.(file)}
        className="rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-primary"
      >
        Validate
      </button>

      <p className="text-sm text-muted" aria-live="polite">
        {file ? `Selected file: ${file.name}` : "No file selected (.csv only)"}
      </p>
    </div>
  );
}
