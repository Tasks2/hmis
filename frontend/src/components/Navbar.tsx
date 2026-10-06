import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Navbar() {
  const { logout } = useAuth();

  return (
    <nav className="border-b bg-white">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
        <Link
          to="/practitioners"
          className="font-bold text-gray-900"
        >
          HMIS
        </Link>

        <div className="flex items-center gap-4">
          <Link
            to="/practitioners"
            className="text-sm text-gray-600 hover:text-gray-900"
          >
            Practitioners
          </Link>

          <Link
            to="/appointments"
            className="text-sm text-gray-600 hover:text-gray-900"
          >
            My appointments
          </Link>

          <button
            onClick={logout}
            className="text-sm text-red-600"
          >
            Log out
          </button>
        </div>
      </div>
    </nav>
  );
}