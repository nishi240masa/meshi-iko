import { useId } from "react";

import styles from "./TimeSlotCard.module.css";
import { TimeSlotData } from "./timeSlotReducer";

type Props = {
  slot: TimeSlotData;
  onChangeTimeSlot: (slot: TimeSlotData) => void;
  onRemoveTimeSlot: (slotId: string) => void;
  selectableTimes: string[];

  slotNumber: number;
};

export function TimeSlotCard({
  slot,
  onChangeTimeSlot,
  onRemoveTimeSlot,
  slotNumber,
  selectableTimes,
}: Props) {
  const id = useId(); // 複数のTimeSlotコンポーネントがある場合に、input要素のidが重複しないようにするためにuseIdを使用

  return (
    <div className={styles.timeSlotCard}>
      <p>時間帯{slotNumber}</p>
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
        {selectableTimes.map((time) => {
          return (
            <option key={time} value={time}>
              {time}
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
        {selectableTimes.map((time) => {
          return (
            <option key={time} value={time}>
              {time}
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
