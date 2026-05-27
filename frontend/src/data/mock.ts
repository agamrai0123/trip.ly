import type { City, TripPost } from "@/types";
import tokyo from "@/assets/city-tokyo.jpg";
import santorini from "@/assets/city-santorini.jpg";
import bali from "@/assets/city-bali.jpg";
import lisbon from "@/assets/city-lisbon.jpg";
import kyoto from "@/assets/city-kyoto.jpg";
import reykjavik from "@/assets/city-reykjavik.jpg";

export const cities: City[] = [
  { id: "tokyo", name: "Tokyo", country: "Japan", image: tokyo, blurb: "Neon nights, tranquil shrines.", postsCount: 184 },
  { id: "santorini", name: "Santorini", country: "Greece", image: santorini, blurb: "Whitewashed cliffs over the Aegean.", postsCount: 122 },
  { id: "bali", name: "Bali", country: "Indonesia", image: bali, blurb: "Rice terraces and ocean temples.", postsCount: 209 },
  { id: "lisbon", name: "Lisbon", country: "Portugal", image: lisbon, blurb: "Trams, tiles, and Atlantic light.", postsCount: 96 },
  { id: "kyoto", name: "Kyoto", country: "Japan", image: kyoto, blurb: "Bamboo groves and quiet tea houses.", postsCount: 158 },
  { id: "reykjavik", name: "Reykjavik", country: "Iceland", image: reykjavik, blurb: "Aurora, geysers, the wild north.", postsCount: 71 },
];

const u = (seed: string) => `https://api.dicebear.com/9.x/notionists/svg?seed=${seed}`;
const ph = (q: string, w = 800, h = 600) => `https://images.unsplash.com/photo-${q}?w=${w}&h=${h}&fit=crop&auto=format&q=80`;

export const tripPosts: TripPost[] = [];