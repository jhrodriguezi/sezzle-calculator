import { afterEach, describe, expect, it, vi } from "vitest";
import { calculate, ApiError } from "./calculatorApi";

function mockFetchResponse(ok: boolean, body: unknown) {
  return vi.fn().mockResolvedValue({
    ok,
    json: () => Promise.resolve(body),
  } as Response);
}

describe("calculate", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("posts the operation payload and resolves with the result", async () => {
    const fetchMock = mockFetchResponse(true, { result: 8 });
    vi.stubGlobal("fetch", fetchMock);

    const result = await calculate("add", 5, 3);

    expect(result).toBe(8);
    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining("/api/calculate"),
      expect.objectContaining({
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ operation: "add", a: 5, b: 3 }),
      }),
    );
  });

  it("throws an ApiError with the server-provided message on failure", async () => {
    vi.stubGlobal("fetch", mockFetchResponse(false, { error: "division by zero" }));

    await expect(calculate("divide", 1, 0)).rejects.toThrow(ApiError);
    await expect(calculate("divide", 1, 0)).rejects.toThrow("division by zero");
  });

  it("falls back to a generic message when the error body is unparsable", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.reject(new Error("bad json")),
    } as unknown as Response);
    vi.stubGlobal("fetch", fetchMock);

    await expect(calculate("add", 1, 1)).rejects.toThrow("Request failed");
  });
});
