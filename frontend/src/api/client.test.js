import { beforeEach, describe, expect, it, vi } from "vitest";
import { getToken, request } from "./client";

describe("API client", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.unstubAllEnvs();
    vi.restoreAllMocks();
  });

  it("uses the local API base path and returns JSON", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      headers: { get: () => "application/json; charset=utf-8" },
      json: async () => ({ status: "ok" }),
    });
    vi.stubGlobal("fetch", fetchMock);

    await expect(request("/healthz")).resolves.toEqual({ status: "ok" });
    expect(fetchMock).toHaveBeenCalledWith("/api/v1/healthz", {
      headers: { "Content-Type": "application/json" },
    });
  });

  it("adds the bearer token and reports API errors", async () => {
    localStorage.setItem("pizzahub_token", "test-token");
    const fetchMock = vi.fn().mockResolvedValue({
      ok: false,
      status: 401,
      headers: { get: () => "application/json" },
      json: async () => ({ error: "authentication required" }),
    });
    vi.stubGlobal("fetch", fetchMock);

    await expect(request("/admin/reviews")).rejects.toThrow("authentication required");
    expect(fetchMock.mock.calls[0][1].headers.Authorization).toBe("Bearer test-token");
    expect(getToken()).toBe("test-token");
  });

  it("does not try to parse a non-JSON success body", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({
      ok: true,
      headers: { get: () => "text/plain" },
      text: async () => "ok",
    }));

    await expect(request("/healthz")).resolves.toBeNull();
  });

  it("removes trailing slashes from the configured API URL", async () => {
    vi.stubEnv("VITE_API_URL", "https://api.example.com/api/v1/");
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      headers: { get: () => "application/json" },
      json: async () => ({ items: [] }),
    });
    vi.stubGlobal("fetch", fetchMock);

    await request("/pizzas");

    expect(fetchMock.mock.calls[0][0]).toBe("https://api.example.com/api/v1/pizzas");
  });
});
