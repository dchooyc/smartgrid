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

export type ColumnProblem = {
  statement: string;
  columns: string[];
  request: string;
};

// Returns the first problem found, or null if the columns are correct.
export function findColumnProblem(titles: string[]): ColumnProblem | null {
  const missing = REQUIRED_COLUMNS.filter((column) => !titles.includes(column));
  if (missing.length > 0) {
    return {
      statement: "The csv file does not contain these required columns:",
      columns: missing,
      request: "Please upload a complete asset file.",
    };
  }

  const duplicates = findDuplicates(titles);
  if (duplicates.length > 0) {
    return {
      statement: "The file contains these duplicate columns:",
      columns: duplicates,
      request: "Please remove the duplicate columns before proceeding.",
    };
  }

  if (titles.length !== REQUIRED_COLUMNS.length) {
    return {
      statement: `The csv file must contain exactly ${REQUIRED_COLUMNS.length} columns, but it has ${titles.length}.`,
      columns: [],
      request: "Please remove the extra columns.",
    };
  }

  return null;
}

// Untitled columns are ignored; the column-count check catches those.
function findDuplicates(titles: string[]): string[] {
  const repeated = titles.filter(
    (title, index) => title !== "" && titles.indexOf(title) !== index,
  );
  return [...new Set(repeated)];
}
