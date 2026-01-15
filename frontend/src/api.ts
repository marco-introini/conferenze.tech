const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export interface User {
  id: string;
  email: string;
  name: string;
  nickname?: string;
  city?: string;
  avatarUrl?: string;
  bio?: string;
  createdAt?: string;
}

export interface Conference {
  id: string;
  title: string;
  date: string;
  location: string;
  website?: string;
  latitude?: number;
  longitude?: number;
  attendees?: Attendee[];
}

export interface Attendee {
  user: {
    id: string;
    nickname?: string;
    city?: string;
  };
  needsRide: boolean;
  hasCar: boolean;
}

export interface Registration {
  id: string;
  conferenceId: string;
  conferenceTitle: string;
  conferenceDate: string;
  conferenceLocation: string;
  status: string;
  role: string;
  needsRide: boolean;
  hasCar: boolean;
  registeredAt: string;
}

let authToken: string | null = null;

export function setAuthToken(token: string | null) {
  authToken = token;
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  if (authToken) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${authToken}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  console.log(`[API] ${options.method || 'GET'} ${endpoint}`, { status: response.status, ok: response.ok });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: "Errore" }));
    console.error(`[API] Error response:`, error);
    throw new Error(error.error || "Errore");
  }

  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}

export const api = {
  login: (email: string, password: string) =>
    request<{ user: User; token: string }>("/api/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  register: (data: {
    email: string;
    password: string;
    name: string;
    nickname?: string;
    city?: string;
    avatarUrl?: string;
    bio?: string;
  }) =>
    request<{ user: User; token: string }>("/api/register", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getConferences: () =>
    request<Conference[]>("/api/conferences", { method: "GET" }),

  getConference: (id: string) =>
    request<Conference>(`/api/conferences/${id}`, { method: "GET" }),

  createConference: (data: {
    title: string;
    date: string;
    location: string;
    website?: string;
    latitude?: number;
    longitude?: number;
  }) =>
    request<Conference>("/api/conferences", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  registerToConference: (
    conferenceId: string,
    data: {
      role: string;
      notes?: string;
      needsRide: boolean;
      hasCar: boolean;
    }
  ) =>
    request<Registration>(`/api/conferences/${conferenceId}/register`, {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getMyRegistrations: () =>
    request<Registration[]>("/api/users/registrations", { method: "GET" }),

  getMe: () =>
    request<User>("/api/me", { method: "GET" }),
};
