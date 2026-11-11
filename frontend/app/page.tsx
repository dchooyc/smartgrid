import CsvUpload from "@/components/uploads/csv-upload";

export default function Home() {
  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-surface font-sans">
      <main className="flex flex-1 w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-background sm:items-start">
        <div className="text-2xl font-semibold text-foreground">{"Smart Grid"}</div>
        <CsvUpload />
      </main>
    </div>
  );
}
