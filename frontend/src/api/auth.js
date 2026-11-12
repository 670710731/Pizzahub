import { request } from "./client";

export async function login(credentials) {
  const result = await request("/auth/login", { method: "POST", body: JSON.stringify(credentials) });
  localStorage.setItem("pizzahub_token", result.access_token);
  localStorage.setItem("pizzahub_user", JSON.stringify(result.user));
  return result;
}
export async function register(details) {
  const result = await request("/auth/register", { method: "POST", body: JSON.stringify(details) });
  localStorage.setItem("pizzahub_token", result.access_token);
  localStorage.setItem("pizzahub_user", JSON.stringify(result.user));
  return result;
}
export function logout() {
  localStorage.removeItem("pizzahub_token");
  localStorage.removeItem("pizzahub_user");
}
export function currentUser() {
  try { return JSON.parse(localStorage.getItem("pizzahub_user") || "null"); }
  catch { return null; }
}