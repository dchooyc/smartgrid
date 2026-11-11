// True if the file is named .csv and its contents parse as CSV.
export async function isCsv(file: File): Promise<boolean> {
  if (!file.name.toLowerCase().endsWith(".csv")) return false;

  try {
    return isCsvContent(await file.text());
  } catch {
    return false; // the file could not be read at all
  }
}

export async function readColumnTitles(file: File): Promise<string[]> {
  // The header row is at the start of the file.
  const start = await file.slice(0, 64 * 1024).text();
  return parseHeaderRow(start);
}

// Well formed: not binary, quotes closed, same column count in every row.
function isCsvContent(text: string): boolean {
  if (text.includes("\0")) return false; // binary data, not text

  let columns: number | null = null; // column count of the first row
  let fields = 1; // fields seen so far in the current row
  let rowHasContent = false;
  let inQuotes = false;
  let fieldStart = true; // at the first character of a field
  let afterQuote = false; // just closed a quoted field

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

function parseHeaderRow(text: string): string[] {
  // Drop Excel's byte-order mark and leading blank lines.
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
