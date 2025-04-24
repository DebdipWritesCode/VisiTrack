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
    <div className="p-6 max-w-4xl mx-auto bg-gray-50 rounded-lg shadow-sm">
      <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">My Appointments</h2>

      <div className="flex justify-center mb-8">
        <div className="bg-white p-1 rounded-lg shadow-sm inline-flex space-x-1">
          <button
            className={`px-6 py-3 rounded-md font-medium transition-all duration-200 ${activeTab === "visited"
                ? "bg-indigo-600 text-white shadow-md"
                : "bg-white text-gray-700 hover:bg-gray-100"
              }`}
            onClick={() => setActiveTab("visited")}
          >
            Visited
          </button>
          <button
            className={`px-6 py-3 rounded-md font-medium transition-all duration-200 ${activeTab === "hosted"
                ? "bg-indigo-600 text-white shadow-md"
                : "bg-white text-gray-700 hover:bg-gray-100"
              }`}
            onClick={() => setActiveTab("hosted")}
          >
            Hosted
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex justify-center items-center h-40">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          <span className="ml-3 text-gray-600">Loading appointments...</span>
        </div>
      ) : (
        <div className="space-y-4">
          {appointments.map((appointment) => (
            <div
              key={appointment.id}
              className="p-5 bg-white border border-gray-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200"
            >
              <div className="flex justify-between items-start">
                <div className="space-y-2">
                  <div className="flex items-center space-x-2">
                    <div className={`w-2 h-2 rounded-full ${appointment.status.String === "pending" ? "bg-yellow-400" :
                        appointment.status.String === "approved" ? "bg-green-500" :
                          appointment.status.String === "cancelled" ? "bg-red-500" : "bg-gray-400"
                      }`}></div>
                    <p className="font-semibold text-lg text-gray-800">
                      {activeTab === "visited"
                        ? `Host: ${appointment.host_name}`
                        : `Visitor: ${appointment.visitor_name}`}
                    </p>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-gray-600">
                    <div className="flex items-center">
                      <svg className="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                      </svg>
                      <span>{new Date(appointment.appointment_date).toLocaleDateString('en-US', {
                        weekday: 'short',
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric'
                      })}</span>
                    </div>

                    <div className="flex items-center">
                      <svg className="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                      </svg>
                      <span>{formatTime(appointment.start_time)} - {formatTime(appointment.end_time)}</span>
                    </div>

                    <div className="flex items-center md:col-span-2">
                      <svg className="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                      </svg>
                      <span className={`${appointment.status.String === "pending" ? "text-yellow-600" :
                          appointment.status.String === "approved" ? "text-green-600" :
                            appointment.status.String === "cancelled" ? "text-red-600" : "text-gray-600"
                        } font-medium`}>
                        {appointment.status.String.charAt(0).toUpperCase() + appointment.status.String.slice(1)}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="mt-2">
                  {activeTab === "hosted" && appointment.status.String === "pending" && (
                    <button
                      onClick={() => handleCancel(appointment.id)}
                      className="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-md transition-colors duration-200 inline-flex items-center"
                    >
                      <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
                      </svg>
                      Cancel
                    </button>
                  )}

                  {activeTab === "visited" && (
                    <button
                      onClick={() => handleViewQR(appointment.qr_code.String)}
                      className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md transition-colors duration-200 inline-flex items-center"
                    >
                      <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"></path>
                      </svg>
                      See QR
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}

          {appointments.length === 0 && (
            <div className="bg-white border border-gray-200 rounded-lg p-10 text-center">
              <svg className="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <p className="text-gray-500 text-lg">No appointments found.</p>
              <p className="text-gray-400 mt-2">Check back later or create a new appointment.</p>
            </div>
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
