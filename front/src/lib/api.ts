import createClient from "openapi-fetch";

import { paths } from "../generated/api";

export const api = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
  // credentials: "include", CORSの設定が必要になるので、コメントアウトしておく
});
