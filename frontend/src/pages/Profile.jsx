import { useState } from "react";
import { changePassword, currentUser } from "../api/auth";

export default function Profile() {
  const user = currentUser();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [busy, setBusy] = useState(false);

  if (!user) return <main className="container py-10">กรุณาเข้าสู่ระบบก่อน</main>;

  async function submit(e) {
    e.preventDefault();
    setError("");
    setSuccess(false);
    if (newPassword !== confirmPassword) {
      setError("รหัสผ่านใหม่และการยืนยันรหัสผ่านไม่ตรงกัน");
      return;
    }
    setBusy(true);
    try {
      await changePassword(currentPassword, newPassword);
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
      setSuccess(true);
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  }

  return <main className="container py-10">
    <section className="card card-body">
      <h1 className="text-2xl font-black">โปรไฟล์สมาชิก</h1>
      <p className="mt-4"><b>ชื่อ:</b> {user.display_name}</p>
      <p><b>อีเมล:</b> {user.email}</p>
      <p><b>สถานะ:</b> <span className="pill">{user.role}</span></p>
    </section>
    <form className="card card-body mt-5" onSubmit={submit}>
      <h2 className="text-xl font-bold">เปลี่ยนรหัสผ่าน</h2>
      {error && <p className="notice" role="alert">{error}</p>}
      {success && <p className="notice" role="status">เปลี่ยนรหัสผ่านเรียบร้อยแล้ว</p>}
      <label className="field">รหัสผ่านปัจจุบัน
        <input type="password" autoComplete="current-password" required value={currentPassword} onChange={e => setCurrentPassword(e.target.value)} />
      </label>
      <label className="field">รหัสผ่านใหม่ (อย่างน้อย 8 ตัวอักษร)
        <input type="password" autoComplete="new-password" minLength="8" required value={newPassword} onChange={e => setNewPassword(e.target.value)} />
      </label>
      <label className="field">ยืนยันรหัสผ่านใหม่
        <input type="password" autoComplete="new-password" minLength="8" required value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} />
      </label>
      <button className="btn" disabled={busy}>{busy ? "กำลังเปลี่ยน..." : "เปลี่ยนรหัสผ่าน"}</button>
    </form>
  </main>;
}