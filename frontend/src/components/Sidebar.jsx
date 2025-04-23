import React from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { logoutUser } from "../utils/auth";

const Sidebar = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    logoutUser();
    navigate("/login");
  };

  const linkClass =
    "block px-4 py-2 rounded-md text-white hover:bg-blue-600 transition";

  return (
    <div className="h-screen w-60 bg-blue-800 p-4 fixed top-0 left-0 shadow-lg">
      <h2 className="text-white text-xl font-bold mb-8">IIITN Visitors</h2>
      <nav className="flex flex-col space-y-3 text-xl gap-8">
        <NavLink to="/dashboard" className={linkClass}>Dashboard</NavLink>
        <NavLink to="/book" className={linkClass}>Book Appointment</NavLink>
        <NavLink to="/appointments" className={linkClass}>My Appointments</NavLink>
        <NavLink to="/availability" className={linkClass}>Manage Availability</NavLink>
        <button onClick={handleLogout} className={`${linkClass} text-left`}>Logout</button>
      </nav>
    </div>
  );
};

export default Sidebar;
