import { useEffect, useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { getAccessToken } from "@/lib/api";

// Derive WebSocket base URL from the HTTP API base URL.
const WS_BASE = (import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080").replace(
  /^http/,
  (s) => (s === "https" ? "wss" : "ws"),
);

/**
 * Opens a WebSocket connection to /ws/notifications and invalidates the
 * ['notifications'] query cache whenever the server pushes a message.
 * Reconnects automatically with a 5-second back-off on disconnect.
 *
 * @param enabled - only connect when the user is authenticated.
 */
export function useNotificationsWS(enabled: boolean) {
  const qc = useQueryClient();
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!enabled) return;

    let active = true;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

    function connect() {
      if (!active) return;
      const token = getAccessToken();
      const url = `${WS_BASE}/ws/notifications${
        token ? `?token=${encodeURIComponent(token)}` : ""
      }`;
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onmessage = () => {
        qc.invalidateQueries({ queryKey: ["notifications"] });
      };

      ws.onclose = () => {
        wsRef.current = null;
        if (active) {
          reconnectTimer = setTimeout(connect, 5_000);
        }
      };
    }

    connect();

    return () => {
      active = false;
      if (reconnectTimer) clearTimeout(reconnectTimer);
      wsRef.current?.close();
      wsRef.current = null;
    };
  }, [enabled, qc]);
}
