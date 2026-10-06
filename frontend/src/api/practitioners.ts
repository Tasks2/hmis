import { apiFetch } from "./client";
import type { AvailableSlot, Practitioner } from "../types/api";

interface AvailabilityResponse {
  date: string;
  practitioner_id: string;
  slots: AvailableSlot[] | null;
}

export function getPractitioners(): Promise<Practitioner[]> {
  return apiFetch<Practitioner[]>("/practitioners");
}

export async function getAvailability(
  practitionerId: string,
  date: string,
): Promise<AvailableSlot[]> {
  const response = await apiFetch<AvailabilityResponse>(
    `/practitioners/${practitionerId}/availability?date=${date}`,
  );

  return response.slots ?? [];
}