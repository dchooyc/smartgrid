import { CsvUpload, DeleteAll } from "@/components/uploads";

export default function Home() {
  return (
    <main className="page-column flex flex-1 flex-col py-10 font-sans">
      <div className="flex flex-1 flex-col gap-8">
        <div className="flex items-center justify-between gap-4">
          <div className="text-2xl font-semibold text-foreground">
            {"Smart Grid"}
          </div>
          <DeleteAll />
        </div>
        <CsvUpload />
      </div>
    </main>
  );
}
