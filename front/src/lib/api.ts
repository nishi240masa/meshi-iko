import createClient, { Middleware } from "openapi-fetch";

import { paths } from "../generated/api";
import { tokenStore } from "./tokenStore";

export const api = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
  // credentials: "include", CORSの設定が必要になるので、コメントアウトしておく
});

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    const token = tokenStore.get();

    if (token) {
      request.headers.set("Authorization", `Bearer ${token}`);
    }

    return request;
  },

  async onResponse({ response }) {
    if (response.status === 401) {
      tokenStore.clear();
      // ログイン画面に遷移
      window.location.href = "/login";
    }
  },
};

api.use(authMiddleware);
