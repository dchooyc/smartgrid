import { useState } from "react";
import type { ValidationIssue, ValidationResult } from "../lib/api";

type Issue = ValidationIssue & { isWarning: boolean };
type View = "type" | "line";

const ERROR_COLOR = "text-red-700";
const WARNING_COLOR = "text-amber-600";

export default function ValidationReport({ result }: { result: ValidationResult }) {
  const [view, setView] = useState<View>("type");

  const issues: Issue[] = [
    ...result.errors.map((error) => ({ ...error, isWarning: false })),
    ...result.warnings.map((warning) => ({ ...warning, isWarning: true })),
  ];

  return (
    <>
      {result.errors.length === 0 && (
        <p className="font-medium text-accent">
          File has been validated successfully! {result.totalRows} rows read!
        </p>
      )}

      {issues.length > 0 && (
        <div className="flex flex-col gap-3 text-left">
          <Summary result={result} issues={issues} />

          <div className="flex gap-2">
            <ViewButton selected={view === "type"} onClick={() => setView("type")}>
              By error type
            </ViewButton>
            <ViewButton selected={view === "line"} onClick={() => setView("line")}>
              By line
            </ViewButton>
          </div>

          <div className="flex max-h-96 flex-col gap-4 overflow-y-auto rounded-lg border border-border bg-background p-4">
            {view === "type" ? (
              <IssuesByType issues={issues} />
            ) : (
              <IssuesByLine issues={issues} />
            )}
          </div>
        </div>
      )}
    </>
  );
}

function Summary({ result, issues }: { result: ValidationResult; issues: Issue[] }) {
  const linesWithErrors = new Set(result.errors.map((error) => error.line)).size;
  const linesToCommit = result.totalRows - linesWithErrors;

  return (
    <div className="flex flex-col gap-2">
      <p className="font-medium text-accent">
        {linesToCommit} lines can be committed.
      </p>
      <p className={`font-medium ${ERROR_COLOR}`}>
        {linesWithErrors} lines have errors.
      </p>

      <details className="rounded-lg border border-border bg-background">
        <summary className="cursor-pointer px-4 py-2 font-medium text-foreground">
          View error counts
        </summary>
        <ul className="border-t border-border px-4 py-2">
          {groupBy(issues, (issue) => issue.type).map(([type, group]) => (
            <li key={type} className={`flex justify-between gap-4 py-1 ${colorOf(group)}`}>
              <span>{labelOf(type, group)}</span>
              <span className="font-semibold">{group.length}</span>
            </li>
          ))}
        </ul>
      </details>
    </div>
  );
}

function IssuesByType({ issues }: { issues: Issue[] }) {
  return groupBy(issues, (issue) => issue.type).map(([type, group]) => (
    <IssueGroup key={type} title={labelOf(type, group)} issues={group}>
      {(issue) => `Line ${issue.line}: ${issue.message}`}
    </IssueGroup>
  ));
}

function IssuesByLine({ issues }: { issues: Issue[] }) {
  const sorted = [...issues].sort((a, b) => a.line - b.line);

  return groupBy(sorted, (issue) => issue.line).map(([line, group]) => (
    <IssueGroup key={line} title={`Line ${line}`} issues={group}>
      {(issue) => (issue.isWarning ? `Warning: ${issue.message}` : issue.message)}
    </IssueGroup>
  ));
}

type IssueGroupProps = {
  title: string;
  issues: Issue[];
  children: (issue: Issue) => string; // the text shown for each issue
};

function IssueGroup({ title, issues, children }: IssueGroupProps) {
  return (
    <section>
      <h3 className={`font-semibold ${colorOf(issues)}`}>
        {title} ({issues.length})
      </h3>
      <ol className="mt-1 list-decimal space-y-1 pl-8">
        {issues.map((issue, index) => (
          <li key={index} className={issue.isWarning ? WARNING_COLOR : ERROR_COLOR}>
            {children(issue)}
          </li>
        ))}
      </ol>
    </section>
  );
}

type ViewButtonProps = {
  selected: boolean;
  onClick: () => void;
  children: string;
};

function ViewButton({ selected, onClick, children }: ViewButtonProps) {
  const colors = selected
    ? "border-foreground bg-foreground text-background"
    : "border-border bg-background text-foreground hover:border-foreground";

  return (
    <button
      type="button"
      aria-pressed={selected}
      onClick={onClick}
      className={`cursor-pointer rounded-full border px-4 py-1.5 font-medium transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent ${colors}`}
    >
      {children}
    </button>
  );
}

// Groups items by key, keeping the order in which each key first appears.
function groupBy<Key>(issues: Issue[], keyOf: (issue: Issue) => Key): [Key, Issue[]][] {
  const groups = new Map<Key, Issue[]>();
  for (const issue of issues) {
    const group = groups.get(keyOf(issue)) ?? [];
    group.push(issue);
    groups.set(keyOf(issue), group);
  }
  return [...groups];
}

// A group is amber only if everything in it is a warning.
function colorOf(issues: Issue[]): string {
  return issues.every((issue) => issue.isWarning) ? WARNING_COLOR : ERROR_COLOR;
}

function labelOf(type: string, issues: Issue[]): string {
  return issues.every((issue) => issue.isWarning) ? `Warning: ${type}` : type;
}
