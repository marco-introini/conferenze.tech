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
  attendees?: {
    user: {
      id: string;
      nickname?: string;
      city?: string;
    };
    needsRide: boolean;
    hasCar: boolean;
  }[];
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
