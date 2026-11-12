import { request } from "./client";

export const getReviews = (slug, params = {}) => {
  const query = new URLSearchParams(params).toString();
  return request(`/pizzas/${encodeURIComponent(slug)}/reviews${query ? `?${query}` : ""}`);
};
export const createReview = (slug, review) =>
  request(`/pizzas/${encodeURIComponent(slug)}/reviews`, {
    method: "POST",
    body: JSON.stringify(review),
  });
export const getAdminReviews = (status = "pending") =>
  request(`/admin/reviews?status=${encodeURIComponent(status)}`);
export const moderateReview = (id, status) =>
  request(`/admin/reviews/${id}`, { method: "PATCH", body: JSON.stringify({ status }) });