export type Operation =
  | "add"
  | "subtract"
  | "multiply"
  | "divide"
  | "exponentiate"
  | "square_root"
  | "percentage";

export interface CalculateRequest {
  operation: Operation;
  a: number;
  b: number;
}

export interface CalculateResponse {
  result: number;
}

export interface ApiErrorResponse {
  error: string;
}
