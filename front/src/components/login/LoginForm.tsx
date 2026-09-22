import { useNavigate } from "@tanstack/react-router";
import { useState } from "react";

import { api } from "../../lib/api";
import styles from "./LoginForm.module.css";
import { loginSchema } from "./LoginSchema";

export function LoginForm() {
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState<string | undefined>(undefined);
  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const userName = formData.get("name");
    const result = loginSchema.safeParse({ name: userName });
    if (result.success) {
      console.log("Name:", result.data.name);
    } else {
      setErrorMessage(result.error.issues[0].message ?? "不正な入力です");
      console.error("Validation error:", result.error);
      return;
    }

    // ログインリクエスト
    try {
      const { data, error } = await api.POST("/users/login", {
        body: { name: result.data.name },
      });
      if (error) {
        setErrorMessage(error.message ?? "ログインに失敗しました");
        console.error("Login error:", error);
        return;
      }

      console.log("Login response:", data);
      //ログイン成功時は投票ページに遷移
      navigate({ to: "/vote" });
    } catch (error) {
      setErrorMessage("ログインに失敗しました");
      console.error("Login request error:", error);
    }
  };
  return (
    <form onSubmit={handleSubmit} className={styles.loginForm}>
      <div className={styles.formContainer}>
        <div className={styles.formGroup}>
          <label htmlFor="name">name</label>
          <input type="text" name="name" id="name" />
        </div>
        {errorMessage && <p className={styles.errorMessage}>{errorMessage}</p>}
      </div>
      <button type="submit" className={styles.submitButton}>
        Login
      </button>
    </form>
  );
}
