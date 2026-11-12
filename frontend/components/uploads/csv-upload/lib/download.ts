import type { GeneratedFile } from "./api";

// Saves the file through a temporary link, since the content is already in memory.
export function download(file: GeneratedFile) {
  const url = URL.createObjectURL(
    new Blob([file.content], { type: "text/csv;charset=utf-8" }),
  );
  const link = document.createElement("a");
  link.href = url;
  link.download = file.name;
  link.click();
  URL.revokeObjectURL(url);
}
