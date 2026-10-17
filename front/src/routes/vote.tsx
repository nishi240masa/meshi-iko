import { createFileRoute } from "@tanstack/react-router";

import { createTimeScale } from "../components/vote/createTimeScale";
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
    return { date };
  },

  component: RouteComponent,
});

function RouteComponent() {
  const { date } = Route.useLoaderData();
  const selectableTimes = createTimeScale("16:00", "21:00"); // 16:00から21:00までの30分刻みの時間を生成
  console.log(date, selectableTimes);
  // 投票可能時間と投票状態を取得して表示切り替え

  return (
    <div className="container">
      <h1>投票ページ</h1>
      <VoteForm selectableTimes={selectableTimes} />
    </div>
  );
}
