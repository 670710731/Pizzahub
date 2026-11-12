import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { register } from "../api/auth";

export default function Register() {
  const [display_name,setName]=useState(""); const [email,setEmail]=useState(""); const [password,setPassword]=useState("");
  const [error,setError]=useState(""); const [busy,setBusy]=useState(false); const navigate=useNavigate();
  async function submit(e){e.preventDefault();setError("");setBusy(true);try{await register({display_name,email,password});navigate("/");window.location.reload();}catch(err){setError(err.message);}finally{setBusy(false);}}
  return <main className="container"><form className="form-card" onSubmit={submit}>
    <h1 className="text-2xl font-black">สมัครสมาชิก</h1><p className="muted">สร้างบัญชีเพื่อเริ่มรีวิวพิซซ่า</p>
    {error && <p className="notice">{error}</p>}
    <label className="field">ชื่อที่แสดง<input required maxLength="100" value={display_name} onChange={e=>setName(e.target.value)}/></label>
    <label className="field">อีเมล<input type="email" required value={email} onChange={e=>setEmail(e.target.value)}/></label>
    <label className="field">รหัสผ่าน (อย่างน้อย 8 ตัว)<input type="password" minLength="8" required value={password} onChange={e=>setPassword(e.target.value)}/></label>
    <button className="btn w-full" disabled={busy}>{busy?"กำลังสมัคร...":"สมัครสมาชิก"}</button>
    <p className="mt-4 text-sm">มีบัญชีแล้ว? <Link className="text-orange-600 font-bold" to="/login">เข้าสู่ระบบ</Link></p>
  </form></main>;
}