export interface AuthResponse {
  token: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface Practitioner {
  id: string;
  first_name: string;
  last_name: string;
  specialty: string;
}

export interface AvailableSlot {
  start_time: string;
  end_time: string;
}

export interface Appointment {
  id: string;
  practitioner_id: string;
  appointment_date: string;
  start_time: string;
  end_time: string;
  status: "BOOKED" | "CANCELLED";
}

export interface CreateAppointmentRequest {
  practitioner_id: string;
  appointment_date: string;
  start_time: string;
}

export interface RescheduleAppointmentRequest {
  appointment_date: string;
  start_time: string;
}