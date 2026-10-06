import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { register } from "../../api/auth";
import { ApiError } from "../../api/client";

export default function RegisterPage() {
  const navigate = useNavigate();

  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    password: "",
  });

  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  function updateField(
    field: keyof typeof form,
    value: string,
  ) {
    setForm((current) => ({
      ...current,
      [field]: value,
    }));
  }

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault();

    setError("");
    setLoading(true);

    try {
      await register(form);
      navigate("/login");
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Unable to register.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-md rounded-lg bg-white p-8 shadow"
      >
        <h1 className="text-2xl font-bold text-gray-900">
          Create account
        </h1>

        <p className="mt-2 text-gray-600">
          Register as a patient.
        </p>

        {error && (
          <div className="mt-4 rounded-md bg-red-50 p-3 text-sm text-red-700">
            {error}
          </div>
        )}

        <div className="mt-6 space-y-4">
          <input
            className="w-full rounded border p-3"
            placeholder="First name"
            value={form.first_name}
            onChange={(e) =>
              updateField("first_name", e.target.value)
            }
            required
          />

          <input
            className="w-full rounded border p-3"
            placeholder="Last name"
            value={form.last_name}
            onChange={(e) =>
              updateField("last_name", e.target.value)
            }
            required
          />

          <input
            className="w-full rounded border p-3"
            type="email"
            placeholder="Email"
            value={form.email}
            onChange={(e) =>
              updateField("email", e.target.value)
            }
            required
          />

          <input
            className="w-full rounded border p-3"
            type="password"
            placeholder="Password"
            value={form.password}
            onChange={(e) =>
              updateField("password", e.target.value)
            }
            required
          />

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded bg-blue-600 p-3 font-medium text-white disabled:opacity-50"
          >
            {loading ? "Creating account..." : "Register"}
          </button>
        </div>

        <p className="mt-6 text-center text-sm text-gray-600">
          Already have an account?{" "}
          <Link
            to="/login"
            className="font-medium text-blue-600"
          >
            Log in
          </Link>
        </p>
      </form>
    </main>
  );
}