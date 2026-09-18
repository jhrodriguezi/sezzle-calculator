import { useState } from "react";
import { calculate, ApiError } from "../api/calculatorApi";
import type { Operation } from "../types/calculator";
import "./Calculator.css";

const OPERATOR_SYMBOLS: Record<Operation, string> = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
  exponentiate: "^",
  percentage: "%",
  square_root: "√",
};

const MAX_DIGITS = 15;

export default function Calculator() {
  const [display, setDisplay] = useState("0");
  const [previousValue, setPreviousValue] = useState<number | null>(null);
  const [pendingOperation, setPendingOperation] = useState<Operation | null>(null);
  const [overwrite, setOverwrite] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [expression, setExpression] = useState("");

  function inputDigit(digit: string) {
    setError(null);
    if (overwrite) {
      setDisplay(digit === "." ? "0." : digit);
      setOverwrite(false);
      return;
    }
    if (display.length >= MAX_DIGITS) return;
    if (digit === "." && display.includes(".")) return;
    setDisplay(display === "0" && digit !== "." ? digit : display + digit);
  }

  function clear() {
    setDisplay("0");
    setPreviousValue(null);
    setPendingOperation(null);
    setOverwrite(true);
    setError(null);
    setExpression("");
  }

  function backspace() {
    if (error) {
      setError(null);
      setDisplay("0");
      setOverwrite(true);
      return;
    }
    if (overwrite) return;
    const next = display.slice(0, -1);
    setDisplay(next === "" || next === "-" ? "0" : next);
  }

  function toggleSign() {
    setError(null);
    if (display === "0") return;
    setDisplay(display.startsWith("-") ? display.slice(1) : `-${display}`);
  }

  async function runCalculation(a: number, b: number, operation: Operation) {
    setIsLoading(true);
    setError(null);
    try {
      const result = await calculate(operation, a, b);
      return result;
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Unable to reach the calculator API.");
      return null;
    } finally {
      setIsLoading(false);
    }
  }

  async function chooseOperation(operation: Operation) {
    const current = Number(display);

    if (pendingOperation && !overwrite && previousValue !== null) {
      const result = await runCalculation(previousValue, current, pendingOperation);
      if (result === null) {
        resetAfterError();
        return;
      }
      setPreviousValue(result);
      setDisplay(String(result));
      setExpression(`${result} ${OPERATOR_SYMBOLS[operation]}`);
    } else {
      setPreviousValue(current);
      setExpression(`${current} ${OPERATOR_SYMBOLS[operation]}`);
    }

    setPendingOperation(operation);
    setOverwrite(true);
  }

  async function equals() {
    if (pendingOperation === null || previousValue === null || overwrite) return;

    const current = Number(display);
    const result = await runCalculation(previousValue, current, pendingOperation);
    if (result === null) {
      resetAfterError();
      return;
    }

    setExpression(`${previousValue} ${OPERATOR_SYMBOLS[pendingOperation]} ${current} =`);
    setDisplay(String(result));
    setPreviousValue(null);
    setPendingOperation(null);
    setOverwrite(true);
  }

  /** Square root is unary: it applies immediately to the current value, with no pending operation. */
  async function applySquareRoot() {
    const current = Number(display);
    const result = await runCalculation(current, 0, "square_root");
    if (result === null) {
      resetAfterError();
      return;
    }

    setExpression(`√(${current}) =`);
    setDisplay(String(result));
    setPreviousValue(null);
    setPendingOperation(null);
    setOverwrite(true);
  }

  function resetAfterError() {
    setPreviousValue(null);
    setPendingOperation(null);
    setExpression("");
    setOverwrite(true);
  }

  return (
    <div className="calculator">
      <div className="screen">
        <div className="expression" aria-live="polite">
          {error ? "\u00A0" : expression || "\u00A0"}
        </div>
        <div
          className={`display${error ? " display-error" : ""}`}
          role={error ? "alert" : undefined}
          aria-live="polite"
          data-testid="display"
        >
          {isLoading ? "…" : error ?? display}
        </div>
      </div>

      <div className="advanced-row">
        <button type="button" className="key advanced" onClick={applySquareRoot}>
          √x
        </button>
        <button type="button" className="key advanced" onClick={() => chooseOperation("exponentiate")}>
          x<sup>y</sup>
        </button>
        <button type="button" className="key advanced" onClick={() => chooseOperation("percentage")}>
          %
        </button>
      </div>

      <div className="keypad">
        <button type="button" className="key function" onClick={clear}>
          AC
        </button>
        <button type="button" className="key function" onClick={toggleSign}>
          +/−
        </button>
        <button type="button" className="key function" onClick={backspace} aria-label="Backspace">
          ⌫
        </button>
        <button type="button" className="key operator" onClick={() => chooseOperation("divide")}>
          ÷
        </button>

        <button type="button" className="key digit" onClick={() => inputDigit("7")}>
          7
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit("8")}>
          8
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit("9")}>
          9
        </button>
        <button type="button" className="key operator" onClick={() => chooseOperation("multiply")}>
          ×
        </button>

        <button type="button" className="key digit" onClick={() => inputDigit("4")}>
          4
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit("5")}>
          5
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit("6")}>
          6
        </button>
        <button type="button" className="key operator" onClick={() => chooseOperation("subtract")}>
          −
        </button>

        <button type="button" className="key digit" onClick={() => inputDigit("1")}>
          1
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit("2")}>
          2
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit("3")}>
          3
        </button>
        <button type="button" className="key operator" onClick={() => chooseOperation("add")}>
          +
        </button>

        <button type="button" className="key digit zero" onClick={() => inputDigit("0")}>
          0
        </button>
        <button type="button" className="key digit" onClick={() => inputDigit(".")}>
          .
        </button>
        <button type="button" className="key equals" onClick={equals}>
          =
        </button>
      </div>
    </div>
  );
}
