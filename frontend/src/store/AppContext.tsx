import { createContext, useContext, useEffect, useState, ReactNode, useCallback } from "react";
import type { ItineraryItem, MyTrip, Place, TripPost } from "@/types";

interface User { email: string; name: string; }

interface AppState {
  user: User | null;
  login: (email: string) => void;
  logout: () => void;
  trips: MyTrip[];
  addPostToTrips: (post: TripPost) => string; // returns tripId
  toggleVisited: (tripId: string, itemId: string) => void;
  reviewPlace: (tripId: string, itemId: string, review: string, rating: number) => void;
  removeItem: (tripId: string, itemId: string) => void;
  updateItemExpense: (tripId: string, itemId: string, expense: number) => void;
  addPlaceToTrip: (tripId: string, place: Place, slot?: ItineraryItem["slot"], dayIndex?: number) => void;
  reorderTrip: (tripId: string, items: ItineraryItem[]) => void;
  createEmptyTrip: (cityId: string, days?: number) => string;
  likes: Record<string, { liked: boolean; count: number }>;
  toggleLike: (postId: string, baseCount: number) => void;
  addComment: (postId: string, text: string) => void;
  commentsByPost: Record<string, { id: string; author: string; avatar: string; text: string; at: string }[]>;
  // Posts created from user reviews
  userPosts: TripPost[];
}

const AppCtx = createContext<AppState | null>(null);

const LS = "wandr_state_v1";

function loadLS<T>(key: string, fb: T): T {
  try { const v = localStorage.getItem(key); return v ? JSON.parse(v) : fb; } catch { return fb; }
}

export function AppProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => loadLS<User | null>(LS + ":user", null));
  const [trips, setTrips] = useState<MyTrip[]>(() => loadLS<MyTrip[]>(LS + ":trips", []));
  const [likes, setLikes] = useState<AppState["likes"]>(() => loadLS(LS + ":likes", {}));
  const [commentsByPost, setComments] = useState<AppState["commentsByPost"]>(() => loadLS(LS + ":comments", {}));
  const [userPosts, setUserPosts] = useState<TripPost[]>(() => loadLS<TripPost[]>(LS + ":userPosts", []));

  useEffect(() => { localStorage.setItem(LS + ":user", JSON.stringify(user)); }, [user]);
  useEffect(() => { localStorage.setItem(LS + ":trips", JSON.stringify(trips)); }, [trips]);
  useEffect(() => { localStorage.setItem(LS + ":likes", JSON.stringify(likes)); }, [likes]);
  useEffect(() => { localStorage.setItem(LS + ":comments", JSON.stringify(commentsByPost)); }, [commentsByPost]);
  useEffect(() => { localStorage.setItem(LS + ":userPosts", JSON.stringify(userPosts)); }, [userPosts]);

  const login = useCallback((email: string) => {
    const name = email.split("@")[0].replace(/\W/g, " ").replace(/\b\w/g, c => c.toUpperCase()) || "Traveler";
    setUser({ email, name });
  }, []);
  const logout = useCallback(() => setUser(null), []);

  const addPostToTrips = useCallback((post: TripPost) => {
    // Create itinerary by spreading places across days, repeating slot pattern.
    // We schedule one place per slot per day in original order, looping if fewer places than slots.
    const slotsOrder: ItineraryItem["slot"][] = ["morning", "lunch", "evening", "night"];
    const items: ItineraryItem[] = [];
    let placeIdx = 0;
    for (let d = 0; d < post.days; d++) {
      for (const slot of slotsOrder) {
        // pick first place matching slot, fall back to round-robin
        const matched = post.places.find(p => p.slot === slot && !items.some(it => it.place.id === p.id && it.dayIndex === d));
        const chosen = matched ?? post.places[placeIdx % post.places.length];
        placeIdx++;
        items.push({
          id: `${post.id}-${d}-${slot}-${chosen.id}-${items.length}`,
          dayIndex: d,
          slot,
          place: { ...chosen },
          visited: false,
        });
      }
    }
    const trip: MyTrip = {
      id: `t-${Date.now()}`,
      cityId: post.cityId,
      cityName: post.cityId.charAt(0).toUpperCase() + post.cityId.slice(1),
      days: post.days,
      items,
      createdFromPostId: post.id,
    };
    setTrips(prev => [trip, ...prev]);
    return trip.id;
  }, []);

  const toggleVisited = useCallback((tripId: string, itemId: string) => {
    setTrips(prev => prev.map(t => t.id === tripId ? {
      ...t, items: t.items.map(it => it.id === itemId ? { ...it, visited: !it.visited } : it)
    } : t));
  }, []);

  const removeItem = useCallback((tripId: string, itemId: string) => {
    setTrips(prev => prev.map(t => t.id === tripId ? { ...t, items: t.items.filter(i => i.id !== itemId) } : t));
  }, []);

  const updateItemExpense = useCallback((tripId: string, itemId: string, expense: number) => {
    setTrips(prev => prev.map(t => t.id === tripId ? {
      ...t, items: t.items.map(it => it.id === itemId ? { ...it, place: { ...it.place, expense } } : it)
    } : t));
  }, []);

  const addPlaceToTrip = useCallback((tripId: string, place: Place, slot?: ItineraryItem["slot"], dayIndex?: number) => {
    setTrips(prev => prev.map(t => {
      if (t.id !== tripId) return t;
      const item: ItineraryItem = {
        id: `it-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
        dayIndex: dayIndex ?? 0,
        slot: slot ?? place.slot,
        place: { ...place },
        visited: false,
      };
      return { ...t, items: [...t.items, item] };
    }));
  }, []);

  const reorderTrip = useCallback((tripId: string, items: ItineraryItem[]) => {
    setTrips(prev => prev.map(t => t.id === tripId ? { ...t, items } : t));
  }, []);

  const createEmptyTrip = useCallback((cityId: string, days = 3) => {
    const trip: MyTrip = {
      id: `t-${Date.now()}`,
      cityId,
      cityName: cityId.charAt(0).toUpperCase() + cityId.slice(1),
      days,
      items: [],
    };
    setTrips(prev => [trip, ...prev]);
    return trip.id;
  }, []);

  const reviewPlace = useCallback((tripId: string, itemId: string, review: string, rating: number) => {
    let captured: { trip?: MyTrip; place?: Place } = {};
    setTrips(prev => prev.map(t => {
      if (t.id !== tripId) return t;
      const items = t.items.map(it => {
        if (it.id !== itemId) return it;
        const place = { ...it.place, review, rating };
        captured = { trip: t, place };
        return { ...it, visited: true, place };
      });
      return { ...t, items };
    }));
    // Auto-create a tiny user post from the review (reviews become posts)
    if (captured.trip && captured.place && user) {
      const np: TripPost = {
        id: `up-${Date.now()}`,
        cityId: captured.trip.cityId,
        author: { name: user.name, avatar: `https://api.dicebear.com/9.x/notionists/svg?seed=${user.name}` },
        title: `${captured.place.name} — quick take`,
        days: 1,
        totalBudget: captured.place.expense,
        bestTime: captured.place.bestTime ?? "Anytime",
        likes: 0,
        cover: captured.place.image,
        summary: review,
        places: [captured.place],
        comments: [],
      };
      setUserPosts(prev => [np, ...prev]);
    }
  }, [user]);

  const toggleLike = useCallback((postId: string, baseCount: number) => {
    setLikes(prev => {
      const cur = prev[postId] ?? { liked: false, count: baseCount };
      const liked = !cur.liked;
      return { ...prev, [postId]: { liked, count: cur.count + (liked ? 1 : -1) } };
    });
  }, []);

  const addComment = useCallback((postId: string, text: string) => {
    if (!user) return;
    setComments(prev => ({
      ...prev,
      [postId]: [
        ...(prev[postId] ?? []),
        { id: `uc-${Date.now()}`, author: user.name, avatar: `https://api.dicebear.com/9.x/notionists/svg?seed=${user.name}`, text, at: "now" },
      ],
    }));
  }, [user]);

  return (
    <AppCtx.Provider value={{
      user, login, logout, trips, addPostToTrips, toggleVisited, reviewPlace,
      removeItem, updateItemExpense, addPlaceToTrip, reorderTrip, createEmptyTrip,
      likes, toggleLike, addComment, commentsByPost, userPosts,
    }}>
      {children}
    </AppCtx.Provider>
  );
}

export function useApp() {
  const ctx = useContext(AppCtx);
  if (!ctx) throw new Error("useApp must be used within AppProvider");
  return ctx;
}
