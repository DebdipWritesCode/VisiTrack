import React, { useState, useEffect, Fragment } from "react";
import { Dialog, Transition } from "@headlessui/react";
import Calendar from "react-calendar";
import 'react-calendar/dist/Calendar.css';
import { toast, ToastContainer } from "react-toastify";
import API from "../utils/api";
import { QrReader } from "react-qr-reader";
import { Bars } from "react-loader-spinner";

const AdminPage = () => {
  const [date, setDate] = useState(new Date());
  const [appointments, setAppointments] = useState([]);
  const [selectedAppointment, setSelectedAppointment] = useState(null);
  const [qrScanActive, setQrScanActive] = useState(false);
  const [scannedResult, setScannedResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [updatingStatus, setUpdatingStatus] = useState(false);

  const formatTime = (h) => {
    const tm = h.split("T")[1].split(":");
    return `${tm[0]}:${tm[1]}`;
  }

  useEffect(() => {
    const fetchAppointments = async () => {
      setLoading(true);
      try {
        const res = await API.get(`/appointments/date`, {
          params: {
            date: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
          }
        });
        setAppointments(res.data);
      } catch (err) {
        console.error(err);
        toast.error("Failed to fetch appointments");
      } finally {
        setLoading(false);
      }
    };
    fetchAppointments();
  }, [date]);

  const handleScan = async (qrData) => {
    if (!qrData) return;

    try {
      const response = await API.get(`/appointments/qr/${qrData}`);
      const appointment = response.data;

      setSelectedAppointment(appointment);
      setQrScanActive(false);
      toast.success("Appointment verified");
    } catch (error) {
      if (error.response?.status === 404) {
        toast.error("No appointment found for this QR.");
      } else {
        toast.error("Error verifying appointment. Please try again.");
      }
    }
  };

  const updateStatus = async (status) => {
    setUpdatingStatus(true);
    try {
      const payload = {
        id: selectedAppointment.id,
        status,
      };
  
      await API.put("/appointments/status", payload);
  
      toast.success(`Marked as ${status}`);
  
      const updatedAppointments = appointments.map(app =>
        app.id === selectedAppointment.id
          ? { ...app, status: { String: status.toLowerCase() } }
          : app
      );
      setAppointments(updatedAppointments);
      setSelectedAppointment(null);
    } catch (err) {
      toast.error("Failed to update status");
    } finally {
      setUpdatingStatus(false);
    }
  };  

  useEffect(() => {
    if (scannedResult) {
      handleScan(scannedResult);
    }
  }, [scannedResult]);

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Admin Dashboard</h1>
        <button
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          onClick={() => setQrScanActive(true)}
        >
          Scan QR
        </button>
      </div>

      <Calendar onChange={setDate} value={date} />

      {qrScanActive && (
        <Transition appear show={qrScanActive} as={Fragment}>
          <Dialog as="div" className="relative z-10" onClose={() => setQrScanActive(false)}>
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300" leave="ease-in duration-200"
              enterFrom="opacity-0" enterTo="opacity-100"
              leaveFrom="opacity-100" leaveTo="opacity-0"
            >
              <div className="fixed inset-0 bg-black opacity-50" />
            </Transition.Child>

            <div className="fixed inset-0 overflow-y-auto">
              <div className="flex min-h-full items-center justify-center p-4">
                <Transition.Child
                  as={Fragment}
                  enter="ease-out duration-300" leave="ease-in duration-200"
                  enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100"
                  leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95"
                >
                  <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white shadow-xl transition-all">
                    <Dialog.Title className="text-xl font-medium text-gray-900 p-4 border-b">
                      {scannedResult ? 'Scanned Result' : 'Scan QR Code'}
                    </Dialog.Title>

                    <div className="w-full relative h-64 md:h-80 flex items-center justify-center">
                      {!scannedResult ? (
                        <QrReader
                          constraints={{
                            facingMode: 'environment',
                            width: { min: 640, ideal: 720, max: 1280 },
                            height: { min: 480, ideal: 720, max: 720 },
                            aspectRatio: 1.3333333
                          }}
                          videoId="qr-video"
                          scanDelay={500}
                          onResult={(result, error) => {
                            if (result?.text) {
                              setScannedResult(result.text);
                            }
                            if (error) {
                              console.error(error);
                            }
                          }}
                          containerStyle={{ width: '100%', height: '100%' }}
                          videoStyle={{ width: '100%', height: '70%', objectFit: 'cover' }}
                        />
                      ) : (
                        <p className="text-center text-blue-600 underline break-words px-4">
                          <a href={scannedResult} target="_blank" rel="noopener noreferrer">
                            {scannedResult}
                          </a>
                        </p>
                      )}
                    </div>

                    <div className="p-4 border-t flex justify-between items-center">
                      {scannedResult && (
                        <button
                          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                          onClick={() => setScannedResult(null)}
                        >
                          Rescan
                        </button>
                      )}
                      <button
                        className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 ml-auto"
                        onClick={() => {
                          setScannedResult(null);
                          setQrScanActive(false);
                        }}
                      >
                        Cancel
                      </button>
                    </div>
                  </Dialog.Panel>
                </Transition.Child>
              </div>
            </div>
          </Dialog>
        </Transition>
      )}

      {/* Show loading spinner when fetching appointments */}
      {loading && (
        <div className="flex justify-center items-center mt-6">
          <Bars
            height="80"
            width="80"
            color="#007bff"
            ariaLabel="bars-loading"
            wrapperClass="flex justify-center items-center"
          />
        </div>
      )}

      {!loading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
          {appointments.length === 0 ? (
            <div className="col-span-full text-center py-10 text-gray-500">
              No appointments found for this date.
            </div>
          ) : (
            appointments.map((appt) => (
              <div
                key={appt.id}
                className="bg-white rounded-lg shadow-md overflow-hidden border-l-4 border-blue-500 hover:shadow-lg transition-shadow duration-300 cursor-pointer"
                onClick={() => setSelectedAppointment(appt)}
              >
                <div className="p-4">
                  <div className="flex justify-between items-center mb-3">
                    <span className={`px-2 py-1 rounded-full text-xs font-semibold
                      ${appt.status.String === 'confirmed' ? 'bg-green-100 text-green-800' :
                        appt.status.String === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                          appt.status.String === 'ongoing' ? 'bg-blue-100 text-blue-800' :
                            appt.status.String === 'completed' ? 'bg-purple-100 text-purple-800' :
                              appt.status.String === 'cancelled' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-800'}`}>
                      {appt.status.String}
                    </span>
                    <span className="text-xs text-gray-500">ID: {appt.id}</span>
                  </div>

                  <div className="mb-4">
                    <div className="flex items-center mb-2">
                      <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center mr-2">
                        <span className="text-blue-600 font-bold">H</span>
                      </div>
                      <span className="font-medium text-gray-900">{appt.host_name}</span>
                    </div>

                    <div className="flex items-center">
                      <div className="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center mr-2">
                        <span className="text-purple-600 font-bold">V</span>
                      </div>
                      <span className="font-medium text-gray-900">{appt.visitor_name}</span>
                    </div>
                  </div>

                  <div className="border-t border-gray-100 pt-3">
                    <div className="flex items-center mb-2">
                      <svg className="w-4 h-4 text-gray-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                      </svg>
                      <span className="text-sm text-gray-700">
                        {new Date(appt.appointment_date).toLocaleDateString(undefined, {
                          weekday: 'short',
                          year: 'numeric',
                          month: 'short',
                          day: 'numeric'
                        })}
                      </span>
                    </div>
                    <div className="flex items-center">
                      <svg className="w-4 h-4 text-gray-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                      </svg>
                      <span className="text-sm text-gray-700">
                        {formatTime(appt.start_time)} - {formatTime(appt.end_time)}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {/* Appointment Detail Modal */}
      <Transition appear show={!!selectedAppointment} as={Fragment}>
        <Dialog as="div" className="relative z-10" onClose={() => setSelectedAppointment(null)}>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300" leave="ease-in duration-200"
            enterFrom="opacity-0" enterTo="opacity-100"
            leaveFrom="opacity-100" leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-black opacity-25" />
          </Transition.Child>

          <div className="fixed inset-0 overflow-y-auto">
            <div className="flex min-h-full items-center justify-center p-4">
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300" leave="ease-in duration-200"
                enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100"
                leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95"
              >
                <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 shadow-xl transition-all">
                  <Dialog.Title className="text-xl font-medium text-gray-900 mb-4">
                    Appointment Details
                  </Dialog.Title>

                  {selectedAppointment && (
                    <>
                      <div className="space-y-3 mb-6">
                        <div className="flex items-center justify-between">
                          <span className="font-semibold">Status:</span>
                          <span className={`px-2 py-1 rounded-full text-xs font-semibold
                            ${selectedAppointment.status.String === 'confirmed' ? 'bg-green-100 text-green-800' :
                              selectedAppointment.status.String === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                selectedAppointment.status.String === 'ongoing' ? 'bg-blue-100 text-blue-800' :
                                  selectedAppointment.status.String === 'completed' ? 'bg-purple-100 text-purple-800' :
                                    selectedAppointment.status.String === 'cancelled' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-800'}`}>
                            {selectedAppointment.status.String}
                          </span>
                        </div>

                        <div className="border-t border-gray-100 pt-2">
                          <p className="font-medium text-gray-700 mb-1">Host</p>
                          <p className="text-gray-900">{selectedAppointment.host_name}</p>
                        </div>

                        <div className="border-t border-gray-100 pt-2">
                          <p className="font-medium text-gray-700 mb-1">Visitor</p>
                          <p className="text-gray-900">{selectedAppointment.visitor_name}</p>
                        </div>

                        <div className="border-t border-gray-100 pt-2">
                          <p className="font-medium text-gray-700 mb-1">Date & Time</p>
                          <p className="text-gray-900">
                            {new Date(selectedAppointment.appointment_date).toLocaleDateString(undefined, {
                              weekday: 'long',
                              year: 'numeric',
                              month: 'long',
                              day: 'numeric'
                            })}
                          </p>
                          <p className="text-gray-900">
                            {formatTime(selectedAppointment.start_time)} - {formatTime(selectedAppointment.end_time)}
                          </p>
                        </div>

                        {selectedAppointment.qr_code && selectedAppointment.qr_code.String && (
                          <div className="border-t border-gray-100 pt-2">
                            <p className="font-medium text-gray-700 mb-1">QR Code</p>
                            <p className="text-gray-900 break-all">{selectedAppointment.qr_code.String}</p>
                          </div>
                        )}
                      </div>

                      <div className="border-t border-gray-200 pt-4">
                        <p className="font-medium text-gray-700 mb-2">Update Status</p>
                        <div className="grid grid-cols-2 gap-2">
                          {['pending', 'ongoing', 'completed', 'cancelled'].map((status) => (
                            <button
                              key={status}
                              onClick={() => updateStatus(status)}
                              disabled={updatingStatus || selectedAppointment.status.String === status}
                              className={`px-3 py-2 rounded text-sm font-medium transition-colors
                                ${selectedAppointment.status.String === status
                                  ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                                  : 'bg-gray-200 hover:bg-gray-300 text-gray-800'}`}
                            >
                              {updatingStatus ? (
                                <span className="flex items-center justify-center">
                                  <svg className="animate-spin h-4 w-4 mr-1" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                  </svg>
                                  Processing
                                </span>
                              ) : (
                                status.charAt(0).toUpperCase() + status.slice(1)
                              )}
                            </button>
                          ))}
                        </div>
                      </div>

                      <div className="mt-6 flex justify-end">
                        <button
                          className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                          onClick={() => setSelectedAppointment(null)}
                        >
                          Close
                        </button>
                      </div>
                    </>
                  )}
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </Dialog>
      </Transition>

      <ToastContainer position="top-right" autoClose={3000} />
    </div>
  );
};

export default AdminPage;