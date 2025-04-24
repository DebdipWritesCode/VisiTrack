import React, { useEffect, useState } from "react";
import API from "../utils/api";
import { getUserId } from "../utils/auth";
import QRCodeModal from "../components/QRCodeModal"; // make sure the path is correct

const MyAppointments = () => {
  const [activeTab, setActiveTab] = useState("visited");
  const [appointments, setAppointments] = useState([]);
  const [loading, setLoading] = useState(false);
  const [qrValue, setQrValue] = useState("");
  const [showQrModal, setShowQrModal] = useState(false);

  const formatTime = (h) => {
    const tm = h.split("T")[1].split(":");
    return `${tm[0]}:${tm[1]}`;
  }

  const userId = getUserId();

  const fetchAppointments = async () => {
    if (!userId) return;

    setLoading(true);
    try {
      const endpoint =
        activeTab === "visited"
          ? `/appointments/visitor/${userId}`
          : `/appointments/host/${userId}`;
      const res = await API.get(endpoint);
      setAppointments(res.data);
      console.log("Appointments fetched:", res.data);
    } catch (err) {
      console.error("Failed to fetch appointments:", err);
    }
    setLoading(false);
  };

  const handleCancel = async (id) => {
    try {
      await API.post(`/appointments/${id}/cancel`);
      fetchAppointments();
    } catch (err) {
      console.error("Failed to cancel appointment:", err);
    }
  };

  const handleViewQR = (qrUrl) => {
    setQrValue(qrUrl);         // Set the QR code value
    setShowQrModal(true);      // Open modal
  };

  useEffect(() => {
    fetchAppointments();
  }, [activeTab]);

  return (
    <div className="p-4 max-w-3xl mx-auto">
      <div className="flex justify-center mb-6 space-x-4">
        <button
          className={`px-4 py-2 rounded cursor-pointer ${activeTab === "visited" ? "bg-blue-600 text-white" : "bg-gray-200"}`}
          onClick={() => setActiveTab("visited")}
        >
          Visited
        </button>
        <button
          className={`px-4 py-2 rounded cursor-pointer ${activeTab === "hosted" ? "bg-blue-600 text-white" : "bg-gray-200"}`}
          onClick={() => setActiveTab("hosted")}
        >
          Hosted
        </button>
      </div>

      {loading ? (
        <p className="text-center text-gray-500">Loading...</p>
      ) : (
        <div className="space-y-4">
          {appointments.map((appointment) => (
            <div
              key={appointment.id}
              className="p-4 border rounded shadow-sm flex justify-between items-center"
            >
              <div>
                <p className="font-semibold">
                  {activeTab === "visited"
                    ? `Host: ${appointment.host_name}`
                    : `Visitor: ${appointment.visitor_name}`}
                </p>
                <p>
                  Date: {new Date(appointment.appointment_date).toLocaleDateString()}
                </p>
                <p>
                  Time: {formatTime(appointment.start_time)} -{" "}
                  {formatTime(appointment.end_time)}
                </p>
                <p>Status: {appointment.status.String}</p>
              </div>

              {activeTab === "hosted" && appointment.status.String === "pending" && (
                <button
                  onClick={() => handleCancel(appointment.id)}
                  className="px-3 py-1 bg-red-500 text-white rounded"
                >
                  Cancel
                </button>
              )}

              {activeTab === "visited" && (
                <button
                  onClick={() => handleViewQR(appointment.qr_code.String)}
                  className="px-3 py-1 bg-green-600 text-white rounded"
                >
                  See QR
                </button>
              )}
            </div>
          ))}

          {appointments.length === 0 && (
            <p className="text-center text-gray-500">No appointments found.</p>
          )}
        </div>
      )}

      <QRCodeModal
        isOpen={showQrModal}
        value={qrValue}
        onClose={() => setShowQrModal(false)}
      />
    </div>
  );
};

export default MyAppointments;
