// バリデーションルール
// 1. 開始時間 < 終了時間であること
// 2. 開始時間と終了時間は選択可能な時間帯の範囲内であること
// 3. 開始時間と終了時間の組み合わせが重複しないこと
// 4. 開始時間と終了時間は交差しないこと

import z from "zod";

// 各スロットの単体のバリデーションスキーマ
export const timeSlotSchema = z
  .object({
    id: z.string(),
    startTime: z.string().min(1, "開始時間は必須です"),
    endTime: z.string().min(1, "終了時間は必須です"),
  })
  .refine((slot) => slot.startTime < slot.endTime, {
    message: "終了時刻は開始時間より後の時刻を選択してください",
    path: ["endTime"],
  });

// 投票フォーム全体のバリデーションスキーマ
export const voteFormSchema = z
  .object({
    status: z.enum(["available", "unavailable"], "いくいかないの選択は必須です"),
    timeSlots: z.array(timeSlotSchema),
  })
  .refine(
    (data) => {
      // いくの場合、一つ以上のスロットが必要
      if (data.status === "available" && data.timeSlots.length === 0) {
        return false;
      }
      return true;
    },
    {
      message: "いくの場合、一つ以上のスロットが必要です",
      path: ["timeSlots"],
    },
  )
  .refine(
    (data) => {
      // スロット同士の重複は許可しない

      // いかないの場合はスロットの重複チェックをスキップ
      if (data.status === "unavailable") {
        return true;
      }
      for (let i = 0; i < data.timeSlots.length; i++) {
        for (let j = i + 1; j < data.timeSlots.length; j++) {
          const a = data.timeSlots[i];
          const b = data.timeSlots[j];
          // 区間が重なっているか判定
          if (a.startTime < b.endTime && b.startTime < a.endTime) {
            return false;
          }
        }
      }
      return true;
    },
    {
      message: "時間帯が重複しています。別の時間帯を選択してください。",
      path: ["timeSlots"],
    },
  );

export type VoteFormData = z.infer<typeof voteFormSchema>;
