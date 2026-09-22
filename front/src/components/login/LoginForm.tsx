import { api } from "../../lib/api";
import { loginSchema } from "./LoginSchema";

export function LoginForm() {
  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const userName = formData.get("name");
    const result = loginSchema.safeParse({ name: userName });
    if (result.success) {
      console.log("Name:", result.data.name);
    } else {
      console.error("Validation error:", result.error);
    }

    if (!result.data?.name) {
      console.log("Name is required");
      return;
    }

    // ログインリクエスト
    const { data, error } = await api.POST("/users/login", {
      body: { name: result.data.name },
    });
    if (error) {
      console.error("Login error:", error);
    }

    console.log("Login response:", data);
  };
  return (
    <form onSubmit={handleSubmit}>
      <label htmlFor="name">name</label>
      <input type="text" name="name" id="name" />
      <button type="submit">Login</button>
    </form>
  );
}
