import { createFileRoute } from "@tanstack/react-router";

import { AuthForm } from "../components/login/AuthForm";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  return (
    <div>
      <h1>Login</h1>
      <AuthForm />
    </div>
  );
}
