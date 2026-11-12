import { currentUser } from "../api/auth";
export default function Profile() {
  const user = currentUser();
  if (!user) return <main className="container py-10">กรุณาเข้าสู่ระบบก่อน</main>;
  return <main className="container py-10"><section className="card card-body">
    <h1 className="text-2xl font-black">โปรไฟล์สมาชิก</h1>
    <p className="mt-4"><b>ชื่อ:</b> {user.display_name}</p>
    <p><b>อีเมล:</b> {user.email}</p>
    <p><b>สถานะ:</b> <span className="pill">{user.role}</span></p>
  </section></main>;
}