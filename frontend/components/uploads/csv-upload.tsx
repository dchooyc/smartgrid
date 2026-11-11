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
        className="text-sm font-medium text-zinc-900 dark:text-zinc-100"
      >
        Upload a CSV file
      </label>

      <input
        id="csv-upload"
        type="file"
        accept=".csv"
        onChange={handleChange}
        className="block w-full cursor-pointer rounded-md border border-zinc-300 text-sm text-zinc-700 file:mr-3 file:cursor-pointer file:border-0 file:bg-zinc-100 file:px-3 file:py-2 file:text-sm file:font-medium file:text-zinc-900 hover:file:bg-zinc-200 dark:border-zinc-700 dark:text-zinc-300 dark:file:bg-zinc-800 dark:file:text-zinc-100 dark:hover:file:bg-zinc-700"
      />

      <button
        type="button"
        disabled={!file}
        onClick={() => file && onValidate?.(file)}
        className="rounded-md bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-700 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
      >
        Validate
      </button>

      <p className="text-sm text-zinc-600 dark:text-zinc-400" aria-live="polite">
        {file ? `Selected file: ${file.name}` : "No file selected (.csv only)"}
      </p>
    </div>
  );
}
