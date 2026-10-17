import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/result")({
  component: ResultPage,
});

function ResultPage() {
  return <div className="container">投票結果を表示するページです</div>;
}
