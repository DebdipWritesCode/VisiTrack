import React, { useEffect, useState } from 'react';
import API from '../utils/api';
import { getUserId } from '../utils/auth';

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
      alert("Availability updated successfully!");
    } catch (err) {
      console.error("Error saving availability", err);
      alert("Failed to save changes.");
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
    <div className="p-6">
      <h2 className="text-xl font-semibold mb-4">Manage Availability</h2>
      <div className="overflow-x-auto">
        <table className="table-auto border border-gray-400">
          <thead>
            <tr>
              <th className="border border-gray-400 px-2 py-1">Day / Time</th>
              {hours.map(hour => (
                <th key={hour} className="border border-gray-400 px-2 py-1">
                  {formatTime(hour)} - {formatTime(hour + 1)}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {weekdays.map((dayName, i) => {
              const day = i + 1;
              return (
                <tr key={day}>
                  <td className="border border-gray-400 px-2 py-1 font-semibold">{dayName}</td>
                  {hours.map(hour => {
                    const key = getSlotKey(day, hour);
                    const isAvailable = availability[key];
                    const slotStatus = status[key];
                    return (
                      <td
                        key={hour}
                        onClick={() => toggleSlot(day, hour)}
                        className={`cursor-pointer border px-3 py-2 text-center ${slotStatus === 'available' ? 'bg-green-400' : 'bg-red-400'} ${editMode ? 'hover:opacity-80' : ''}`}
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

      <div className="mt-4 space-x-2">
        {!editMode ? (
          <button
            className="bg-blue-600 text-white px-4 py-2 rounded cursor-pointer hover:bg-blue-700"
            onClick={() => setEditMode(true)}
          >
            Change Availability
          </button>
        ) : (
          <>
            <button
              className="bg-green-600 text-white px-4 py-2 rounded disabled:opacity-50 cursor-pointer hover:bg-green-700"
              onClick={handleSave}
              disabled={isSaving}
            >
              {isSaving ? 'Saving...' : 'Save Changes'}
            </button>
            <button
              className="bg-gray-500 text-white px-4 py-2 rounded"
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
    </div>
  );
};

export default ManageAvailability;
