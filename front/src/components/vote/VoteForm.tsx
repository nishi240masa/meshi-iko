import { useReducer, useState } from "react";

import { TimeSlotCard } from "./TimeSlotCard";
import { initialTimeSlotState, timeSlotReducer } from "./timeSlotReducer";

type Props = {
  selectableTimes: string[];
};

export function VoteForm({ selectableTimes }: Props) {
  const [available, setAvailable] = useState<boolean | null>(null);
  const [state, dispatch] = useReducer(timeSlotReducer, initialTimeSlotState());

  const handleVote = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    console.log("Form data:", Object.fromEntries(formData.entries()));

    // 投票結果を送信前に一次元の配列に変換する
    const voteTimeSlots = state.slots.flatMap(({ startTime, endTime }) => [startTime, endTime]);

    console.log(voteTimeSlots);
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
              onChange={() => setAvailable(true)}
            />
          </div>
          <div>
            <label htmlFor="unavailable">いかない</label>
            <input
              type="radio"
              name="rsvp"
              id="unavailable"
              value="unavailable"
              onChange={() => setAvailable(false)}
            />
          </div>
        </div>
      </div>
      {available && (
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
        disabled={(state.slots.length === 0 && available === true) || available === null}
      >
        Submit
      </button>
    </form>
  );
}
