import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, expect, it, vi } from "vitest";
import App from "./App";

beforeEach(() => {
  localStorage.clear();
  vi.stubGlobal("fetch", vi.fn().mockResolvedValue({
    ok: true,
    headers: { get: () => "application/json" },
    json: async () => ({ items: [] }),
  }));
});

it("renders the homepage and empty-state message", async () => {
  render(
    <MemoryRouter initialEntries={["/"]}>
      <App />
    </MemoryRouter>
  );

  expect(screen.getByRole("link", { name: /PizzaHub/i })).toBeInTheDocument();
  expect(await screen.findByText("ยังไม่มีข้อมูลพิซซ่าในระบบ กรุณาเพิ่มข้อมูลจาก Backend ก่อน")).toBeInTheDocument();
});
