import { describe, it, expect, vi, beforeEach, type MockInstance } from "vitest";
import { render, screen, waitFor, act } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter, Routes, Route } from "react-router-dom";
import type { DragEndEvent } from "@dnd-kit/core";

// ── dnd-kit mocks ────────────────────────────────────────────────────────────
// Capture onDragEnd so tests can trigger it directly.
let capturedOnDragEnd: ((e: DragEndEvent) => void) | null = null;

vi.mock("@dnd-kit/core", () => ({
  DndContext: ({ children, onDragEnd }: { children: React.ReactNode; onDragEnd: (e: DragEndEvent) => void }) => {
    capturedOnDragEnd = onDragEnd;
    return <div>{children}</div>;
  },
  PointerSensor: class {},
  closestCenter: vi.fn(),
  useSensor: vi.fn(),
  useSensors: vi.fn(() => []),
}));

vi.mock("@dnd-kit/sortable", () => ({
  SortableContext: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  arrayMove: <T,>(arr: T[], from: number, to: number): T[] => {
    const r = [...arr];
    const [item] = r.splice(from, 1);
    r.splice(to, 0, item);
    return r;
  },
  useSortable: () => ({
    attributes: {},
    listeners: {},
    setNodeRef: vi.fn(),
    transform: null,
    transition: undefined,
    isDragging: false,
  }),
  verticalListSortingStrategy: vi.fn(),
}));

vi.mock("@dnd-kit/utilities", () => ({
  CSS: { Transform: { toString: () => "" } },
}));

// ── API mocks ────────────────────────────────────────────────────────────────
const mockReorderItems = vi.fn();
const mockInviteCollaborator = vi.fn();
const mockFetchTrip = vi.fn();
const mockFetchTripDays = vi.fn();
const mockFetchDayItems = vi.fn();
const mockFetchCollaborators = vi.fn();
const mockDeleteItem = vi.fn();
const mockRemoveCollaborator = vi.fn();
const mockSearchPlaces = vi.fn();

vi.mock("@/lib/api", () => ({
  fetchTrip: (...args: unknown[]) => mockFetchTrip(...args),
  fetchTripDays: (...args: unknown[]) => mockFetchTripDays(...args),
  fetchDayItems: (...args: unknown[]) => mockFetchDayItems(...args),
  reorderItems: (...args: unknown[]) => mockReorderItems(...args),
  createDayItem: vi.fn(),
  deleteItem: (...args: unknown[]) => mockDeleteItem(...args),
  fetchCollaborators: (...args: unknown[]) => mockFetchCollaborators(...args),
  inviteCollaborator: (...args: unknown[]) => mockInviteCollaborator(...args),
  removeCollaborator: (...args: unknown[]) => mockRemoveCollaborator(...args),
  searchPlaces: (...args: unknown[]) => mockSearchPlaces(...args),
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

// ── Fixtures ─────────────────────────────────────────────────────────────────
const TRIP = {
  id: "trip-1",
  title: "Paris Adventure",
  destination: "Paris, France",
  cover_image_url: "",
  start_date: "2026-07-01",
  end_date: "2026-07-07",
  status: "active",
  visibility: "private",
  budget_total: 2000,
  currency: "USD",
};

const DAY = { id: "day-1", trip_id: "trip-1", day_number: 1, date: "2026-07-01", notes: "" };

const ITEMS = [
  { id: "item-1", day_id: "day-1", trip_id: "trip-1", title: "Eiffel Tower", type: "activity", order_index: 0, description: "", location: "", place_id: "", start_time: "", end_time: "", cost: 0, currency: "USD" },
  { id: "item-2", day_id: "day-1", trip_id: "trip-1", title: "Louvre", type: "activity", order_index: 1, description: "", location: "", place_id: "", start_time: "", end_time: "", cost: 0, currency: "USD" },
];

// ── Helpers ───────────────────────────────────────────────────────────────────
function makeQC() {
  return new QueryClient({ defaultOptions: { queries: { retry: false } } });
}

function renderTripDetail(tripId = "trip-1") {
  const qc = makeQC();
  return render(
    <QueryClientProvider client={qc}>
      <MemoryRouter initialEntries={[`/trips/${tripId}`]}>
        <Routes>
          <Route path="/trips/:id" element={<TripDetail />} />
          <Route path="/trips" element={<div>trips list</div>} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  );
}

import TripDetail from "@/pages/TripDetail";

// ── DaySection / reorder tests ────────────────────────────────────────────────
describe("TripDetail — dnd-kit item reorder", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    capturedOnDragEnd = null;
    mockFetchTrip.mockResolvedValue(TRIP);
    mockFetchTripDays.mockResolvedValue([DAY]);
    mockFetchDayItems.mockResolvedValue(ITEMS);
    mockFetchCollaborators.mockResolvedValue([]);
    mockReorderItems.mockResolvedValue({});
  });

  it("renders itinerary items", async () => {
    renderTripDetail();

    await waitFor(() => {
      expect(screen.getByText("Eiffel Tower")).toBeInTheDocument();
      expect(screen.getByText("Louvre")).toBeInTheDocument();
    });
  });

  it("calls reorderItems with reversed IDs when drag swaps two items", async () => {
    renderTripDetail();

    // Wait for items to render so DndContext mounts and captures onDragEnd
    await waitFor(() => expect(screen.getByText("Eiffel Tower")).toBeInTheDocument());

    expect(capturedOnDragEnd).not.toBeNull();

    // Simulate dragging item-1 onto item-2 (swap)
    await act(async () => {
      capturedOnDragEnd!({
        active: { id: "item-1", data: { current: undefined }, rect: { current: { initial: null, translated: null } } },
        over: { id: "item-2", data: { current: undefined }, rect: { width: 0, height: 0, left: 0, top: 0, right: 0, bottom: 0 } },
        collisions: [],
        activatorEvent: new Event("pointermove"),
        delta: { x: 0, y: 0 },
      } as unknown as DragEndEvent);
    });

    await waitFor(() => {
      expect(mockReorderItems).toHaveBeenCalledWith("trip-1", ["item-2", "item-1"]);
    });
  });

  it("does not call reorderItems when item is dropped on itself", async () => {
    renderTripDetail();

    await waitFor(() => expect(screen.getByText("Eiffel Tower")).toBeInTheDocument());

    await act(async () => {
      capturedOnDragEnd!({
        active: { id: "item-1", data: { current: undefined }, rect: { current: { initial: null, translated: null } } },
        over: { id: "item-1", data: { current: undefined }, rect: { width: 0, height: 0, left: 0, top: 0, right: 0, bottom: 0 } },
        collisions: [],
        activatorEvent: new Event("pointermove"),
        delta: { x: 0, y: 0 },
      } as unknown as DragEndEvent);
    });

    await waitFor(() => {
      expect(mockReorderItems).not.toHaveBeenCalled();
    });
  });
});

// ── CollaboratorsPanel tests ──────────────────────────────────────────────────
describe("TripDetail — CollaboratorsPanel invite form", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockFetchTrip.mockResolvedValue(TRIP);
    mockFetchTripDays.mockResolvedValue([DAY]);
    mockFetchDayItems.mockResolvedValue(ITEMS);
    mockFetchCollaborators.mockResolvedValue([]);
    mockInviteCollaborator.mockResolvedValue({});
  });

  async function openCollabPanel() {
    const user = userEvent.setup();
    renderTripDetail();

    await waitFor(() => expect(screen.getByRole("button", { name: /share/i })).toBeInTheDocument());
    await user.click(screen.getByRole("button", { name: /share/i }));

    await waitFor(() => expect(screen.getByPlaceholderText(/invite by email/i)).toBeInTheDocument());
    return user;
  }

  it("renders the invite form when Share is clicked", async () => {
    await openCollabPanel();
    expect(screen.getByPlaceholderText(/invite by email/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /invite/i })).toBeInTheDocument();
  });

  it("does not call inviteCollaborator when email is empty", async () => {
    const user = await openCollabPanel();
    await user.click(screen.getByRole("button", { name: /invite/i }));
    expect(mockInviteCollaborator).not.toHaveBeenCalled();
  });

  it("calls inviteCollaborator with trimmed email and selected role", async () => {
    const user = await openCollabPanel();

    await user.type(screen.getByPlaceholderText(/invite by email/i), "alice@example.com");
    await user.click(screen.getByRole("button", { name: /invite/i }));

    await waitFor(() => {
      expect(mockInviteCollaborator).toHaveBeenCalledWith("trip-1", "alice@example.com", "viewer");
    });
  });

  it("shows 'No collaborators yet' message when list is empty", async () => {
    await openCollabPanel();
    await waitFor(() => {
      expect(screen.getByText(/no collaborators yet/i)).toBeInTheDocument();
    });
  });

  it("renders existing collaborators with name and role", async () => {
    mockFetchCollaborators.mockResolvedValue([
      { user_id: "u2", name: "Bob", email: "bob@example.com", role: "editor", avatar_url: "" },
    ]);
    const user = await openCollabPanel();
    await waitFor(() => expect(screen.getByText("Bob")).toBeInTheDocument());
    expect(screen.getByText("editor")).toBeInTheDocument();
    void user; // suppress unused var warning
  });
});
