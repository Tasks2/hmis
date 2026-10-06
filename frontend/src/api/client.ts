const API_URL = import.meta.env.VITE_API_URL;

interface ApiErrorBody {
  error?: {
    code?: string;
    message?: string;
  };
}

export class ApiError extends Error {
  status: number;
  code?: string;

  constructor(
    message: string,
    status: number,
    code?: string,
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token = localStorage.getItem("hmis_token");


  const headers = new Headers(options.headers);

  headers.set("Content-Type", "application/json");

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
  if (response.status === 401) {
    localStorage.removeItem("hmis_token");
  }

  let errorBody: ApiErrorBody = {};

  try {
    errorBody = await response.json();
  } catch {
    // Ignore invalid/non-JSON error responses.
  }

  throw new ApiError(
    errorBody.error?.message ??
      "An unexpected error occurred",
    response.status,
    errorBody.error?.code,
  );
}

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}