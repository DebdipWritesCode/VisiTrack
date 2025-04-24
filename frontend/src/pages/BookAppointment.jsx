import React, { useEffect, useState } from "react";
import API from "../utils/api";
import { getUserId } from "../utils/auth";
import TimeSlotModal from "../components/TimeSlotModal";
import ConfirmationDialog from "../components/ConfirmationDialog";
import QRCodeModal from "../components/QRCodeModal";

const BookAppointment = () => {
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState("");
  const [selectedUserId, setSelectedUserId] = useState(null);
  const [showTimeModal, setShowTimeModal] = useState(false);
  const [availableSlots, setAvailableSlots] = useState([]);
  const [selectedSlot, setSelectedSlot] = useState(null);
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);
  const [qrCodeValue, setQrCodeValue] = useState(null);
  const [showQrModal, setShowQrModal] = useState(false);

  useEffect(() => {
    const fetchPopularUsers = async () => {
      try {
        const res = await API.get("/users/popular");

        const filteredUsers = res.data.filter(
          (user) => user.role.String !== "admin" && user.id !== parseInt(getUserId())
        )

        setUsers(filteredUsers);
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

  const handleBookClick = async (userId) => {
    try {
      setSelectedUserId(userId);
      const res = await API.get(`/availability/${userId}`);
      setAvailableSlots(res.data); // assuming array of strings like ["09:00", "09:30"]
      setShowTimeModal(true);
    } catch (err) {
      console.error("Failed to fetch availability:", err);
    }
  };

  const handleSlotSelect = (slot) => {
    setSelectedSlot(slot);
    setShowConfirmDialog(true);
  };

  const handleConfirmBooking = async () => {
    try {
      const visitorId = getUserId();

      const payload = {
        visitor_id: parseInt(visitorId),
        host_id: selectedUserId,
        appointment_date: selectedSlot.date,
        start_time: selectedSlot.start_time,
        end_time: selectedSlot.end_time,
        qr_code: `appointment-${Date.now()}`
      };

      const res = await API.post("/appointments", payload);

      setQrCodeValue(res.data.qr_code.String); // If backend sends it
      setShowConfirmDialog(false);
      setShowQrModal(true); // Show QR modal
      setSelectedUserId(null);
      setSelectedSlot(null);
    } catch (err) {
      console.error("Booking failed:", err);
    }
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

      <TimeSlotModal
        isOpen={showTimeModal}
        onClose={() => setShowTimeModal(false)}
        timeSlots={availableSlots}
        onSelect={handleSlotSelect}
      />

      <ConfirmationDialog
        isOpen={showConfirmDialog}
        onConfirm={handleConfirmBooking}
        onCancel={() => setShowConfirmDialog(false)}
      />

      <QRCodeModal
        isOpen={showQrModal}
        onClose={() => setShowQrModal(false)}
        value={qrCodeValue}
      />
    </div>
  );
};

export default BookAppointment;
