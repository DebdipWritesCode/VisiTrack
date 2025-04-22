import React from "react";
import Sidebar from "./Sidebar";

const Layout = ({ children }) => {
  return (
    <div className="flex">
      <Sidebar />
      <main className="ml-60 p-6 w-full">{children}</main>
    </div>
  );
};

export default Layout;
