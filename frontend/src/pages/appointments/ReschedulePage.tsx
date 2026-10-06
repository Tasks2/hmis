import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import { getAppointments, rescheduleAppointment } from "../../api/appointments";
import { getAvailability } from "../../api/practitioners";
import { ApiError } from "../../api/client";
import Navbar from "../../components/Navbar";

export default function ReschedulePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [date, setDate] = useState(
    new Date().toISOString().split("T")[0],
  );

  const [error, setError] = useState("");

  const { data: appointments } = useQuery({
    queryKey: ["appointments"],
    queryFn: getAppointments,
  });

  const appointment = appointments?.find(
    (item) => item.id === id,
  );

  const { data: slots, isLoading } = useQuery({
    queryKey: [
      "reschedule-availability",
      appointment?.practitioner_id,
      date,
    ],
    queryFn: () =>
      getAvailability(
        appointment!.practitioner_id,
        date,
      ),
    enabled: Boolean(
      appointment?.practitioner_id && date,
    ),
  });

  const mutation = useMutation({
    mutationFn: (startTime: string) =>
      rescheduleAppointment(id!, {
        appointment_date: date,
        start_time: startTime,
      }),
    onSuccess: () => {
      navigate("/appointments");
    },
  });

  async function handleReschedule(startTime: string) {
    setError("");

    try {
      await mutation.mutateAsync(startTime);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Unable to reschedule appointment.");
      }
    }
  }

  if (!appointment) {
    return (
      <main className="mx-auto max-w-4xl px-4 py-8">
        <p className="text-gray-600">
          Appointment not found.
        </p>
      </main>
    );
  }

  return (
    <main className="mx-auto max-w-4xl px-4 py-8">
     <Navbar />
      <Link
        to="/appointments"
        className="text-sm text-blue-600"
      >
        ← Back to appointments
      </Link>

      <h1 className="mt-6 text-3xl font-bold">
        Reschedule appointment
      </h1>

      <p className="mt-2 text-gray-600">
        Current appointment:{" "}
        {appointment.appointment_date} at{" "}
        {appointment.start_time}
      </p>

      {error && (
        <div className="mt-6 rounded-lg bg-red-50 p-4 text-red-700">
          {error}
        </div>
      )}

      <div className="mt-6">
        <label
          htmlFor="date"
          className="block text-sm font-medium"
        >
          New date
        </label>

        <input
          id="date"
          type="date"
          value={date}
          onChange={(event) => setDate(event.target.value)}
          className="mt-2 rounded border p-3"
        />
      </div>

      <section className="mt-8">
        {isLoading && (
          <p className="text-gray-600">
            Loading available slots...
          </p>
        )}

        {!isLoading &&
          (!slots || slots.length === 0) && (
            <p className="text-gray-600">
              No available slots for this date.
            </p>
          )}

        <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
          {slots?.map((slot) => (
            <button
              key={slot.start_time}
              disabled={mutation.isPending}
              onClick={() =>
                handleReschedule(slot.start_time)
              }
              className="rounded-lg border bg-white p-4 text-left hover:border-blue-500 hover:bg-blue-50 disabled:opacity-50"
            >
              {mutation.isPending
                ? "Updating..."
                : slot.start_time}
              <span className="block text-sm text-gray-500">
                until {slot.end_time}
              </span>
            </button>
          ))}
        </div>
      </section>
    </main>
  );
}