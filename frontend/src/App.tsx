import { Navigate, Route, Routes } from "react-router-dom";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import ProtectedRoute from "./routes/ProtectedRoute";
import PractitionerListPage from "./pages/practitioners/PractitionerListPage";
import AvailabilityPage from "./pages/practitioners/AvailabilityPage";
import AppointmentsPage from "./pages/appointments/AppointmentsPage";
import ReschedulePage from "./pages/appointments/ReschedulePage";

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route element={<ProtectedRoute />}>
        <Route
          path="/practitioners"
          element={<PractitionerListPage />}
        />

        <Route
          path="/practitioners/:id/availability"
          element={<AvailabilityPage />}
        />
        <Route
           path="/appointments"
           element={<AppointmentsPage />}
        />
        <Route path="/appointments/:id/reschedule"
        element={<ReschedulePage/>}
        />

      </Route>

      <Route
        path="*"
        element={<Navigate to="/practitioners" replace />}
      />
    </Routes>
  );
}

export default App;