import React, { useEffect, useState } from 'react';
import API from '../utils/api';
import { getUserId } from '../utils/auth';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const weekdays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];
const hours = Array.from({ length: 9 }, (_, i) => 9 + i); // 9 to 17

const formatTime = (h) => `${String(h).padStart(2, '0')}:${String(0).padStart(2, '0')}`;
const getSlotKey = (day, hour) => `${day}-${hour}`;

const ManageAvailability = () => {
  const [availability, setAvailability] = useState({});
  const [editMode, setEditMode] = useState(false);
  const [initialAvailability, setInitialAvailability] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [status, setStatus] = useState({});

  const userId = getUserId();

  useEffect(() => {
    API.get(`/availability/${userId}`)
      .then(res => {
        const availMap = {};
        const statusMap = {};
        res.data.forEach(slot => {
          const day = slot.day_of_week;
          const hour = new Date(slot.start_time).getUTCHours();
          const key = getSlotKey(day, hour);
          availMap[key] = true;
          statusMap[key] = slot.status.String || 'available';
        });
        setAvailability({ ...availMap });
        setInitialAvailability({ ...availMap });
        setStatus({ ...statusMap });
        setIsLoading(false);
      })
      .catch(err => {
        console.error("Failed to fetch availability", err);
        setIsLoading(false);
      });
  }, [userId]);

  const toggleSlot = (day, hour) => {
    if (!editMode) return;
    const key = getSlotKey(day, hour);
    const newStatus = status[key] === 'available' ? 'not_available' : 'available';

    setAvailability(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
    setStatus(prev => ({
      ...prev,
      [key]: newStatus
    }));
  };

  const handleSave = async () => {
    setIsSaving(true);
    const promises = [];

    for (let day = 1; day <= 5; day++) {
      for (let hour = 9; hour <= 17; hour++) {
        const key = getSlotKey(day, hour);
        const was = initialAvailability[key] || false;
        const now = availability[key] || false;
        const slotStatus = status[key] || 'available';

        if (was !== now) {
          const startTime = formatTime(hour);
          const endTime = formatTime(hour + 1);

          const requestPayload = {
            user_id: parseInt(userId),
            day_of_week: day,
            start_time: startTime,
            end_time: endTime,
            status: slotStatus
          };

          console.log("Sending request payload:", requestPayload);

          promises.push(
            API.put('/availability/status', requestPayload)
              .catch(err => {
                if (err.response?.status === 409) {
                  console.warn('Slot already exists:', key);
                } else {
                  throw err;
                }
              })
          );
        }
      }
    }

    try {
      await Promise.all(promises);
      setInitialAvailability({ ...availability });
      setEditMode(false);
      toast.success("Availability updated successfully!");
    } catch (err) {
      console.error("Error saving availability", err);
      toast.error("Failed to update availability. Please try again.");
    } finally {
      setIsSaving(false);
    }
  };


  if (isLoading) {
    return (
      <div className="p-6 text-center text-lg font-medium">
        Loading availability...
      </div>
    );
  }

  return (
    <div className="p-6 bg-white rounded-xl shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-gray-800 text-center">Manage Availability</h2>

      <div className="overflow-x-auto rounded-md border border-gray-200">
        <table className="min-w-full text-sm text-gray-700 text-center">
          <thead className="bg-gray-100 text-gray-700 font-semibold">
            <tr>
              <th className="border border-gray-300 px-4 py-3">Day / Time</th>
              {hours.map(hour => (
                <th key={hour} className="border border-gray-300 px-4 py-3">
                  {formatTime(hour)} - {formatTime(hour + 1)}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {weekdays.map((dayName, i) => {
              const day = i + 1;
              return (
                <tr key={day} className="even:bg-gray-50">
                  <td className="border border-gray-300 px-4 py-2 font-medium">{dayName}</td>
                  {hours.map(hour => {
                    const key = getSlotKey(day, hour);
                    const isAvailable = availability[key];
                    const slotStatus = status[key];
                    return (
                      <td
                        key={hour}
                        onClick={() => toggleSlot(day, hour)}
                        className={`border px-4 py-2 transition-colors duration-200 ease-in-out ${slotStatus === 'available'
                            ? 'bg-green-100 text-green-700'
                            : 'bg-red-100 text-red-700'
                          } ${editMode ? 'cursor-pointer hover:opacity-80' : 'cursor-default'}`}
                      >
                        {slotStatus === 'available' ? '✔️' : '❌'}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <div className="mt-6 flex justify-center space-x-4">
        {!editMode ? (
          <button
            className="bg-blue-600 text-white px-6 py-2 rounded-lg font-medium shadow hover:bg-blue-700 transition"
            onClick={() => setEditMode(true)}
          >
            Change Availability
          </button>
        ) : (
          <>
            <button
              className="bg-green-600 text-white px-6 py-2 rounded-lg font-medium shadow hover:bg-green-700 transition disabled:opacity-50"
              onClick={handleSave}
              disabled={isSaving}
            >
              {isSaving ? 'Saving...' : 'Save Changes'}
            </button>
            <button
              className="bg-gray-400 text-white px-6 py-2 rounded-lg font-medium shadow hover:bg-gray-500 transition"
              onClick={() => {
                setAvailability({ ...initialAvailability });
                setEditMode(false);
              }}
              disabled={isSaving}
            >
              Cancel
            </button>
          </>
        )}
      </div>

      <ToastContainer position="top-right" autoClose={3000} />
    </div>
  );
};

export default ManageAvailability;
