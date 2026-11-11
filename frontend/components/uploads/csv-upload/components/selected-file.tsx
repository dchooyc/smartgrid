import { FaFileCsv, FaXmark } from "react-icons/fa6";

type SelectedFileProps = {
  file: File;
  onRemove: () => void;
  removeDisabled?: boolean;
};

export default function SelectedFile({
  file,
  onRemove,
  removeDisabled = false,
}: SelectedFileProps) {
  return (
    <div className="flex flex-1 items-center justify-center rounded-lg border-2 border-border bg-background p-6">
      <div className="relative flex max-w-full flex-col items-center gap-3 px-8 pt-6">
        <button
          type="button"
          onClick={onRemove}
          disabled={removeDisabled}
          aria-label={`Remove ${file.name}`}
          className="absolute right-0 top-0 flex size-7 cursor-pointer items-center justify-center rounded-full bg-foreground text-background transition-colors hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50"
        >
          <FaXmark className="size-4" aria-hidden />
        </button>
        <FaFileCsv className="size-20 text-accent" aria-hidden />
        <p className="max-w-full break-all text-center text-sm font-medium text-foreground">
          {file.name}
        </p>
      </div>
    </div>
  );
}
