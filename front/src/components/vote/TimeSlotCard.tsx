import { useId } from "react";

import { createTimeScale } from "./createTimeScale";
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
    <div>
      <div>
        <p>区間{number}</p>
        <div>
          <label htmlFor={`${id}-startTime`}>開始時間</label>
          <select
            name={`timeSlot-${slot.id}-startTime`}
            id={`${id}-startTime`}
            value={slot.startTime}
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
          <label htmlFor={`${id}-endTime`}>終了時間</label>
          <select
            name={`timeSlot-${slot.id}-endTime`}
            id={`${id}-endTime`}
            value={slot.endTime}
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
        </div>
      </div>
      <button type="button" onClick={() => onRemoveTimeSlot(slot.id)}>
        削除
      </button>
    </div>
  );
}
