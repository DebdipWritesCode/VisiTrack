import React, { useRef } from "react";
import QRCode from "react-qr-code";

const QRCodeModal = ({ isOpen, onClose, value }) => {
  const qrCodeRef = useRef(null);

  if (!isOpen) return null;

  const downloadQRCode = () => {

    const svg = qrCodeRef.current;
    const svgData = new XMLSerializer().serializeToString(svg);
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    const img = new Image();
    
    canvas.width = svg.width.baseVal.value;
    canvas.height = svg.height.baseVal.value;
    
    img.onload = () => {
      ctx.fillStyle = "white";
      ctx.fillRect(0, 0, canvas.width, canvas.height);
      
      ctx.drawImage(img, 0, 0);
      
      const downloadLink = document.createElement("a");
      downloadLink.download = "qr-code.png";
      downloadLink.href = canvas.toDataURL("image/png");
      downloadLink.click();
    };
    
    img.src = "data:image/svg+xml;base64," + btoa(svgData);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black opacity-40" />

      <div className="bg-white p-6 rounded shadow-lg max-w-sm w-full flex flex-col items-center relative z-10">
        <p className="mb-4 text-center font-semibold">Show this QR code at the entrance:</p>
        <QRCode ref={qrCodeRef} value={value} />
        <div className="mt-4 text-center flex space-x-3">
          <button
            onClick={downloadQRCode}
            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          >
            Download QR
          </button>
          <button
            onClick={onClose}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default QRCodeModal;