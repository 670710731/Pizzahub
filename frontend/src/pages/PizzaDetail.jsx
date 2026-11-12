import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getPizza } from "../api/pizzas";
import { getReviews } from "../api/reviews";
import { currentUser } from "../api/auth";
import ReviewForm from "../components/ReviewForm";

export default function PizzaDetail() {
  const { slug } = useParams();
  const [pizza, setPizza] = useState(null);
  const [reviews, setReviews] = useState([]);
  const [error, setError] = useState("");
  const user = currentUser();
  const loadReviews = useCallback(() => getReviews(slug).then(d => setReviews(d.items || [])).catch(e => setError(e.message)), [slug]);
  useEffect(() => {
    getPizza(slug).then(setPizza).catch(e => setError(e.message));
    loadReviews();
  }, [slug, loadReviews]);
  if (error && !pizza) return <main className="container"><p className="notice">{error}</p><Link to="/" className="btn">กลับหน้าแรก</Link></main>;
  if (!pizza) return <main className="container py-10">กำลังโหลด...</main>;
  const image = pizza.images?.[0]?.image_url || "https://images.unsplash.com/photo-1513104890138-7c749659a591?w=1000";
  return <main className="container py-8">
    <Link to="/" className="muted">← กลับหน้าแรก</Link>
    <div className="grid mt-5" style={{gridTemplateColumns:"repeat(auto-fit,minmax(280px,1fr))"}}>
      <img className="pizza-image rounded-2xl" style={{height:340}} src={image} alt={pizza.name}/>
      <section>
        <span className="pill">{pizza.categories?.map(c=>c.name).join(" / ") || "Pizza"}</span>
        <h1 className="text-3xl font-black my-3">{pizza.name}</h1>
        <p className="rating text-xl">★ {Number(pizza.rating_average || 0).toFixed(1)} <span className="muted text-sm">({pizza.review_count || 0} รีวิว)</span></p>
        <p className="my-4">{pizza.description}</p>
        <h3 className="font-bold">ขนาดและราคา</h3>
        {pizza.variants?.length ? pizza.variants.map(v=><p key={v.id} className="my-2">{v.size_label} — {Number(v.price).toLocaleString("th-TH")} {v.currency || "THB"}</p>) : <p className="muted">ยังไม่มีข้อมูลราคา</p>}
      </section>
    </div>
    <h2 className="section-title">รีวิวจากสมาชิก</h2>
    {reviews.map(r=><article className="card card-body mb-3" key={r.id}>
      <strong>{r.title || "รีวิวจากสมาชิก"}</strong><p className="muted text-sm">โดย {r.display_name}</p><p className="rating">{"★".repeat(r.rating || 0)}</p><p>{r.body}</p>
    </article>)}
    {!reviews.length && <p className="muted mb-4">ยังไม่มีรีวิว เป็นคนแรกที่รีวิวได้เลย</p>}
    {user ? <ReviewForm slug={slug} onCreated={loadReviews}/> : <p className="notice">กรุณา <Link className="font-bold" to="/login">เข้าสู่ระบบ</Link> เพื่อเขียนรีวิว</p>}
  </main>;
}