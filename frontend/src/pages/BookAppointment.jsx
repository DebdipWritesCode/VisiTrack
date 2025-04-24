import React, { useEffect, useState } from "react";
import API from "../utils/api";
import { getUserId } from "../utils/auth";
import TimeSlotModal from "../components/TimeSlotModal";
import ConfirmationDialog from "../components/ConfirmationDialog";
import QRCodeModal from "../components/QRCodeModal";
import { debounce } from "lodash";
import { Bars } from "react-loader-spinner";

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
  const [loading, setLoading] = useState(false); // Loading state

  // Fetch users when the component loads
  useEffect(() => {
    const fetchPopularUsers = async () => {
      try {
        setLoading(true); // Set loading to true when fetching
        const res = await API.get("/users/popular");

        const filteredUsers = res.data.filter(
          (user) => user.role.String !== "admin" && user.id !== parseInt(getUserId())
        );

        setUsers(filteredUsers);
      } catch (err) {
        console.error("Error fetching popular users:", err);
      } finally {
        setLoading(false); // Set loading to false after fetching is done
      }
    };

    fetchPopularUsers();
  }, []);

  const debouncedSearch = debounce(async (query) => {
    try {
      setLoading(true); // Set loading to true when searching
      const res = await API.get(`/users/search?query=${query}`);

      const filteredUsers = res.data.filter(
        (user) => user.role.String !== "admin" && user.id !== parseInt(getUserId())
      );

      setUsers(filteredUsers);
    } catch (err) {
      console.error("Error searching users:", err);
    } finally {
      setLoading(false); // Set loading to false after search is done
    }
  }, 100);

  const handleSearchChange = (e) => {
    const value = e.target.value;
    setSearch(value);

    setLoading(true); // Show the loading spinner when typing
    debouncedSearch(value);
  };

  const handleBookClick = async (userId) => {
    try {
      setSelectedUserId(userId);
      setLoading(true); // Set loading when fetching availability
      const res = await API.get(`/availability/${userId}`);
      setAvailableSlots(res.data);
      setShowTimeModal(true);
    } catch (err) {
      console.error("Failed to fetch availability:", err);
    } finally {
      setLoading(false); // Stop loading after availability is fetched
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

      setQrCodeValue(res.data.qr_code.String);
      setShowConfirmDialog(false);
      setShowQrModal(true); 
      setSelectedUserId(null);
      setSelectedSlot(null);
    } catch (err) {
      console.error("Booking failed:", err);
    }
  };

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-6 text-gray-800 border-b pb-3">Book an Appointment</h2>

      <div className="relative mb-6">
        <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
          <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </div>
        <input
          type="text"
          placeholder="Search by name..."
          className="pl-10 border border-gray-300 rounded-lg px-4 py-3 w-full focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 shadow-sm"
          value={search}
          onChange={handleSearchChange}
        />
      </div>

      <div className="space-y-4">
        {loading ? (
          <div className="text-center py-8 bg-gray-50 rounded-lg border border-gray-200 flex justify-center items-center">
            <Bars width="30" color="#4B8BBE" />
          </div>
        ) : users.length === 0 ? (
          <div className="text-center py-8 bg-gray-50 rounded-lg border border-gray-200">
            <p className="text-gray-500">No users found matching your search criteria</p>
          </div>
        ) : (
          users.map((user) => (
            <div
              key={user.id}
              className="flex items-center justify-between p-5 bg-white rounded-lg shadow-sm border border-gray-100 hover:shadow-md transition-shadow duration-200"
            >
              <div className="flex items-center space-x-4">
                <div className="bg-indigo-100 text-indigo-700 rounded-full w-10 h-10 flex items-center justify-center font-medium">
                  {user.first_name.charAt(0)}{user.last_name.charAt(0)}
                </div>
                <div>
                  <p className="font-semibold text-gray-800">
                    {user.first_name} {user.last_name}
                  </p>
                  <div className="flex items-center text-sm text-gray-500 mt-1">
                    <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
                    </svg>
                    {user.phone_number}
                  </div>
                </div>
              </div>
              <button
                onClick={() => handleBookClick(user.id)}
                className="bg-indigo-600 text-white px-5 py-2 rounded-lg hover:bg-indigo-700 focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 flex items-center"
              >
                <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
                Book
              </button>
            </div>
          ))
        )}
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
