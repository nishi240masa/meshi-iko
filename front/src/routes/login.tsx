import { createFileRoute } from "@tanstack/react-router";

import { LoginForm } from "../components/login/LoginForm";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  return (
    <div>
      <h1>Login</h1>
      <LoginForm />
    </div>
  );
}
