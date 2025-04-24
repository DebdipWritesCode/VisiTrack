import React, { useEffect, useState } from 'react'
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import API from '../utils/api'
import { toast } from 'react-toastify'
import { getUserId } from '../utils/auth'

const COLORS = ['#4F46E5', '#10B981', '#F59E0B']

const Dashboard = () => {
  const id = getUserId()
  const [stats, setStats] = useState(null)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const res = await API.get(`/users/${id}/stats`)
        console.log('Stats response:', res.data)
        setStats(res.data)
      } catch (err) {
        toast.error("Failed to load stats")
      }
    }

    fetchStats()
  }, [id])

  if (!stats) return (
    <div className="flex items-center justify-center h-64 text-gray-600 text-lg">
      <div className="animate-pulse flex flex-col items-center">
        <div className="h-8 w-8 border-4 border-t-indigo-600 border-gray-200 rounded-full animate-spin mb-4"></div>
        Loading dashboard...
      </div>
    </div>
  )

  const pieData = [
    { name: 'Hosted', value: stats.appointments_hosted.Int32 },
    { name: 'Visited', value: stats.appointments_visited.Int32 },
  ]

  const barData = [
    {
      name: 'Appointments',
      Hosted: stats.appointments_hosted.Int32,
      Visited: stats.appointments_visited.Int32,
      Pending: stats.pending_appointments,
    },
  ]

  return (
    <div className="p-8 space-y-10 bg-gray-50 min-h-screen">
      <h1 className="text-3xl font-bold text-center text-gray-800 mb-2">
        User Appointment Dashboard
      </h1>
      <div className="max-w-6xl mx-auto">
        <div className="bg-indigo-100 text-indigo-800 p-4 rounded-lg text-center mb-8 shadow-sm">
          <p className="font-medium">
            Welcome! Here's an overview of your appointment activity.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100">
            <h2 className="text-xl font-semibold mb-6 text-center text-gray-700">Hosted vs Visited</h2>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pieData}
                  dataKey="value"
                  nameKey="name"
                  outerRadius={100}
                  label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                  labelLine={false}
                >
                  {pieData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [`${value} appointments`, null]} />
                <Legend verticalAlign="bottom" height={36} />
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100">
            <h2 className="text-xl font-semibold mb-6 text-center text-gray-700">Appointment Overview</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={barData}>
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip cursor={{ fill: 'rgba(0, 0, 0, 0.05)' }} />
                <Legend verticalAlign="bottom" height={36} />
                <Bar dataKey="Hosted" fill={COLORS[0]} radius={[4, 4, 0, 0]} />
                <Bar dataKey="Visited" fill={COLORS[1]} radius={[4, 4, 0, 0]} />
                <Bar dataKey="Pending" fill={COLORS[2]} radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-10">
          <div className="bg-white p-5 rounded-xl shadow-md border-l-4 border-indigo-500">
            <h3 className="text-lg font-medium text-gray-700">Hosted</h3>
            <p className="text-3xl font-bold text-indigo-600">{stats.appointments_hosted.Int32}</p>
          </div>
          
          <div className="bg-white p-5 rounded-xl shadow-md border-l-4 border-emerald-500">
            <h3 className="text-lg font-medium text-gray-700">Visited</h3>
            <p className="text-3xl font-bold text-emerald-600">{stats.appointments_visited.Int32}</p>
          </div>
          
          <div className="bg-white p-5 rounded-xl shadow-md border-l-4 border-amber-500">
            <h3 className="text-lg font-medium text-gray-700">Pending</h3>
            <p className="text-3xl font-bold text-amber-600">{stats.pending_appointments}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard