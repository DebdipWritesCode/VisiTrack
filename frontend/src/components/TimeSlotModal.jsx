import React, { useState } from "react";
import { Dialog, DialogPanel, DialogTitle } from "@headlessui/react";

const TimeSlotModal = ({ isOpen, onClose, timeSlots, onSelect }) => {
  const [selectedSlot, setSelectedSlot] = useState(null);
  const [selectedDate, setSelectedDate] = useState(() => new Date().toISOString().split("T")[0]);

  const formatTime = (h) => `${String(h).padStart(2, "0")}:00`;

  const extractHour = (isoString) => {
    const date = new Date(isoString);
    return date.getUTCHours(); // or getHours() if working with local time
  };

  const formatDay = (dayNum) => {
    const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    return days[dayNum] || "Invalid";
  };

  const handleConfirm = () => {
    if (selectedSlot) {
      const combineDateTime = (date, timeStr) =>
        new Date(`${date}T${timeStr}:00Z`).toISOString();
  
      onSelect({
        ...selectedSlot,
        date: new Date(`${selectedDate}T00:00:00Z`).toISOString(),
        start_time: combineDateTime(selectedDate, extractHour(selectedSlot.start_time).toString().padStart(2, "0")),
        end_time: combineDateTime(selectedDate, extractHour(selectedSlot.end_time).toString().padStart(2, "0")),
      });
  
      onClose();
    }
  };  

  return (
    <Dialog open={isOpen} onClose={onClose} className="fixed inset-0 z-50">
      <div className="bg-black opacity-40 fixed inset-0" />
      <div className="fixed inset-0 flex items-center justify-center">
        <DialogPanel className="bg-white p-6 rounded w-full max-w-md">
          <DialogTitle className="text-xl font-bold mb-4">Select a Time Slot</DialogTitle>

          <label className="block text-sm mb-2 font-medium">Choose a Date:</label>
          <input
            type="date"
            className="mb-4 w-full px-4 py-2 border rounded"
            value={selectedDate}
            min={new Date().toISOString().split("T")[0]}
            onChange={(e) => setSelectedDate(e.target.value)}
          />

          <div className="space-y-2 max-h-[300px] overflow-y-auto">
            {timeSlots.map((slot) => (
              <button
                key={slot.id}
                className={`w-full text-left px-4 py-2 rounded ${
                  selectedSlot?.id === slot.id ? "bg-blue-500 text-white" : "bg-gray-100"
                }`}
                onClick={() => setSelectedSlot(slot)}
              >
                {formatDay(slot.day_of_week)} | {formatTime(extractHour(slot.start_time))} - {formatTime(extractHour(slot.end_time))}
              </button>
            ))}
          </div>

          <div className="mt-4 flex justify-end space-x-2">
            <button onClick={onClose} className="px-4 py-2 bg-gray-300 rounded">Cancel</button>
            <button
              onClick={handleConfirm}
              disabled={!selectedSlot}
              className="px-4 py-2 bg-blue-600 text-white rounded disabled:opacity-50"
            >
              Confirm
            </button>
          </div>
        </DialogPanel>
      </div>
    </Dialog>
  );
};

export default TimeSlotModal;
