import { useState } from "react";
import { createReview } from "../api/reviews";

export default function ReviewForm({ slug, onCreated }) {
  const [rating, setRating] = useState(5);
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [busy, setBusy] = useState(false);

  async function submit(e) {
    e.preventDefault(); setError(""); setSuccess(false); setBusy(true);
    try {
      await createReview(slug, { rating: Number(rating), title, body });
      setTitle(""); setBody(""); setSuccess(true); onCreated?.();
    } catch (err) { setError(err.message); }
    finally { setBusy(false); }
  }
  return <form onSubmit={submit} className="card card-body">
    <h3 className="font-bold text-lg">เขียนรีวิวของคุณ</h3>
    {error && <p className="notice">{error}</p>}
    {success && <p className="notice" role="status">ส่งรีวิวแล้ว รีวิวจะปรากฏหลังผู้ดูแลอนุมัติ</p>}
    <label className="field">คะแนน
      <select value={rating} onChange={e => setRating(e.target.value)}>
        {[5,4,3,2,1].map(n => <option key={n} value={n}>{n} ดาว</option>)}
      </select>
    </label>
    <label className="field">หัวข้อรีวิว<input value={title} onChange={e => setTitle(e.target.value)} maxLength="150" /></label>
    <label className="field">รายละเอียด<textarea rows="4" value={body} onChange={e => setBody(e.target.value)} maxLength="5000" /></label>
    <button className="btn" disabled={busy}>{busy ? "กำลังส่ง..." : "ส่งรีวิว"}</button>
  </form>;
}