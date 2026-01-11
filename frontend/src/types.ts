export interface User {
  id: string;
  nickname: string;
  city: string;
}

export interface Conference {
  id: string;
  title: string;
  date: string;
  location: string;
  website: string;
  latitude?: number;
  longitude?: number;
  attendees: {
    user: User;
    needsRide: boolean;
    hasCar: boolean;
  }[];
}
