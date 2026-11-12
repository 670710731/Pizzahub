import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../api/auth";

export default function Login() {
  const [email,setEmail]=useState(""); const [password,setPassword]=useState("");
  const [error,setError]=useState(""); const [busy,setBusy]=useState(false); const navigate=useNavigate();
  async function submit(e){e.preventDefault();setError("");setBusy(true);try{await login({email,password});navigate("/");window.location.reload();}catch(err){setError(err.message);}finally{setBusy(false);}}
  return <main className="container"><form className="form-card" onSubmit={submit}>
    <h1 className="text-2xl font-black">เข้าสู่ระบบ</h1><p className="muted">ยินดีต้อนรับกลับสู่ PizzaHub</p>
    {error && <p className="notice">{error}</p>}
    <label className="field">อีเมล<input type="email" required value={email} onChange={e=>setEmail(e.target.value)}/></label>
    <label className="field">รหัสผ่าน<input type="password" required value={password} onChange={e=>setPassword(e.target.value)}/></label>
    <button className="btn w-full" disabled={busy}>{busy?"กำลังเข้าสู่ระบบ...":"เข้าสู่ระบบ"}</button>
    <p className="mt-4 text-sm">ยังไม่มีบัญชี? <Link className="text-orange-600 font-bold" to="/register">สมัครสมาชิก</Link></p>
  </form></main>;
}