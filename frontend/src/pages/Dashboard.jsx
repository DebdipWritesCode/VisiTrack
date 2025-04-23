// import React, { useEffect, useState } from "react";
// import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";
// import API from "../utils/api"; // your axios instance
// import { getUserId } from "../utils/auth"; // localStorage utility

// const Dashboard = () => {
//   const [appointmentStats, setAppointmentStats] = useState({
//     hosted: 0,
//     visiting: 0,
//   });

//   useEffect(() => {
//     const fetchDashboardData = async () => {
//       try {
//         const userId = getUserId();

//         const hostedRes = await API.get(`/appointments/hosted/${userId}`);
//         const visitingRes = await API.get(`/appointments/visiting/${userId}`);

//         setAppointmentStats({
//           hosted: hostedRes.data.pending || 0,
//           visiting: visitingRes.data.pending || 0,
//         });
//       } catch (err) {
//         console.error("Error fetching dashboard data:", err);
//       }
//     };

//     fetchDashboardData();
//   }, []);

//   const chartData = [
//     { name: "Hosting", count: appointmentStats.hosted },
//     { name: "Visiting", count: appointmentStats.visiting },
//   ];

//   return (
//     <div className="p-6 w-full">
//       <h2 className="text-2xl font-bold mb-4">Dashboard</h2>
//       <div className="bg-white p-4 rounded shadow-md">
//         <h3 className="text-lg font-semibold mb-2">Pending Appointments</h3>
//         <ResponsiveContainer width="100%" height={300}>
//           <BarChart data={chartData}>
//             <XAxis dataKey="name" />
//             <YAxis />
//             <Tooltip />
//             <Bar dataKey="count" fill="#3b82f6" radius={[4, 4, 0, 0]} />
//           </BarChart>
//         </ResponsiveContainer>
//       </div>
//     </div>
//   );
// };

// export default Dashboard;
import React from 'react'

const Dashboard = () => {
  return (
    <div>Dashboard</div>
  )
}

export default Dashboard