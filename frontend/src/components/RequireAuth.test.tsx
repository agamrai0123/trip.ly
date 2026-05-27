import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router-dom";
import { RequireAuth } from "@/components/RequireAuth";

// Mock AppContext so tests control auth state
vi.mock("@/store/AppContext", () => ({
  useApp: vi.fn(),
}));

import { useApp } from "@/store/AppContext";

const Protected = () => <div>protected content</div>;

function renderWithRouter(initialEntry = "/") {
  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <Routes>
        <Route path="/" element={<div>login page</div>} />
        <Route
          path="/dashboard"
          element={
            <RequireAuth>
              <Protected />
            </RequireAuth>
          }
        />
      </Routes>
    </MemoryRouter>
  );
}

describe("RequireAuth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows a loading spinner while auth state is loading", () => {
    vi.mocked(useApp).mockReturnValue({
      user: null,
      loading: true,
      login: vi.fn(),
      logout: vi.fn(),
    });

    render(
      <MemoryRouter>
        <RequireAuth>
          <div>content</div>
        </RequireAuth>
      </MemoryRouter>
    );

    // The spinner is a div with animate-spin class — check it's not showing content
    expect(screen.queryByText("content")).not.toBeInTheDocument();
    const spinner = document.querySelector(".animate-spin");
    expect(spinner).toBeInTheDocument();
  });

  it("redirects to / when user is not authenticated", () => {
    vi.mocked(useApp).mockReturnValue({
      user: null,
      loading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderWithRouter("/dashboard");

    expect(screen.queryByText("protected content")).not.toBeInTheDocument();
    expect(screen.getByText("login page")).toBeInTheDocument();
  });

  it("renders children when user is authenticated", () => {
    vi.mocked(useApp).mockReturnValue({
      user: { id: "u1", email: "test@example.com", name: "Test", avatar_url: "", provider: "google" },
      loading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderWithRouter("/dashboard");

    expect(screen.getByText("protected content")).toBeInTheDocument();
  });
});
