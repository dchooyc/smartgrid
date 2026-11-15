import { CsvUpload } from "@/components/uploads";

export default function Home() {
  return (
    <main className="page-column flex flex-1 flex-col py-10 font-sans">
      <div className="flex flex-col gap-8">
        <div className="text-2xl font-semibold text-foreground">
          {"Smart Grid"}
        </div>
        <CsvUpload />
      </div>
    </main>
  );
}
