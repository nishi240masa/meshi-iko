import { useNavigate } from "@tanstack/react-router";
import { useReducer, useState } from "react";

import { api } from "../../lib/api";
import { createTimeScale } from "./createTimeScale";
import { TimeSlotCard } from "./TimeSlotCard";
import { initialTimeSlotState, timeSlotReducer } from "./timeSlotReducer";
import styles from "./VoteForm.module.css";
import { voteFormSchema } from "./VoteSchema";

type Props = {
  selectableTimes: string[];
};

export function VoteForm({ selectableTimes }: Props) {
  const navigate = useNavigate();
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
        headers: {
          Authorization: `Bearer 9f2c1d0b8a7e6f5d4c3b2a1908172635`,
        },
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
      navigate({ to: "/result" });
    } catch (error) {
      setErrorMessage("投票に失敗しました");
      console.error("Vote request error:", error);
    }
  };

  return (
    <form onSubmit={handleVote} className={styles.voteForm}>
      <div className={styles.voteStatusContainer}>
        <p>本日の投票</p>
        <div className={styles.radioGroup}>
          <div className={styles.radioItem}>
            <input
              type="radio"
              name="rsvp"
              id="available"
              value="available"
              onChange={() => setStatus("available")}
            />
            <label htmlFor="available">いく</label>
          </div>
          <div className={styles.radioItem}>
            <input
              type="radio"
              name="rsvp"
              id="unavailable"
              value="unavailable"
              onChange={() => {
                dispatch({ type: "RESET" });
                setStatus("unavailable");
              }}
            />
            <label htmlFor="unavailable">いかない</label>
          </div>
        </div>
      </div>
      <div className={styles.timeSlotInfo}>
        <p>参加できる時間帯を選択してください</p>
        <p>
          選択可能な時間帯:<span>{selectableTimes[0]}</span>~
          <span>{selectableTimes[selectableTimes.length - 1]}</span>
        </p>
      </div>
      <div className={styles.timeSlotContainer}>
        {state.slots.length > 0 && (
          <div className={styles.timeSlotList}>
            {state.slots.map((timeSlot, index) => (
              <TimeSlotCard
                key={timeSlot.id}
                slot={timeSlot}
                onChangeTimeSlot={(updatedSlot) =>
                  dispatch({ type: "UPDATE", payload: { ...updatedSlot } })
                }
                onRemoveTimeSlot={() => dispatch({ type: "REMOVE", payload: { id: timeSlot.id } })}
                number={index + 1}
              />
            ))}
          </div>
        )}
      </div>

      {status === "available" && (
        <button
          type="button"
          disabled={status !== "available"}
          className={styles.addButton}
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
          時間帯を追加
        </button>
      )}
      <p>投票結果は後から確認・変更できます</p>
      <button className={styles.submitButton} type="submit">
        投票する
      </button>
      {errorMessage && <p className={styles.errorMessage}>{errorMessage}</p>}
    </form>
  );
}
