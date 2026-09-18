import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import Calculator from "./Calculator";
import { calculate, ApiError } from "../api/calculatorApi";

vi.mock("../api/calculatorApi", () => ({
  calculate: vi.fn(),
  ApiError: class ApiError extends Error {},
}));

const mockedCalculate = vi.mocked(calculate);

describe("Calculator", () => {
  it("renders with an initial display of 0", () => {
    render(<Calculator />);
    expect(screen.getByTestId("display")).toHaveTextContent("0");
  });

  it("performs an addition and displays the result from the API", async () => {
    mockedCalculate.mockResolvedValueOnce(5);
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "2" }));
    await user.click(screen.getByRole("button", { name: "+" }));
    await user.click(screen.getByRole("button", { name: "3" }));
    await user.click(screen.getByRole("button", { name: "=" }));

    expect(mockedCalculate).toHaveBeenCalledWith("add", 2, 3);
    expect(await screen.findByTestId("display")).toHaveTextContent("5");
  });

  it("shows an error message returned by the API instead of a result", async () => {
    mockedCalculate.mockRejectedValueOnce(new ApiError("division by zero"));
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "1" }));
    await user.click(screen.getByRole("button", { name: "÷" }));
    await user.click(screen.getByRole("button", { name: "0" }));
    await user.click(screen.getByRole("button", { name: "=" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("division by zero");
  });

  it("clears the display when AC is pressed", async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "7" }));
    expect(screen.getByTestId("display")).toHaveTextContent("7");

    await user.click(screen.getByRole("button", { name: "AC" }));
    expect(screen.getByTestId("display")).toHaveTextContent("0");
  });
});
