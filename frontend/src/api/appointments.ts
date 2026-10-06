import { apiFetch } from "./client";
import type {
  Appointment,
  CreateAppointmentRequest,
  RescheduleAppointmentRequest,
} from "../types/api";

export function createAppointment(
  data: CreateAppointmentRequest,
): Promise<Appointment> {
  return apiFetch<Appointment>("/appointments", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function getAppointments(): Promise<Appointment[]> {
  return apiFetch<Appointment[]>("/appointments");
}

export function cancelAppointment(
  id: string,
): Promise<void> {
  return apiFetch<void>(`/appointments/${id}`, {
    method: "DELETE",
  });
}

export function rescheduleAppointment(
  id: string,
  data: RescheduleAppointmentRequest,
): Promise<Appointment> {
  return apiFetch<Appointment>(`/appointments/${id}`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });
}