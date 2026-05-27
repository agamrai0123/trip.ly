import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { type ApiUser } from "@/lib/api";
import { Compass } from "lucide-react";

/**
 * Handles the redirect from the backend OAuth callback.
 *
 * The backend redirects here as:
 *   /auth/callback?access_token=<jwt>&user=<base64-json>
 *
 * We read the token + user, store them in AuthContext, then send the user
 * to their dashboard.
 */
const AuthCallback = () => {
  const { login } = useApp();
  const nav = useNavigate();
  const handled = useRef(false);

  useEffect(() => {
    if (handled.current) return;
    handled.current = true;

    const params = new URLSearchParams(window.location.search);
    const token = params.get("access_token");
    const userB64 = params.get("user");

    if (!token || !userB64) {
      // Malformed callback — send back to login
      nav("/", { replace: true });
      return;
    }

    try {
      const user: ApiUser = JSON.parse(atob(userB64));
      login(token, user);
      nav("/dashboard", { replace: true });
    } catch {
      nav("/", { replace: true });
    }
  }, [login, nav]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-4">
      <div className="grid h-14 w-14 place-items-center rounded-2xl bg-cta shadow-glow">
        <Compass className="h-7 w-7 text-primary-foreground" />
      </div>
      <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      <p className="text-sm text-muted-foreground">Signing you in…</p>
    </div>
  );
};

export default AuthCallback;
