import { createFileRoute } from "@tanstack/react-router";

import { VoteForm } from "../components/vote/VoteForm";

export const Route = createFileRoute("/vote")({
  loader: async () => {
    // 投票可能時間と投票状態の取得
    // YYYY--MM--DD
    const date = new Intl.DateTimeFormat("sv-SE", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      timeZone: "Asia/Tokyo", //JST指定
    }).format(new Date());
    console.log("Current date and time:", date);
    return { date };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { date } = Route.useLoaderData();
  console.log(date);

  // 投票可能時間と投票状態を取得して表示切り替え

  return (
    <div>
      <h1>投票ページ</h1>
      <VoteForm
        selectableTimes={[
          "16:00",
          "16:30",
          "17:00",
          "17:30",
          "18:00",
          "18:30",
          "19:00",
          "19:30",
          "20:00",
          "20:30",
          "21:00",
        ]}
      />
    </div>
  );
}
