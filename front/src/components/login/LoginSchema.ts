import z from "zod";

// ログイン時にユーザが入力する値のバリデーションスキーマ
export const loginSchema = z.object({
  name: z.string().trim().min(1, "Nameは必須です").max(32, "Nameは最大32文字です"),
});
