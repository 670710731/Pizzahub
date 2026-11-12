import { request } from "./client";

export const getPizzas = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return request(`/pizzas${query ? `?${query}` : ""}`);
};
export const getPizza = (slug) => request(`/pizzas/${encodeURIComponent(slug)}`);
export const getCategories = () => request("/categories");