import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/result")({
  component: ResultPage,
});

function ResultPage() {
  return (
    <div className="container">
      <h1>結果確認ページ</h1>
    </div>
  );
}
