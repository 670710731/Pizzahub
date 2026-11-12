import { useEffect, useState } from "react";
import { getPizzas } from "../api/pizzas";
import PizzaCard from "../components/PizzaCard";

export default function Home() {
  const [pizzas, setPizzas] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    getPizzas({ page: 1, page_size: 30, sort: "rating" })
      .then(data => setPizzas(data.items || []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, []);
  return <main className="container">
    <section className="hero">
      <p className="pill">ค้นพบพิซซ่าจานโปรดของคุณ</p>
      <h1>รีวิวพิซซ่า<br/>อร่อยจริง ต้องบอกต่อ 🍕</h1>
      <p className="text-lg">ค้นหา เปรียบเทียบ และแบ่งปันประสบการณ์พิซซ่าที่คุณชื่นชอบ</p>
    </section>
    <h2 className="section-title">พิซซ่ายอดนิยม</h2>
    {error && <p className="notice">เชื่อมต่อ Backend ไม่สำเร็จ: {error}</p>}
    {loading ? <p>กำลังโหลดข้อมูล...</p> : pizzas.length ?
      <div className="grid">{pizzas.map(p => <PizzaCard key={p.id} pizza={p}/>)}</div> :
      !error && <p className="notice">ยังไม่มีข้อมูลพิซซ่าในระบบ กรุณาเพิ่มข้อมูลจาก Backend ก่อน</p>}
  </main>;
}