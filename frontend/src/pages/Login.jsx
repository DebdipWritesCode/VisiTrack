import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import API from "../utils/api";
import { saveUserId } from "../utils/auth";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const Login = () => {
  const [phoneNumber, setPhoneNumber] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const res = await API.post("/auth/login", {
        phone_number: phoneNumber,
      });
      
      toast.success("Login successful!");
      saveUserId(res.data.id);
  
      if (res.data.role.String === "admin") {
        navigate("/admin");
      } else {
        navigate("/dashboard");
      }
    } catch (error) {
      console.error("Login error:", error);
      if (error.response && error.response.status === 401) {
        toast.error("Invalid phone number or password. Please try again.");
      } else {
        toast.error("An error occurred. Please try again later.");
      }
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-blue-100">
      <div className="bg-white p-8 rounded-lg shadow-lg w-96">
        <h2 className="text-2xl font-bold text-blue-700 mb-6 text-center">Login</h2>
        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Phone Number</label>
            <input
              type="tel"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
              className="w-full mt-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <button
            type="submit"
            className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
          >
            Login
          </button>
        </form>

        <p className="text-sm mt-4 text-center">
          New here?{" "}
          <Link to="/signup" className="text-blue-600 hover:underline">
            Create an account
          </Link>
        </p>
      </div>

      <ToastContainer position="top-right" autoClose={3000} />
    </div>
  );
};

export default Login;
