"use client";

import { useRef, useState, type ChangeEvent, type ReactNode } from "react";
import { FaCloudArrowUp, FaFileCsv, FaXmark } from "react-icons/fa6";

type CsvUploadProps = {
  /** Called with the selected file when Validate is clicked and its columns are correct. */
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

/** Column titles every asset file must have, as in grid_assets.csv. */
const REQUIRED_COLUMNS = [
  "asset_id",
  "parent_asset_id",
  "asset_type",
  "asset_name",
  "operational_status",
  "voltage_kv",
  "rating_kva",
  "manufacturer",
  "model",
  "serial_number",
  "commissioned_date",
];

/** Returns the column titles from the first non-blank row of CSV text. */
function parseHeaderRow(text: string): string[] {
  // Drop the byte-order mark Excel adds, and any leading blank lines.
  const csv = text.replace(/^\uFEFF/, "").replace(/^[\r\n]+/, "");

  const titles: string[] = [];
  let title = "";
  let inQuotes = false;

  for (let i = 0; i < csv.length; i++) {
    const char = csv[i];

    if (inQuotes) {
      if (char !== '"') title += char;
      else if (csv[i + 1] === '"') {
        title += '"'; // escaped quote ("")
        i++;
      } else inQuotes = false;
    } else if (char === '"' && title === "") {
      inQuotes = true;
    } else if (char === ",") {
      titles.push(title.trim());
      title = "";
    } else if (char === "\n" || char === "\r") {
      break; // end of the header row
    } else {
      title += char;
    }
  }

  titles.push(title.trim());
  return titles;
}

export default function CsvUpload({ onValidate }: CsvUploadProps) {
  // Only one file at a time: selecting is only possible while this is null.
  const [file, setFile] = useState<File | null>(null);

  // Messages are shown in a dialog rather than alert(), which cannot format text.
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [message, setMessage] = useState<ReactNode>(null);

  function showMessage(content: ReactNode) {
    setMessage(content);
    dialogRef.current?.showModal();
  }

  async function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const input = event.target;
    const selected = input.files?.[0];
    if (!selected) return;

    if (await isCsv(selected)) {
      setFile(selected);
    } else {
      input.value = ""; // allow picking the same file again
      showMessage(<p>The file uploaded cannot be read as a .csv file!</p>);
    }
  }

  async function handleValidate() {
    if (!file) return;

    // The header row is at the start, so there is no need to read the whole file.
    const titles = parseHeaderRow(await file.slice(0, 64 * 1024).text());

    const missing = REQUIRED_COLUMNS.filter((column) => !titles.includes(column));
    if (missing.length > 0) {
      showMessage(
        <>
          <p>The csv file does not contain these required columns:</p>
          <p className="break-words font-bold">{missing.join(", ")}</p>
          <p>Please upload a complete asset file.</p>
        </>,
      );
      return;
    }

    // Titles that appear more than once, each listed once. Untitled columns
    // are left to the column-count check below.
    const duplicates = [
      ...new Set(
        titles.filter(
          (title, index) => title !== "" && titles.indexOf(title) !== index,
        ),
      ),
    ];
    if (duplicates.length > 0) {
      showMessage(
        <>
          <p>The file contains these duplicate columns:</p>
          <p className="break-words font-bold">{duplicates.join(", ")}</p>
          <p>Please remove the duplicate columns before proceeding.</p>
        </>,
      );
      return;
    }

    // All required columns are present, so any difference in count is extras.
    if (titles.length !== REQUIRED_COLUMNS.length) {
      showMessage(
        <p>
          The csv file must contain exactly {REQUIRED_COLUMNS.length} columns,
          but it has {titles.length}. Please remove the extra columns.
        </p>,
      );
      return;
    }

    onValidate?.(file);
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
        onClick={handleValidate}
        className="rounded-full bg-primary px-6 py-4 text-lg font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-primary"
      >
        Validate
      </button>

      <dialog
        ref={dialogRef}
        className="m-auto w-full max-w-md rounded-xl border border-border bg-background p-6 text-foreground shadow-xl backdrop:bg-foreground/40"
      >
        <div className="flex flex-col gap-3 text-sm">{message}</div>
        <form method="dialog" className="mt-6 flex justify-end">
          <button
            autoFocus
            className="cursor-pointer rounded-full bg-primary px-6 py-2 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            OK
          </button>
        </form>
      </dialog>
    </div>
  );
}
