import { ofetch, type FetchOptions } from "ofetch";
import { HTTP_DEFAULT_CONFIG } from "../constants/http";
import { optionizeDeep } from "../utils/useful";

const http = ofetch.create({
  baseURL: HTTP_DEFAULT_CONFIG.BASE_URL,
  headers: { "Content-Type": HTTP_DEFAULT_CONFIG.CONTENT_TYPE },
  timeout: HTTP_DEFAULT_CONFIG.TIMEOUT,
  onResponse({ response }) {
    response._data = optionizeDeep(response._data);
  }
});

export { http };
export type { FetchOptions };
