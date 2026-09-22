import * as React from "react";
import ReactDOM from "react-dom/client";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { routeTree } from './routeTree.gen'

const router = createRouter({
  routeTree,
})

// Tanstack Routerの型定義を拡張して型補完できるようにする
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

ReactDOM.createRoot(
  document.getElementById('root')!,
).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)