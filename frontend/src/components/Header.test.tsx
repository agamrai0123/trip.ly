import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter } from "react-router-dom";

// ── Mock useNotificationsWS so no real WebSocket is opened ───────────────────
vi.mock("@/hooks/useNotificationsWS", () => ({ useNotificationsWS: vi.fn() }));

// ── Mock next-themes ──────────────────────────────────────────────────────────
vi.mock("next-themes", () => ({
  useTheme: () => ({ theme: "light", setTheme: vi.fn() }),
}));

// ── API mocks ─────────────────────────────────────────────────────────────────
const mockMarkAllNotificationsRead = vi.fn().mockResolvedValue(undefined);
const mockMarkNotificationRead = vi.fn().mockResolvedValue(undefined);
const mockFetchTrips = vi.fn().mockResolvedValue([]);
const mockFetchNotifications = vi.fn();

vi.mock("@/lib/api", () => ({
  fetchTrips: (...args: unknown[]) => mockFetchTrips(...args),
  fetchNotifications: (...args: unknown[]) => mockFetchNotifications(...args),
  markAllNotificationsRead: (...args: unknown[]) => mockMarkAllNotificationsRead(...args),
  markNotificationRead: (...args: unknown[]) => mockMarkNotificationRead(...args),
}));

// ── Mock AppContext ────────────────────────────────────────────────────────────
vi.mock("@/store/AppContext", () => ({ useApp: vi.fn() }));

import { useApp } from "@/store/AppContext";
import { Header } from "@/components/Header";

const MOCK_USER = {
  id: "u1",
  email: "alice@example.com",
  name: "Alice Smith",
  avatar_url: "",
  provider: "google",
};

const MOCK_NOTIFICATIONS = [
  { id: "n1", type: "trip_invite", read: false, created_at: new Date().toISOString(), user_id: "u1", trip_id: "trip-1", message: "" },
  { id: "n2", type: "comment", read: true, created_at: new Date().toISOString(), user_id: "u1", trip_id: "trip-1", message: "" },
];

function renderHeader() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Header />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe("Header — notifications bell", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useApp).mockReturnValue({
      user: MOCK_USER,
      loading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });
    mockFetchTrips.mockResolvedValue([]);
  });

  it("shows unread count badge when there are unread notifications", async () => {
    mockFetchNotifications.mockResolvedValue(MOCK_NOTIFICATIONS);
    renderHeader();

    // Wait for query to settle
    await waitFor(() => {
      const badge = document.querySelector("span.rounded-full.bg-destructive");
      expect(badge).toBeInTheDocument();
      expect(badge!.textContent).toBe("1"); // 1 unread in MOCK_NOTIFICATIONS
    });
  });

  it("does not show badge when all notifications are read", async () => {
    mockFetchNotifications.mockResolvedValue(
      MOCK_NOTIFICATIONS.map((n) => ({ ...n, read: true }))
    );
    renderHeader();

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /notifications/i })).toBeInTheDocument();
    });

    const badge = document.querySelector("span.rounded-full.bg-destructive");
    expect(badge).not.toBeInTheDocument();
  });

  it("shows 'Mark all read' button when there are unread notifications", async () => {
    mockFetchNotifications.mockResolvedValue(MOCK_NOTIFICATIONS);
    renderHeader();

    // Open the dropdown
    const user = userEvent.setup();
    await waitFor(() =>
      expect(screen.getByRole("button", { name: /notifications/i })).toBeInTheDocument()
    );
    await user.click(screen.getByRole("button", { name: /notifications/i }));

    await waitFor(() => {
      expect(screen.getByText(/mark all read/i)).toBeInTheDocument();
    });
  });

  it("calls markAllNotificationsRead when 'Mark all read' is clicked", async () => {
    mockFetchNotifications.mockResolvedValue(MOCK_NOTIFICATIONS);
    const user = userEvent.setup();
    renderHeader();

    await waitFor(() =>
      expect(screen.getByRole("button", { name: /notifications/i })).toBeInTheDocument()
    );
    await user.click(screen.getByRole("button", { name: /notifications/i }));

    await waitFor(() => {
      expect(screen.getByText(/mark all read/i)).toBeInTheDocument();
    });
    await user.click(screen.getByText(/mark all read/i));

    await waitFor(() => {
      expect(mockMarkAllNotificationsRead).toHaveBeenCalledTimes(1);
    });
  });

  it("shows notification items in dropdown list", async () => {
    mockFetchNotifications.mockResolvedValue(MOCK_NOTIFICATIONS);
    const user = userEvent.setup();
    renderHeader();

    await waitFor(() =>
      expect(screen.getByRole("button", { name: /notifications/i })).toBeInTheDocument()
    );
    await user.click(screen.getByRole("button", { name: /notifications/i }));

    await waitFor(() => {
      expect(screen.getByText("trip invite")).toBeInTheDocument();
      expect(screen.getByText("comment")).toBeInTheDocument();
    });
  });

  it("shows 'No notifications yet.' when notification list is empty", async () => {
    mockFetchNotifications.mockResolvedValue([]);
    const user = userEvent.setup();
    renderHeader();

    await waitFor(() =>
      expect(screen.getByRole("button", { name: /notifications/i })).toBeInTheDocument()
    );
    await user.click(screen.getByRole("button", { name: /notifications/i }));

    await waitFor(() => {
      expect(screen.getByText(/no notifications yet/i)).toBeInTheDocument();
    });
  });
});
