import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Routes, Route } from "react-router-dom";

// Mock API URLs before component import
vi.mock("@/lib/api", () => ({
  googleLoginUrl: () => "http://localhost:8080/auth/google/login",
  githubLoginUrl: () => "http://localhost:8080/auth/github/login",
}));

// Mock AppContext
vi.mock("@/store/AppContext", () => ({
  useApp: vi.fn(),
}));

// Mock the hero image so vite asset import doesn't break jsdom
vi.mock("@/assets/hero-ocean.jpg", () => ({ default: "hero.jpg" }));

import { useApp } from "@/store/AppContext";
import Login from "@/pages/Login";

function renderLogin() {
  return render(
    <MemoryRouter initialEntries={["/"]}>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/dashboard" element={<div>dashboard</div>} />
      </Routes>
    </MemoryRouter>
  );
}

describe("Login", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useApp).mockReturnValue({
      user: null,
      loading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });
  });

  afterEach(() => {
    // Restore window.location (and any other stubbed globals)
    vi.unstubAllGlobals();
  });

  it("renders Google and GitHub sign-in buttons", () => {
    renderLogin();
    expect(screen.getByRole("button", { name: /continue with google/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /continue with github/i })).toBeInTheDocument();
  });

  it("clicking Google button navigates to Google OAuth URL", async () => {
    const user = userEvent.setup();
    // jsdom ignores cross-origin href assignments; stub window.location so we
    // can observe what the component sets on window.location.href.
    const locationStub = { href: "" };
    vi.stubGlobal("location", locationStub);

    renderLogin();
    await user.click(screen.getByRole("button", { name: /continue with google/i }));

    expect(locationStub.href).toBe("http://localhost:8080/auth/google/login");
  });

  it("clicking GitHub button navigates to GitHub OAuth URL", async () => {
    const user = userEvent.setup();
    const locationStub = { href: "" };
    vi.stubGlobal("location", locationStub);

    renderLogin();
    await user.click(screen.getByRole("button", { name: /continue with github/i }));

    expect(locationStub.href).toBe("http://localhost:8080/auth/github/login");
  });

  it("redirects authenticated users to /dashboard", () => {
    vi.mocked(useApp).mockReturnValue({
      user: { id: "u1", email: "x@y.com", name: "X", avatar_url: "", provider: "google" },
      loading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderLogin();

    expect(screen.getByText("dashboard")).toBeInTheDocument();
  });

  it("shows a loading spinner while auth state resolves", () => {
    vi.mocked(useApp).mockReturnValue({
      user: null,
      loading: true,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderLogin();

    const spinner = document.querySelector(".animate-spin");
    expect(spinner).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /continue with google/i })).not.toBeInTheDocument();
  });
});
