import { Link, NavLink, useNavigate } from "react-router-dom";
import { currentUser, logout } from "../api/auth";

export default function Navbar() {
  const navigate = useNavigate();
  const user = currentUser();
  const signOut = () => { logout(); navigate("/"); window.location.reload(); };
  return (
    <header className="nav">
      <div className="container nav-inner">
        <Link className="brand" to="/">🍕 PizzaHub</Link>
        <nav className="nav-links">
          <NavLink to="/">หน้าแรก</NavLink>
          {user ? <>
            <NavLink to="/profile">โปรไฟล์</NavLink>
            {user.role === "admin" && <NavLink to="/admin">จัดการระบบ</NavLink>}
            <button className="btn secondary" onClick={signOut}>ออกจากระบบ</button>
          </> : <>
            <NavLink to="/login">เข้าสู่ระบบ</NavLink>
            <Link className="btn" to="/register">สมัครสมาชิก</Link>
          </>}
        </nav>
      </div>
    </header>
  );
}