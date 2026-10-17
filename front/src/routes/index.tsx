import "./__root.css";

import { createFileRoute, Link } from "@tanstack/react-router";

import { Greet } from "../components/Greet";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  return (
    <main className="container">
      <h1>Welcome to Tauri + React</h1>
      <Greet />
      <Link to="/login">ログイン</Link>
    </main>
  );
}
