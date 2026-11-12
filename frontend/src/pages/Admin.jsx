import { useEffect, useState } from "react";
import { getAdminReviews, moderateReview } from "../api/reviews";
import { currentUser } from "../api/auth";

export default function Admin() {
  const user = currentUser();
  const [reviews,setReviews]=useState([]); const [error,setError]=useState("");
  const load=()=>getAdminReviews("pending").then(d=>setReviews(d.items||[])).catch(e=>setError(e.message));
  useEffect(()=>{if(user?.role==="admin") load();},[user?.role]);
  if(user?.role!=="admin") return <main className="container py-10"><p className="notice">หน้านี้สำหรับผู้ดูแลระบบเท่านั้น</p></main>;
  async function update(id,status){try{await moderateReview(id,status);load();}catch(e){setError(e.message);}}
  return <main className="container py-8">
    <h1 className="text-3xl font-black">Admin Dashboard</h1><p className="muted">จัดการรีวิวที่รอตรวจสอบ</p>
    {error&&<p className="notice">{error}</p>}
    <h2 className="section-title">รีวิวรออนุมัติ ({reviews.length})</h2>
    {reviews.map(r=><article className="card card-body mb-3" key={r.id}>
      <div className="flex justify-between gap-3 flex-wrap"><b>{r.title||"รีวิว"}</b><span className="rating">{"★".repeat(r.rating||0)}</span></div>
      <p className="muted text-sm">พิซซ่า: {r.pizza_slug} · โดย {r.display_name}</p><p className="my-2">{r.body}</p>
      <div className="flex gap-2"><button className="btn" onClick={()=>update(r.id,"published")}>อนุมัติ</button><button className="btn danger" onClick={()=>update(r.id,"rejected")}>ปฏิเสธ</button></div>
    </article>)}
    {!reviews.length&&!error&&<p className="notice">ไม่มีรีวิวรอตรวจสอบ</p>}
  </main>;
}