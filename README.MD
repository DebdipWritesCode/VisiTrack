# VisiTrack – Intelligent Visitor Management System 🚪

A full-stack Visitor Management System designed for offices or gated venues to streamline the appointment process between visitors and hosts. This project features OTP-based login, smart availability booking, QR code-based check-in/check-out, and an intuitive admin panel for security personnel.

---

## 📚 Tech Stack

### 🔧 Backend
- **Language:** Go (Golang)
- **Framework:** Gin
- **ORM/DB Layer:** SQLC (type-safe SQL queries)
- **Database:** PostgreSQL
- **Authentication:** OTP-based login via mobile number

### 💻 Frontend
- **Library:** React
- **Language:** TypeScript (Optional: JS support)
- **State Management:** React Context / Zustand (Optional)
- **Styling:** TailwindCSS or CSS Modules
- **QR Code:** QRCode.js or similar libraries

---

## ⚙️ Features

### 👥 User Flow
- **Login via OTP** using mobile number.
- **Sign Up** with first name, last name, and mobile number.
- **Dashboard** shows:
  - Appointments booked by the user.
  - Appointments other people have booked with the user.
  - Graphical summaries of appointments.
- **Sidebar Navigation**:
  - Dashboard
  - Profile (with availability settings)
  - Appointments (incoming and outgoing)
  - Book Appointment
  - Logout
- **Profile Availability**:
  - Calendar + time slots from 9 AM to 6 PM
- **Book Appointment**:
  - Search for hosts or view popular hosts.
  - Choose from only available slots.
  - Confirm via modal, then get a **QR Code** + **SMS Link**
- **QR Code**:
  - Used for entry (check-in) and exit (check-out).
  - Scanned by security guards to manage appointment lifecycle.

### 🛡️ Admin Flow (Security Guards)
- View all appointments for any day.
- Scan QR codes:
  - First scan sets appointment as **"Ongoing"**.
  - Second scan sets appointment as **"Completed"**.
- See live status updates for visits.

---

## 🧩 Database Schema Overview

### Tables:
- `users`: user info including role (user/admin)
- `appointments`: all appointment data
- `availability`: time slots per user
- `otps`: OTP-based login verification
- `appointment_logs`: check-in/out tracking
- `appointment_stats`: tracks popular hosts

(Full SQL schema in `/db/schema.sql`)

---

## 🚀 Getting Started

### Backend (Go)
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend (React)
```bash
cd frontend
npm install
npm run dev
```

### PostgreSQL Setup
Make sure your `.env` file contains:
```
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=visitordb
DB_HOST=localhost
DB_PORT=5432
```

---

## 📦 Folder Structure

```
/backend
  ├── api/
  ├── db/ (SQLC queries, migrations)
  └── main.go

/frontend
  ├── src/
  ├── components/
  ├── pages/
  └── App.tsx

/db
  ├── schema.sql
  └── migrations/
```

---

## 🧠 Author
Made with ❤️ by Debdip Mukherjee
