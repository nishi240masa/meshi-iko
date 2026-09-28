export type TimeSlotData = {
  id: string;
  startTime: string;
  endTime: string;
};

export type TimeSlotState = {
  slots: TimeSlotData[];
};

export type TimeSlotAction =
  | { type: "ADD"; payload: { startTime: string; endTime: string } }
  | { type: "REMOVE"; payload: { id: string } }
  | { type: "UPDATE"; payload: TimeSlotData };

export const initialTimeSlotState = (): TimeSlotState => ({
  slots: [],
});

// TimeSlotの開始時間と終了時間を管理するためのreducer関数
export function timeSlotReducer(state: TimeSlotState, action: TimeSlotAction): TimeSlotState {
  switch (action.type) {
    // 追加
    case "ADD":
      return {
        ...state,
        slots: [
          ...state.slots,
          {
            id: crypto.randomUUID(),
            startTime: action.payload.startTime, //初期値は空文字
            endTime: action.payload.endTime,
          },
        ],
      };
    // 削除
    case "REMOVE":
      return {
        ...state,
        slots: state.slots.filter((slot) => slot.id !== action.payload.id),
      };
    // 更新
    case "UPDATE":
      return {
        ...state,
        // 更新対象のスロットを見つけて更新
        slots: state.slots.map((slot) =>
          slot.id === action.payload.id ? { ...slot, ...action.payload } : slot,
        ),
      };
    default:
      throw Error("Unknown action: " + action);
  }
}
