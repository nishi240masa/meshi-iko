import { loginSchema } from "./LoginSchema";

export function LoginForm() {
  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
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
    fetch(`${import.meta.env.VITE_API_LOCAL_BASE_URL}/users/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name: result.data.name }),
    }).then((response) => {
      if (response.ok) {
        console.log("Login successful");
      } else {
        console.error("Login failed");
      }
    });
  };
  return (
    <form onSubmit={handleSubmit}>
      <label htmlFor="name">name</label>
      <input type="text" name="name" id="name" />
      <button type="submit">Login</button>
    </form>
  );
}
