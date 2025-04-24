import React from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { logoutUser } from "../utils/auth";

const Sidebar = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    logoutUser();
    navigate("/login");
  };

  const navLinkClass = ({ isActive }) =>
    `flex items-center px-4 py-3 rounded-lg transition duration-200 ${
      isActive
        ? "bg-indigo-600 text-white font-medium shadow-md"
        : "text-indigo-100 hover:bg-indigo-700 hover:text-white hover:shadow-md"
    }`;

  return (
    <div className="h-screen w-64 bg-gradient-to-b from-indigo-900 to-indigo-500 p-6 fixed top-0 left-0 shadow-xl flex flex-col">
      <div className="mb-10">
        <h2 className="text-white text-2xl font-bold">
          <span className="border-b-2 border-indigo-300 pb-1">IIITN</span> 
          <span className="text-indigo-200"> Visitors</span>
        </h2>
      </div>
      
      <nav className="flex flex-col space-y-4 flex-grow">
        <NavLink to="/dashboard" className={navLinkClass}>
          Dashboard
        </NavLink>
        <NavLink to="/book" className={navLinkClass}>
          Book Appointment
        </NavLink>
        <NavLink to="/appointments" className={navLinkClass}>
          My Appointments
        </NavLink>
        <NavLink to="/availability" className={navLinkClass}>
          Manage Availability
        </NavLink>
        
        <div className="flex-grow"></div>
        
        <button 
          onClick={handleLogout} 
          className="flex items-center px-4 py-3 mt-8 rounded-lg text-white hover:bg-rose-600 transition duration-200 text-left bg-red-500"
        >
          Logout
        </button>
      </nav>
    </div>
  );
};

export default Sidebar;