import { FaCircleCheck } from "react-icons/fa6";
import type { GeneratedFile } from "../lib/api";

type GeneratedFilesProps = {
  validFile: GeneratedFile;
  invalidFile?: GeneratedFile; // left out when no rows have errors
};

export default function GeneratedFiles({
  validFile,
  invalidFile,
}: GeneratedFilesProps) {
  return (
    <ul className="flex flex-col gap-2 text-left">
      <FileRow file={validFile} color="text-green-800">
        <FaCircleCheck className="size-4 shrink-0" aria-hidden />
        <span className="break-all">{validFile.name}</span>
      </FileRow>
      {invalidFile && (
        <FileRow file={invalidFile} color="text-orange-600">
          <span className="break-all">{invalidFile.name} (To be fixed)</span>
        </FileRow>
      )}
    </ul>
  );
}

type FileRowProps = {
  file: GeneratedFile;
  color: string;
  children: React.ReactNode;
};

function FileRow({ file, color, children }: FileRowProps) {
  return (
    <li className="flex items-center justify-between gap-4 rounded-lg border border-border bg-background px-4 py-2">
      <span className={`flex items-center gap-2 font-medium ${color}`}>
        {children}
      </span>
      <button
        type="button"
        onClick={() => download(file)}
        className="cursor-pointer rounded-full bg-primary px-4 py-1.5 font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
      >
        Download
      </button>
    </li>
  );
}

// Saves the file through a temporary link, since the content is already in memory.
function download(file: GeneratedFile) {
  const url = URL.createObjectURL(
    new Blob([file.content], { type: "text/csv;charset=utf-8" }),
  );
  const link = document.createElement("a");
  link.href = url;
  link.download = file.name;
  link.click();
  URL.revokeObjectURL(url);
}
