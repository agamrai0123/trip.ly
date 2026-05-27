export type Slot = "morning" | "lunch" | "evening" | "night";

export interface Place {
  id: string;
  name: string;
  category: "stay" | "food" | "sight" | "activity";
  slot: Slot;
  lat: number;
  lng: number;
  address: string;
  expense: number; // USD
  review: string;
  rating: number; // 0-5
  bestTime?: string;
  image: string;
}

export interface TripPost {
  id: string;
  cityId: string;
  author: { name: string; avatar: string };
  title: string;
  days: number;
  totalBudget: number;
  bestTime: string;
  likes: number;
  liked?: boolean;
  cover: string;
  summary: string;
  places: Place[]; // ordered as a single-day flow; can be reused per day in itinerary
  comments: Comment[];
}

export interface Comment {
  id: string;
  author: string;
  avatar: string;
  text: string;
  at: string;
}

export interface City {
  id: string;
  name: string;
  country: string;
  image: string;
  blurb: string;
  postsCount: number;
}

export interface ItineraryItem {
  id: string;
  dayIndex: number; // 0..days-1
  slot: Slot;
  place: Place;
  visited: boolean;
}

export interface MyTrip {
  id: string;
  cityId: string;
  cityName: string;
  days: number;
  items: ItineraryItem[];
  createdFromPostId?: string;
}
