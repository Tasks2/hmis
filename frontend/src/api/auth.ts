import { apiFetch } from "./client";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
} from "../types/api";

export function register(
  data: RegisterRequest,
): Promise<void> {
  return apiFetch<void>("/auth/register", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function login(
  data: LoginRequest,
): Promise<AuthResponse> {
  return apiFetch<AuthResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(data),
  });
}