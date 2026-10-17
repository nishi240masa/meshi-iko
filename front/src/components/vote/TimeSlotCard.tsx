import { useId } from "react";

import { createTimeScale } from "./createTimeScale";
import styles from "./TimeSlotCard.module.css";
import { TimeSlotData } from "./timeSlotReducer";

type Props = {
  slot: TimeSlotData;
  onChangeTimeSlot: (slot: TimeSlotData) => void;
  onRemoveTimeSlot: (slotId: string) => void;

  number: number;
};

export function TimeSlotCard({ slot, onChangeTimeSlot, onRemoveTimeSlot, number }: Props) {
  const id = useId(); // 複数のTimeSlotコンポーネントがある場合に、input要素のidが重複しないようにするためにuseIdを使用
  const selectableTimeRange = createTimeScale("16:00", "21:00"); // 選択可能な時間帯の範囲を作成

  return (
    <div className={styles.timeSlotCard}>
      <p>時間帯{number}</p>
      <label htmlFor={`${id}-startTime`} className={styles.visuallyHidden}>
        開始時間
      </label>
      <select
        name={`timeSlot-${slot.id}-startTime`}
        id={`${id}-startTime`}
        value={slot.startTime}
        className={styles.timeSelect}
        onChange={(e) => onChangeTimeSlot({ ...slot, startTime: e.target.value })}
      >
        {selectableTimeRange.map((range) => {
          return (
            <option key={range} value={range}>
              {range}
            </option>
          );
        })}
      </select>
      <p>～</p>
      <label htmlFor={`${id}-endTime`} className={styles.visuallyHidden}>
        終了時間
      </label>
      <select
        name={`timeSlot-${slot.id}-endTime`}
        id={`${id}-endTime`}
        value={slot.endTime}
        className={styles.timeSelect}
        onChange={(e) => onChangeTimeSlot({ ...slot, endTime: e.target.value })}
      >
        {selectableTimeRange.map((range) => {
          return (
            <option key={range} value={range}>
              {range}
            </option>
          );
        })}
      </select>
      <button
        type="button"
        onClick={() => onRemoveTimeSlot(slot.id)}
        className={styles.removeButton}
      >
        削除
      </button>
    </div>
  );
}
