import { createContext, useCallback, useContext, useEffect, useState, ReactNode } from "react";
import { type ApiUser, authRefresh, authLogout, setAccessToken } from "@/lib/api";

interface AppState {
  user: ApiUser | null;
  /** true while attempting a silent refresh on first load */
  loading: boolean;
  /** Called by /auth/callback after a successful OAuth redirect */
  login: (token: string, user: ApiUser) => void;
  logout: () => Promise<void>;
}

const AppCtx = createContext<AppState | null>(null);

export function AppProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<ApiUser | null>(null);
  const [loading, setLoading] = useState(true);

  // On mount: silent refresh to restore session from httpOnly cookie
  useEffect(() => {
    authRefresh()
      .then(({ access_token, user: u }) => {
        setAccessToken(access_token);
        setUser(u);
      })
      .catch(() => { /* no session — stay logged out */ })
      .finally(() => setLoading(false));
  }, []);

  const login = useCallback((token: string, u: ApiUser) => {
    setAccessToken(token);
    setUser(u);
  }, []);

  const logout = useCallback(async () => {
    try { await authLogout(); } catch { /* best-effort */ }
    setUser(null);
    setAccessToken(null);
  }, []);

  return (
    <AppCtx.Provider value={{ user, loading, login, logout }}>
      {children}
    </AppCtx.Provider>
  );
}

export function useApp() {
  const ctx = useContext(AppCtx);
  if (!ctx) throw new Error("useApp must be used within AppProvider");
  return ctx;
}
