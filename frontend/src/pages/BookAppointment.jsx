import React, { useEffect, useState } from "react";
import API from "../utils/api";
import { getUserId } from "../utils/auth";

const BookAppointment = () => {
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState("");

  useEffect(() => {
    const fetchPopularUsers = async () => {
      try {
        const res = await API.get("/users/popular");
        setUsers(res.data);
      } catch (err) {
        console.error("Error fetching popular users:", err);
      }
    };

    fetchPopularUsers();
  }, []);

  const handleSearchChange = (e) => {
    setSearch(e.target.value);
  };

  const filteredUsers = users.filter((user) =>
    `${user.first_name} ${user.last_name}`
      .toLowerCase()
      .includes(search.toLowerCase())
  );

  const handleBookClick = (userId) => {
    console.log("Book clicked for user:", userId);
    // TODO: Show timeslot picker modal
  };

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4">Book an Appointment</h2>

      <input
        type="text"
        placeholder="Search by name..."
        className="border rounded px-4 py-2 mb-4 w-full"
        value={search}
        onChange={handleSearchChange}
      />

      <div className="space-y-4">
        {filteredUsers.map((user) => (
          <div
            key={user.id}
            className="flex items-center justify-between p-4 bg-white rounded shadow"
          >
            <div>
              <p className="font-semibold">
                {user.first_name} {user.last_name}
              </p>
              <p className="text-sm text-gray-500">Phone: {user.phone_number}</p>
            </div>
            <button
              onClick={() => handleBookClick(user.id)}
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            >
              Book
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default BookAppointment;
