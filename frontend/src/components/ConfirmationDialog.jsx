import React from "react";

const ConfirmationDialog = ({ isOpen, onConfirm, onCancel }) => {
  if (!isOpen) return null;

  return (
    <>
      <div className="fixed inset-0 bg-black opacity-40 z-40" />

      <div className="fixed inset-0 flex items-center justify-center z-50">
        <div className="bg-white p-6 rounded shadow-lg w-full max-w-sm">
          <p className="mb-4">Are you sure you want to confirm this booking?</p>
          <div className="flex justify-end gap-4">
            <button
              onClick={onCancel}
              className="px-4 py-2 bg-gray-300 rounded"
            >
              No
            </button>
            <button
              onClick={onConfirm}
              className="px-4 py-2 bg-green-500 text-white rounded"
            >
              Yes
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default ConfirmationDialog;
