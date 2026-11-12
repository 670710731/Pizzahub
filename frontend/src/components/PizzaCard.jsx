import { Link } from "react-router-dom";

const fallback = "https://images.unsplash.com/photo-1513104890138-7c749659a591?w=800";
export default function PizzaCard({ pizza }) {
  const image = pizza.images?.[0]?.image_url || fallback;
  return (
    <article className="card">
      <img className="pizza-image" src={image} alt={pizza.name} />
      <div className="card-body">
        <div className="flex justify-between items-center gap-2">
          <h3 className="font-bold text-lg">{pizza.name}</h3>
          <span className="rating">★ {Number(pizza.rating_average || 0).toFixed(1)}</span>
        </div>
        <p className="muted text-sm my-2">{pizza.description || "ยังไม่มีคำอธิบาย"}</p>
        <p className="text-sm muted mb-3">{pizza.review_count || 0} รีวิว</p>
        <Link className="btn" to={`/pizzas/${pizza.slug}`}>ดูรายละเอียด</Link>
      </div>
    </article>
  );
}