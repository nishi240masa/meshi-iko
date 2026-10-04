import { useReducer, useState } from "react";

import { api } from "../../lib/api";
import { createTimeScale } from "./createTimeScale";
import { TimeSlotCard } from "./TimeSlotCard";
import { initialTimeSlotState, timeSlotReducer } from "./timeSlotReducer";
import { voteFormSchema } from "./VoteSchema";

type Props = {
  selectableTimes: string[];
};

export function VoteForm({ selectableTimes }: Props) {
  const [status, setStatus] = useState<"available" | "unavailable" | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [state, dispatch] = useReducer(timeSlotReducer, initialTimeSlotState());

  const handleVote = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    console.log("Form data:", Object.fromEntries(formData.entries()));

    const result = voteFormSchema.safeParse({
      status: status,
      timeSlots: state.slots,
    });

    if (!result.success) {
      setErrorMessage(result.error.issues[0].message ?? "不正な入力です");
      console.error("Validation error:", result.error);
      return;
    }

    setErrorMessage(null); // バリデーションエラーがない場合はエラーメッセージをクリア

    // 各スロットの開始時間と終了時間を30分単位で分割して配列に格納
    // すべてのスロットの一次元の配列に結合
    const slotSplits: string[] =
      status === "available"
        ? state.slots.map((slot) => createTimeScale(slot.startTime, slot.endTime)).flat()
        : [];

    console.log(slotSplits);

    // 投票結果を送信
    try {
      const { data, error } = await api.PUT("/answers/me", {
        body: {
          status: result.data.status,
          timeSlots: slotSplits,
        },
      });

      if (error) {
        setErrorMessage(error.message ?? "投票に失敗しました");
        console.error("Vote error:", error);
        return;
      }

      console.log("Vote response:", data);
    } catch (error) {
      setErrorMessage("投票に失敗しました");
      console.error("Vote request error:", error);
    }
  };

  return (
    <form onSubmit={handleVote}>
      <div>
        <label htmlFor="vote">いくいかない</label>
        <div>
          <div>
            <label htmlFor="available">いく</label>
            <input
              type="radio"
              name="rsvp"
              id="available"
              value="available"
              onChange={() => setStatus("available")}
            />
          </div>
          <div>
            <label htmlFor="unavailable">いかない</label>
            <input
              type="radio"
              name="rsvp"
              id="unavailable"
              value="unavailable"
              onChange={() => setStatus("unavailable")}
            />
          </div>
        </div>
      </div>
      {status === "available" && (
        <div>
          <p>時間帯を設定</p>
          <p>
            選択可能な時間帯<span>{selectableTimes[0]}</span>~
            <span>{selectableTimes[selectableTimes.length - 1]}</span>
          </p>
          <div>
            <p>時間帯を追加</p>
            <button
              type="button"
              onClick={() =>
                dispatch({
                  type: "ADD",
                  payload: {
                    startTime: selectableTimes[0],
                    endTime: selectableTimes[1] ?? selectableTimes[0],
                  },
                })
              }
            >
              追加
            </button>
            {state.slots.map((timeSlot, index) => (
              <div key={timeSlot.id}>
                <TimeSlotCard
                  slot={timeSlot}
                  onChangeTimeSlot={(updatedSlot) =>
                    dispatch({ type: "UPDATE", payload: { ...updatedSlot } })
                  }
                  onRemoveTimeSlot={() =>
                    dispatch({ type: "REMOVE", payload: { id: timeSlot.id } })
                  }
                  number={index + 1}
                />
              </div>
            ))}
          </div>
        </div>
      )}
      <button
        type="submit"
        // 投票可能時間を設定していないか、いくいかないの選択がされていない場合はSubmitボタンを無効化
        disabled={(state.slots.length === 0 && status === "available") || status === null}
      >
        Submit
      </button>
      {errorMessage && <p>{errorMessage}</p>}
    </form>
  );
}
