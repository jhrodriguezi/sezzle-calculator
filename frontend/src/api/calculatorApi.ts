import type {
  ApiErrorResponse,
  CalculateRequest,
  CalculateResponse,
  Operation,
} from "../types/calculator";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

export class ApiError extends Error {}

/**
 * The only module allowed to call fetch() against the calculator API.
 */
export async function calculate(operation: Operation, a: number, b: number): Promise<number> {
  const payload: CalculateRequest = { operation, a, b };

  const response = await fetch(`${API_BASE_URL}/api/calculate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    const message = (data as ApiErrorResponse | null)?.error ?? "Request failed";
    throw new ApiError(message);
  }

  return (data as CalculateResponse).result;
}
