import { http, type FetchOptions } from "./ofetch";
import { Result } from "../utils/rust";
import { ApiError } from "../errors";
import { ErrorLevel } from "../constants";

export async function request<Response>(
  path: string,
  options?: FetchOptions<"json">
): Promise<Result<Response, ApiError>> {
  try {
    const data = await http<Response>(path, options);

    return Result.ok(data);
  } catch (error) {
    let message = "未知类型的错误";

    if (error instanceof Error) {
      message = `[API]: 请求失败，响应为： ${error.message}`;
    } else {
      message = String(error);
    }

    const err = new ApiError({ level: ErrorLevel.Error, message });
    return Result.err(err);
  }
}
