import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { getAppointments } from "../../api/appointments";
import Navbar from "../../components/Navbar";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { cancelAppointment } from "../../api/appointments";

export default function AppointmentsPage() {
  const {
    data: appointments,
    isLoading,
    isError,
    refetch,
  } = useQuery({
    queryKey: ["appointments"],
    queryFn: getAppointments,
  });

//    const today = new Date()
//   .toISOString()
//   .split("T")[0];

    const queryClient = useQueryClient();

    const cancelMutation = useMutation({
    mutationFn: cancelAppointment,
    onSuccess: () => {
        queryClient.invalidateQueries({
        queryKey: ["appointments"],
        });
    },
    });

  return (
    <main className="mx-auto max-w-5xl px-4 py-8">
    <Navbar/>
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            My Appointments
          </h1>

          <p className="mt-1 text-gray-600">
            View and manage your appointments.
          </p>
        </div>

        <Link
          to="/practitioners"
          className="rounded bg-blue-600 px-4 py-2 text-white"
        >
          Book appointment
        </Link>
      </div>

      {isLoading && (
        <p className="text-gray-600">
          Loading appointments...
        </p>
      )}

      {isError && (
        <div className="rounded-lg bg-red-50 p-6">
          <p className="text-red-700">
            Unable to load appointments.
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
        (!appointments || appointments.length === 0) && (
          <div className="rounded-lg border bg-white p-8 text-center">
            <p className="text-gray-600">
              You don't have any appointments yet.
            </p>
          </div>
        )}

      <div className="space-y-4">
        {appointments?.map((appointment) => (
          <article
            key={appointment.id}
            className="rounded-lg border bg-white p-6 shadow-sm"
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="font-semibold text-gray-900">
                  {appointment.appointment_date}
                </p>

                <p className="mt-1 text-gray-600">
                  {appointment.start_time} –{" "}
                  {appointment.end_time}
                </p>
              </div>

              <span className="rounded-full bg-green-100 px-3 py-1 text-sm text-green-700">
                {appointment.status}
              </span>
            </div>

            {appointment.status ===  "BOOKED" && (
            <div className="mt-4">
              <Link
                to={`/appointments/${appointment.id}/reschedule`}
                className="text-sm font-medium text-blue-600"
              >
                Reschedule
              </Link>

              <button
                onClick={() => {
                    if (
                    window.confirm(
                        "Cancel this appointment?",
                    )
                    ) {
                    cancelMutation.mutate(appointment.id);
                    }
                }}
                disabled={cancelMutation.isPending}
                className="ml-4 text-sm font-medium text-red-600 disabled:opacity-50"
                >
                {cancelMutation.isPending
                    ? "Cancelling..."
                    : "Cancel"}
                </button>
            </div>
            )}
          </article>
        ))}
      </div>
    </main>
  );
}