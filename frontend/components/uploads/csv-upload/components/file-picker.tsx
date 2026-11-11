import type { ChangeEvent } from "react";
import { FaCloudArrowUp } from "react-icons/fa6";

type FilePickerProps = {
  onSelect: (file: File) => void;
};

export default function FilePicker({ onSelect }: FilePickerProps) {
  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = ""; // allow picking the same file again
    if (file) onSelect(file);
  }

  return (
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
  );
}
