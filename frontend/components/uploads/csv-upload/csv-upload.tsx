"use client";

import { useState } from "react";
import {
  FilePicker,
  MessageDialog,
  SelectedFile,
  ValidationStatus,
} from "./components";
import type { Message } from "./components/message-dialog";
import type { BackendStatus } from "./components/validation-status";
import { validateOnBackend } from "./lib/api";
import { findColumnProblem } from "./lib/columns";
import { isCsv, readColumnTitles } from "./lib/csv";

export default function CsvUpload() {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<BackendStatus>({ state: "idle" });
  const [message, setMessage] = useState<Message>(null);
  const loading = status.state === "loading";

  async function handleSelect(selected: File) {
    if (!(await isCsv(selected))) {
      setMessage("The file uploaded cannot be read as a .csv file!");
      return;
    }

    setFile(selected);
    setStatus({ state: "idle" });
  }

  function handleRemove() {
    setFile(null);
    setStatus({ state: "idle" });
  }

  async function handleValidate() {
    if (!file || loading) return;
    setStatus({ state: "idle" });

    const problem = findColumnProblem(await readColumnTitles(file));
    if (problem) {
      setMessage(problem);
      return;
    }

    setStatus({ state: "loading" });
    try {
      const result = await validateOnBackend(file);
      setStatus({ state: "done", result });
    } catch (error) {
      const reason = error instanceof Error ? error.message : String(error);
      setStatus({ state: "error", message: reason });
    }
  }

  return (
    <div className="flex w-full flex-1 flex-col gap-4 rounded-xl border border-border bg-surface p-6 sm:p-10">
      {file ? (
        <SelectedFile
          file={file}
          onRemove={handleRemove}
          removeDisabled={loading}
        />
      ) : (
        <FilePicker onSelect={handleSelect} />
      )}

      <button
        type="button"
        disabled={!file || loading}
        onClick={handleValidate}
        className="rounded-full bg-primary px-6 py-4 text-lg font-semibold text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-primary"
      >
        Validate
      </button>

      <ValidationStatus status={status} onReset={handleRemove} />

      <MessageDialog message={message} onClose={() => setMessage(null)} />
    </div>
  );
}
