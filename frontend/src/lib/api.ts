import axios from "axios";

const BASE = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

// In-memory access token store — never persisted to disk.
let _accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  _accessToken = token;
}
export function getAccessToken(): string | null {
  return _accessToken;
}

export const api = axios.create({
  baseURL: BASE,
  withCredentials: true, // send httpOnly refresh-token cookie
});

// Attach JWT on every request
api.interceptors.request.use((config) => {
  if (_accessToken) {
    config.headers.Authorization = `Bearer ${_accessToken}`;
  }
  return config;
});

let _refreshing: Promise<string> | null = null;

// On 401: try refresh once, then retry original request
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const original = err.config;
    if (err.response?.status === 401 && !original._retry) {
      original._retry = true;
      try {
        if (!_refreshing) {
          _refreshing = api
            .post<{ data: { access_token: string } }>("/auth/refresh")
            .then((r) => {
              const token = r.data.data.access_token;
              setAccessToken(token);
              return token;
            })
            .finally(() => {
              _refreshing = null;
            });
        }
        const token = await _refreshing;
        original.headers.Authorization = `Bearer ${token}`;
        return api(original);
      } catch {
        setAccessToken(null);
        window.location.href = "/";
        return Promise.reject(err);
      }
    }
    return Promise.reject(err);
  }
);

// ── Auth ────────────────────────────────────────────────────────────────────

export function googleLoginUrl() {
  return `${BASE}/auth/google/login`;
}
export function githubLoginUrl() {
  return `${BASE}/auth/github/login`;
}

export async function authRefresh(): Promise<{ access_token: string; user: ApiUser }> {
  const r = await api.post<{ data: { access_token: string; user: ApiUser } }>("/auth/refresh");
  return r.data.data;
}

export async function authMe(): Promise<ApiUser> {
  const r = await api.get<{ data: ApiUser }>("/auth/me");
  return r.data.data;
}

export async function authLogout(): Promise<void> {
  await api.post("/auth/logout");
  setAccessToken(null);
}

// ── Trips ───────────────────────────────────────────────────────────────────

export async function fetchTrips(): Promise<ApiTrip[]> {
  const r = await api.get<{ data: ApiTrip[] }>("/api/v1/trips");
  return r.data.data ?? [];
}

export async function fetchTrip(id: string): Promise<ApiTrip> {
  const r = await api.get<{ data: ApiTrip }>(`/api/v1/trips/${id}`);
  return r.data.data;
}

export async function createTrip(body: CreateTripBody): Promise<ApiTrip> {
  const r = await api.post<{ data: ApiTrip }>("/api/v1/trips", body);
  return r.data.data;
}

export async function updateTrip(id: string, body: Partial<CreateTripBody>): Promise<ApiTrip> {
  const r = await api.patch<{ data: ApiTrip }>(`/api/v1/trips/${id}`, body);
  return r.data.data;
}

export async function deleteTrip(id: string): Promise<void> {
  await api.delete(`/api/v1/trips/${id}`);
}

export async function fetchTripDays(tripId: string): Promise<ApiDay[]> {
  const r = await api.get<{ data: ApiDay[] }>(`/api/v1/trips/${tripId}/days`);
  return r.data.data ?? [];
}

export async function createTripDay(tripId: string, body: { day_number: number; notes?: string }): Promise<ApiDay> {
  const r = await api.post<{ data: ApiDay }>(`/api/v1/trips/${tripId}/days`, body);
  return r.data.data;
}

export async function fetchDayItems(tripId: string, dayId: string): Promise<ApiItem[]> {
  const r = await api.get<{ data: ApiItem[] }>(`/api/v1/trips/${tripId}/days/${dayId}/items`);
  return r.data.data ?? [];
}

export async function createDayItem(tripId: string, dayId: string, body: CreateItemBody): Promise<ApiItem> {
  const r = await api.post<{ data: ApiItem }>(`/api/v1/trips/${tripId}/days/${dayId}/items`, body);
  return r.data.data;
}

export async function updateItem(tripId: string, itemId: string, body: Partial<CreateItemBody>): Promise<ApiItem> {
  const r = await api.patch<{ data: ApiItem }>(`/api/v1/trips/${tripId}/items/${itemId}`, body);
  return r.data.data;
}

export async function deleteItem(tripId: string, itemId: string): Promise<void> {
  await api.delete(`/api/v1/trips/${tripId}/items/${itemId}`);
}

export async function reorderItems(tripId: string, itemIds: string[]): Promise<void> {
  await api.patch(`/api/v1/trips/${tripId}/items/reorder`, { item_ids: itemIds });
}

// ── Search ──────────────────────────────────────────────────────────────────

export async function searchTrips(q: string): Promise<ApiTrip[]> {
  const r = await api.get<{ data: ApiTrip[] }>("/api/v1/search/trips", { params: { q } });
  return r.data.data ?? [];
}

export async function searchPlaces(q: string): Promise<ApiPlace[]> {
  const r = await api.get<{ data: ApiPlace[] }>("/api/v1/search/places", { params: { q } });
  return r.data.data ?? [];
}

// ── Collaborators ────────────────────────────────────────────────────────────

export async function fetchCollaborators(tripId: string): Promise<ApiCollaborator[]> {
  const r = await api.get<{ data: ApiCollaborator[] }>(`/api/v1/trips/${tripId}/collaborators`);
  return r.data.data ?? [];
}

export async function inviteCollaborator(tripId: string, email: string, role: string): Promise<ApiCollaborator> {
  const r = await api.post<{ data: ApiCollaborator }>(`/api/v1/trips/${tripId}/collaborators`, { email, role });
  return r.data.data;
}

export async function removeCollaborator(tripId: string, userId: string): Promise<void> {
  await api.delete(`/api/v1/trips/${tripId}/collaborators/${userId}`);
}

// ── API Types ────────────────────────────────────────────────────────────────

export interface ApiUser {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  provider: string;
}

export interface ApiTrip {
  id: string;
  user_id: string;
  title: string;
  destination: string;
  cover_image_url: string;
  start_date: string | null;
  end_date: string | null;
  status: string;
  visibility: string;
  budget_total: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface ApiDay {
  id: string;
  trip_id: string;
  day_number: number;
  date: string | null;
  notes: string;
}

export interface ApiItem {
  id: string;
  day_id: string;
  trip_id: string;
  title: string;
  description: string;
  location: string;
  place_id: string;
  type: string;
  start_time: string;
  end_time: string;
  cost: number;
  currency: string;
  order_index: number;
  created_at: string;
  updated_at: string;
}

export interface ApiPlace {
  id: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
  type: string;
}

export interface ApiCollaborator {
  user_id: string;
  trip_id: string;
  role: string;
  email: string;
  name: string;
  avatar_url: string;
}

export interface CreateTripBody {
  title: string;
  destination: string;
  cover_image_url?: string;
  start_date?: string;
  end_date?: string;
  status?: string;
  visibility?: string;
  budget_total?: number;
  currency?: string;
}

export interface CreateItemBody {
  title: string;
  description?: string;
  location?: string;
  place_id?: string;
  type?: string;
  start_time?: string;
  end_time?: string;
  cost?: number;
  currency?: string;
  order_index?: number;
}
