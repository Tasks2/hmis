import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getPractitioners } from "../../api/practitioners";
import { useAuth } from "../../context/AuthContext";
import Navbar from "../../components/Navbar";

export default function PractitionerListPage() {
  const { logout } = useAuth();

  const {
    data: practitioners,
    isLoading,
    isError,
    refetch,
  } = useQuery({
    queryKey: ["practitioners"],
    queryFn: getPractitioners,
  });

  if (isLoading) {
    return (
      <main className="mx-auto max-w-5xl px-4 py-8">
        <p className="text-gray-600">
          Loading practitioners...
        </p>
      </main>
    );
  }

  if (isError) {
    return (
      <main className="mx-auto max-w-5xl px-4 py-8">
        <div className="rounded-lg bg-red-50 p-6">
          <p className="text-red-700">
            Unable to load practitioners.
          </p>

          <button
            onClick={() => refetch()}
            className="mt-3 rounded bg-red-600 px-4 py-2 text-white"
          >
            Try again
          </button>
        </div>
      </main>
    );
  }

  return (
    
    <main className="mx-auto max-w-5xl px-4 py-8">
      <Navbar/>
      <header className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            Practitioners
          </h1>

          <p className="mt-1 text-gray-600">
            Choose a practitioner to view available
            appointments.
          </p>
        </div>

        <button
          onClick={logout}
          className="rounded border px-4 py-2 text-gray-700"
        >
          Log out
        </button>
      </header>

      {!practitioners || practitioners.length === 0 ? (
        <div className="rounded-lg border bg-white p-8 text-center">
          <p className="text-gray-600">
            No practitioners are currently available.
          </p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {practitioners.map((practitioner) => (
            <article
              key={practitioner.id}
              className="rounded-lg border bg-white p-6 shadow-sm"
            >
              <h2 className="text-xl font-semibold text-gray-900">
                Dr. {practitioner.first_name}{" "}
                {practitioner.last_name}
              </h2>

              <p className="mt-1 text-gray-600">
                {practitioner.specialty}
              </p>

              <Link
                to={`/practitioners/${practitioner.id}/availability`}
                className="mt-5 inline-block rounded bg-blue-600 px-4 py-2 font-medium text-white"
              >
                View availability
              </Link>
            </article>
          ))}
        </div>
      )}
    </main>
  );
}