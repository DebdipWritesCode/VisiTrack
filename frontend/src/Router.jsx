import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";

import Login from "./pages/Login";
import Signup from "./pages/Signup";

import Dashboard from "./pages/Dashboard";
import BookAppointment from "./pages/BookAppointment";
import MyAppointments from "./pages/MyAppointments";
import ManageAvailability from "./pages/ManageAvailability";
import QRPage from "./pages/QRPage";

import Layout from "./components/Layout";
import { getUserId } from "./utils/auth";

const AppRouter = () => {
  const userId = getUserId();

  return (
    <Routes>
      <Route path="/" element={userId ? <Navigate to="/dashboard" /> : <Navigate to="/login" />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<Signup />} />

      <Route element={<Layout />}>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/book" element={<BookAppointment />} />
        <Route path="/appointments" element={<MyAppointments />} />
        <Route path="/availability" element={<ManageAvailability />} />
      </Route>

      <Route path="/qr/:id" element={<QRPage />} /> {/* Optional sidebar */}
    </Routes>
  );
};

export default AppRouter;
