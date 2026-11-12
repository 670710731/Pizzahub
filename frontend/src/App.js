import { Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import PizzaDetail from "./pages/PizzaDetail";
import Profile from "./pages/Profile";
import Admin from "./pages/Admin";

export default function App() {
  return <>
    <Navbar />
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/pizzas/:slug" element={<PizzaDetail />} />
      <Route path="/profile" element={<Profile />} />
      <Route path="/admin" element={<Admin />} />
      <Route path="*" element={<main className="container py-12"><h1 className="text-2xl font-bold">ไม่พบหน้าที่ต้องการ</h1></main>} />
    </Routes>
    <footer className="footer">© {new Date().getFullYear()} PizzaHub · รีวิวพิซซ่าที่คุณชื่นชอบ</footer>
  </>;
}