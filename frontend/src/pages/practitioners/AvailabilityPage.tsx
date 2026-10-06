import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import { getAvailability } from "../../api/practitioners";
import { createAppointment } from "../../api/appointments";
import { ApiError } from "../../api/client";
import Navbar from "../../components/Navbar";

export default function AvailabilityPage() {
  const { id } = useParams<{ id: string }>();

  const [date, setDate] = useState(
    new Date().toISOString().split("T")[0],
  );

  const [message, setMessage] = useState("");
  const today = new Date().toISOString().split("T")[0];

  const {
    data: slots,
    isLoading,
    isError,
    refetch,
  } = useQuery({
    queryKey: ["availability", id, date],
    queryFn: () => getAvailability(id!, date),
    enabled: Boolean(id && date),
  });


  const bookingMutation = useMutation({
    mutationFn: (startTime: string) =>
      createAppointment({
        practitioner_id: id!,
        appointment_date: date,
        start_time: startTime,
      }),

    onSuccess: () => {
      setMessage("Appointment booked successfully.");

      refetch();
    },
  });

  async function handleBooking(startTime: string) {
    setMessage("");

    try {
      await bookingMutation.mutateAsync(startTime);
    } catch (err) {
      if (err instanceof ApiError) {
        setMessage(err.message);
      } else {
        setMessage("Unable to book appointment.");
      }
    }
  }

  return (
    <main className="mx-auto max-w-4xl px-4 py-8">
    <Navbar/>
      <Link
        to="/practitioners"
        className="text-sm text-blue-600"
      >
        ← Back to practitioners
      </Link>

      <div className="mt-6">
        <h1 className="text-3xl font-bold text-gray-900">
          Availability
        </h1>

        <div className="mt-6">
          <label
            htmlFor="date"
            className="block text-sm font-medium text-gray-700"
          >
            Appointment date
          </label>

          <input
            id="date"
            type="date"
            min={today}
            value={date}
            onChange={(event) => {
              setDate(event.target.value);
              setMessage("");
            }}
            className="mt-2 rounded border p-3"
          />
        </div>
      </div>

      {message && (
        <div className="mt-6 rounded-lg bg-blue-50 p-4 text-blue-700">
          {message}
        </div>
      )}

      <section className="mt-8">
        {isLoading && (
          <p className="text-gray-600">
            Loading available slots...
          </p>
        )}

        {isError && (
          <div className="rounded-lg bg-red-50 p-6">
            <p className="text-red-700">
              Unable to load availability.
            </p>

            <button
              onClick={() => refetch()}
              className="mt-3 rounded bg-red-600 px-4 py-2 text-white"
            >
              Try again
            </button>
          </div>
        )}

        {!isLoading &&
          !isError &&
          (!slots || slots.length === 0) && (
            <div className="rounded-lg border bg-white p-8 text-center">
              <p className="text-gray-600">
                No available appointments for this date.
              </p>
            </div>
          )}

        {!isLoading &&
          !isError &&
          slots &&
          slots.length > 0 && (
            <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
              {slots.map((slot) => (
                <button
                  key={`${slot.start_time}-${slot.end_time}`}
                  onClick={() =>
                    handleBooking(slot.start_time)
                  }
                  disabled={bookingMutation.isPending}
                  className="rounded-lg border bg-white p-4 text-left shadow-sm hover:border-blue-500 hover:bg-blue-50 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <p className="font-semibold text-gray-900">
                    {bookingMutation.isPending
                      ? "Booking..."
                      : slot.start_time}
                  </p>

                  <p className="text-sm text-gray-500">
                    until {slot.end_time}
                  </p>
                </button>
              ))}
            </div>
          )}
      </section>
    </main>
  );
}