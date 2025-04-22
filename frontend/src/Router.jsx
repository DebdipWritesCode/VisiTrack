import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";

import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Dashboard from "./pages/Dashboard";
import BookAppointment from "./pages/BookAppointment";
import MyAppointments from "./pages/MyAppointments";
import ManageAvailability from "./pages/ManageAvailability";
import QRPage from "./pages/QRPage";

const AppRouter = () => {
  const userId = localStorage.getItem("user_id");

  return (
    <Routes>
      <Route path="/" element={userId ? <Navigate to="/dashboard" /> : <Navigate to="/login" />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<Signup />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/book" element={<BookAppointment />} />
      <Route path="/appointments" element={<MyAppointments />} />
      <Route path="/availability" element={<ManageAvailability />} />
      <Route path="/qr/:id" element={<QRPage />} />
    </Routes>
  );
};

export default AppRouter;
