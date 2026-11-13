import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { changePassword, currentUser } from "../api/auth";
import Profile from "./Profile";

vi.mock("../api/auth", () => ({
  changePassword: vi.fn(),
  currentUser: vi.fn(),
}));

describe("Profile password change", () => {
  afterEach(() => {
    cleanup();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    currentUser.mockReturnValue({
      display_name: "Test User",
      email: "test@example.com",
      role: "customer",
    });
  });

  it("submits a matching new password and confirms success", async () => {
    changePassword.mockResolvedValue(null);
    render(<Profile />);
    fireEvent.change(screen.getByLabelText("รหัสผ่านปัจจุบัน"), {
      target: { value: "old-password" },
    });
    fireEvent.change(screen.getByLabelText("รหัสผ่านใหม่ (อย่างน้อย 8 ตัวอักษร)"), {
      target: { value: "new-password" },
    });
    fireEvent.change(screen.getByLabelText("ยืนยันรหัสผ่านใหม่"), {
      target: { value: "new-password" },
    });
    fireEvent.click(screen.getByRole("button", { name: "เปลี่ยนรหัสผ่าน" }));

    expect(await screen.findByRole("status")).toHaveTextContent(
      "เปลี่ยนรหัสผ่านเรียบร้อยแล้ว",
    );
    expect(changePassword).toHaveBeenCalledWith("old-password", "new-password");
  });

  it("does not submit when the new passwords do not match", async () => {
    render(<Profile />);
    fireEvent.change(screen.getByLabelText("รหัสผ่านปัจจุบัน"), {
      target: { value: "old-password" },
    });
    fireEvent.change(screen.getByLabelText("รหัสผ่านใหม่ (อย่างน้อย 8 ตัวอักษร)"), {
      target: { value: "new-password" },
    });
    fireEvent.change(screen.getByLabelText("ยืนยันรหัสผ่านใหม่"), {
      target: { value: "different-password" },
    });
    fireEvent.click(screen.getByRole("button", { name: "เปลี่ยนรหัสผ่าน" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "รหัสผ่านใหม่และการยืนยันรหัสผ่านไม่ตรงกัน",
    );
    expect(changePassword).not.toHaveBeenCalled();
  });
});
