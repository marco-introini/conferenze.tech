const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export interface User {
  id: string;
  email: string;
  name: string;
  nickname?: string;
  city?: string;
  avatarUrl?: string;
  bio?: string;
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

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: "Credenziali errate" }));
    throw new Error(error.error || "Credenziali errate");
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
    request<Conference>("/api/conference?id=" + id, { method: "GET" }),

  createConference: (data: {
    title: string;
    date: string;
    location: string;
    website?: string;
    latitude?: number;
    longitude?: number;
  }) =>
    request<Conference>("/api/conferences/create", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  registerToConference: (
    userId: string,
    data: {
      conferenceId: string;
      role: string;
      notes?: string;
      needsRide: boolean;
      hasCar: boolean;
    }
  ) =>
    request<{ id: string }>("/api/register-to-conference?userId=" + userId, {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getMyRegistrations: (userId: string) =>
    request<Registration[]>("/api/my-registrations?userId=" + userId, {
      method: "GET",
    }),

  getMe: (userId: string) =>
    request<User>("/api/me?userId=" + userId, { method: "GET" }),
};
