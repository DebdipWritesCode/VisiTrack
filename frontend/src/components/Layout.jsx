import React from "react";
import Sidebar from "./Sidebar";
import { Outlet } from "react-router-dom"; // Add this

const Layout = () => {
  return (
    <div className="flex">
      <Sidebar />
      <main className="ml-60 p-6 w-full">
        <Outlet /> {/* This will render the content of the nested route */}
      </main>
    </div>
  );
};

export default Layout;
