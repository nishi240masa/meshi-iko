import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/vote")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>投票ページ</div>;
}
