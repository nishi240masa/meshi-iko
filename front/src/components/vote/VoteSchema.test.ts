import { describe, expect, test } from "vitest";

import { voteFormSchema } from "./VoteSchema";

describe("voteFormSchema", () => {
  describe("正常系", () => {
    test("「いく」で有効なスロットが1つある場合、成功する", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [{ id: "1", startTime: "16:00", endTime: "17:00" }],
      });
      expect(result.success).toBe(true);
    });

    test("「いく」で時間が重複しない複数スロットがある場合、成功する", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [
          { id: "1", startTime: "16:00", endTime: "17:00" },
          { id: "2", startTime: "18:00", endTime: "19:00" },
        ],
      });
      expect(result.success).toBe(true);
    });

    test("終了と開始が接しているスロット（例: 17:00終了と17:00開始）は重複とみなさず成功する", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [
          { id: "1", startTime: "16:00", endTime: "17:00" },
          { id: "2", startTime: "17:00", endTime: "18:00" },
        ],
      });
      expect(result.success).toBe(true);
    });

    test("「いかない」でスロットが空の場合、成功する", () => {
      const result = voteFormSchema.safeParse({
        status: "unavailable",
        timeSlots: [],
      });
      expect(result.success).toBe(true);
    });
  });

  describe("異常系", () => {
    test("status が未選択の場合、エラーになる", () => {
      const result = voteFormSchema.safeParse({
        status: null,
        timeSlots: [],
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe("いくいかないの選択は必須です");
      }
    });

    test("「いく」なのにスロットが空の場合、エラーになる", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [],
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe("いくの場合、一つ以上のスロットが必要です");
      }
    });

    test("開始時刻と終了時刻が同じ場合、エラーになる", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [{ id: "1", startTime: "16:00", endTime: "16:00" }],
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe(
          "終了時刻は開始時間より後の時刻を選択してください",
        );
      }
    });

    test("開始時刻が終了時刻より遅い場合、エラーになる", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [{ id: "1", startTime: "17:00", endTime: "16:00" }],
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe(
          "終了時刻は開始時間より後の時刻を選択してください",
        );
      }
    });

    test("スロット同士の時間帯が重複している場合、エラーになる", () => {
      const result = voteFormSchema.safeParse({
        status: "available",
        timeSlots: [
          { id: "1", startTime: "16:00", endTime: "17:30" },
          { id: "2", startTime: "17:00", endTime: "18:00" },
        ],
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe(
          "時間帯が重複しています。別の時間帯を選択してください。",
        );
      }
    });
  });
});
