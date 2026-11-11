"use client";

import { useState, type ChangeEvent } from "react";
import { FaCloudArrowUp, FaFileCsv, FaXmark } from "react-icons/fa6";

type CsvUploadProps = {
  /** Called with the selected file when Validate is clicked. */
  onValidate?: (file: File) => void;
};

/**
 * Parses text as comma-separated values and reports whether it is well formed:
 * not binary, at least one row, quotes closed properly, and the same number of
 * columns in every row.
 */
function isCsvContent(text: string): boolean {
  if (text.includes("\0")) return false; // binary data, not text

  let columns: number | null = null; // column count of the first row
  let fields = 1; // fields seen so far in the current row
  let rowHasContent = false;
  let inQuotes = false;
  let fieldStart = true; // at the first character of a field
  let afterQuote = false; // just closed a quoted field

  // Returns false if the finished row's column count differs from the first row.
  function endRow(): boolean {
    const blank = fields === 1 && !rowHasContent;
    if (!blank) {
      if (columns === null) columns = fields;
      else if (fields !== columns) return false;
    }
    fields = 1;
    rowHasContent = false;
    fieldStart = true;
    afterQuote = false;
    return true;
  }

  for (let i = 0; i < text.length; i++) {
    const char = text[i];

    if (inQuotes) {
      if (char !== '"') continue;
      if (text[i + 1] === '"') i++; // escaped quote ("")
      else {
        inQuotes = false;
        afterQuote = true;
      }
    } else if (char === ",") {
      fields++;
      rowHasContent = true;
      fieldStart = true;
      afterQuote = false;
    } else if (char === "\n" || char === "\r") {
      if (char === "\r" && text[i + 1] === "\n") i++;
      if (!endRow()) return false;
    } else if (afterQuote) {
      return false; // text directly after a closing quote
    } else if (char === '"' && fieldStart) {
      inQuotes = true;
      rowHasContent = true;
      fieldStart = false;
    } else {
      rowHasContent = true;
      fieldStart = false;
    }
  }

  if (inQuotes) return false; // quote never closed
  if (!endRow()) return false;
  return columns !== null; // false when the file has no rows at all
}

/** True if the file is named .csv and its contents can be read as CSV. */
async function isCsv(file: File): Promise<boolean> {
  if (!file.name.toLowerCase().endsWith(".csv")) return false;

  try {
    return isCsvContent(await file.text());
  } catch {
    return false; // the file could not be read at all
  }
}

export default function CsvUpload({ onValidate }: CsvUploadProps) {
  // Only one file at a time: selecting is only possible while this is null.
  const [file, setFile] = useState<File | null>(null);

  async function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const input = event.target;
    const selected = input.files?.[0];
    if (!selected) return;

    if (await isCsv(selected)) {
      setFile(selected);
    } else {
      input.value = ""; // allow picking the same file again
      alert("The file uploaded cannot be read as a .csv file!");
    }
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
