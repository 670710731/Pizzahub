import { fireEvent, render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { createReview } from "../api/reviews";
import ReviewForm from "./ReviewForm";

vi.mock("../api/reviews", () => ({
  createReview: vi.fn(),
}));

describe("ReviewForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("confirms successful submissions are waiting for moderation", async () => {
    createReview.mockResolvedValue({ status: "pending" });
    const onCreated = vi.fn();
    render(<ReviewForm slug="margherita" onCreated={onCreated} />);

    fireEvent.change(screen.getByLabelText("หัวข้อรีวิว"), {
      target: { value: "อร่อยมาก" },
    });
    fireEvent.change(screen.getByLabelText("รายละเอียด"), {
      target: { value: "แป้งนุ่มและชีสเข้มข้น" },
    });
    fireEvent.click(screen.getByRole("button", { name: "ส่งรีวิว" }));

    expect(await screen.findByRole("status")).toHaveTextContent(
      "ส่งรีวิวแล้ว รีวิวจะปรากฏหลังผู้ดูแลอนุมัติ",
    );
    expect(createReview).toHaveBeenCalledWith("margherita", {
      rating: 5,
      title: "อร่อยมาก",
      body: "แป้งนุ่มและชีสเข้มข้น",
    });
    expect(onCreated).toHaveBeenCalledOnce();
  });
});
