import { createFileRoute } from "@tanstack/react-router";

import { AuthForm } from "../components/login/AuthForm";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  return (
    <div className="container">
      <h1>ログイン</h1>
      <AuthForm />
    </div>
  );
}
