import { Link, useLocation, useNavigate } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchTrips, fetchNotifications, markAllNotificationsRead, markNotificationRead, type ApiNotification } from "@/lib/api";
import { useNotificationsWS } from "@/hooks/useNotificationsWS";
import { useTheme } from "next-themes";
import { Bell, Compass, MapPin, Briefcase, Users2, LogOut, User, Settings, Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { formatDistanceToNow } from "date-fns";

export const Header = () => {
  const { user, logout } = useApp();
  const loc = useLocation();
  const nav = useNavigate();
  const qc = useQueryClient();
  const { theme, setTheme } = useTheme();
  const onAuth = loc.pathname === "/";

  // Real-time notifications via WebSocket (supplements 30s polling)
  useNotificationsWS(Boolean(user));

  const { data: trips = [] } = useQuery({
    queryKey: ["trips"],
    queryFn: fetchTrips,
    enabled: Boolean(user),
  });

  const { data: notifications = [] } = useQuery({
    queryKey: ["notifications"],
    queryFn: fetchNotifications,
    enabled: Boolean(user),
    refetchInterval: 30_000,
  });

  const unread = notifications.filter((n: ApiNotification) => !n.read);

  const markRead = useMutation({
    mutationFn: (id: string) => markNotificationRead(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notifications"] }),
  });

  const markAll = useMutation({
    mutationFn: markAllNotificationsRead,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notifications"] }),
  });

  if (onAuth || !user) return null;

  const initials = user.name.split(" ").map((s: string) => s[0]).slice(0, 2).join("").toUpperCase();

  const handleLogout = async () => {
    await logout();
    nav("/", { replace: true });
  };

  return (
    <header className="sticky top-0 z-40 glass">
      <div className="container flex h-16 items-center justify-between">
        <Link to="/dashboard" className="flex items-center gap-2 group">
          <div className="grid h-9 w-9 place-items-center rounded-xl bg-cta shadow-glow">
            <Compass className="h-5 w-5 text-primary-foreground" />
          </div>
          <span className="font-display text-xl font-semibold tracking-tight">wandr</span>
        </Link>

        <nav className="hidden items-center gap-1 md:flex">
          <NavBtn to="/dashboard" icon={<MapPin className="h-4 w-4" />} active={loc.pathname.startsWith("/dashboard") || loc.pathname.startsWith("/city")}>
            Discover
          </NavBtn>
          <NavBtn to="/trips" icon={<Briefcase className="h-4 w-4" />} active={loc.pathname.startsWith("/trips")}>
            My Trips {trips.length > 0 && <span className="ml-1 rounded-full bg-primary/20 px-1.5 text-xs text-primary">{trips.length}</span>}
          </NavBtn>
          <NavBtn to="/collaborations" icon={<Users2 className="h-4 w-4" />} active={loc.pathname.startsWith("/collaborations")}>
            Collaborate
          </NavBtn>
          <NavBtn to="/settings" icon={<Settings className="h-4 w-4" />} active={loc.pathname === "/settings"}>
            Settings
          </NavBtn>
        </nav>

        <div className="flex items-center gap-2">
          {/* Dark / Light theme toggle */}
          <button
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            className="rounded-full p-2 text-muted-foreground transition hover:bg-secondary focus:outline-none focus-visible:ring-2 focus-visible:ring-primary"
            aria-label="Toggle theme"
          >
            {theme === "dark" ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
          </button>

          {/* Notifications bell */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button className="relative rounded-full p-2 text-muted-foreground transition hover:bg-secondary focus:outline-none focus-visible:ring-2 focus-visible:ring-primary" aria-label="Notifications">
                <Bell className="h-5 w-5" />
                {unread.length > 0 && (
                  <span className="absolute -right-0.5 -top-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-destructive text-[10px] font-bold text-white">
                    {unread.length > 9 ? "9+" : unread.length}
                  </span>
                )}
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-80 max-h-96 overflow-y-auto">
              <div className="flex items-center justify-between px-3 py-2">
                <span className="text-sm font-semibold">Notifications</span>
                {unread.length > 0 && (
                  <button onClick={() => markAll.mutate()} className="text-xs text-primary hover:underline">
                    Mark all read
                  </button>
                )}
              </div>
              <DropdownMenuSeparator />
              {notifications.length === 0 ? (
                <div className="px-3 py-6 text-center text-xs text-muted-foreground">No notifications yet.</div>
              ) : (
                notifications.slice(0, 20).map((n: ApiNotification) => (
                  <DropdownMenuItem
                    key={n.id}
                    onClick={() => { if (!n.read) markRead.mutate(n.id); }}
                    className={`flex flex-col items-start gap-0.5 px-3 py-2 ${!n.read ? "bg-primary/5" : ""}`}
                  >
                    <span className="text-xs font-medium capitalize">{n.type.replace(/_/g, " ")}</span>
                    <span className="text-xs text-muted-foreground">
                      {formatDistanceToNow(new Date(n.created_at), { addSuffix: true })}
                    </span>
                  </DropdownMenuItem>
                ))
              )}
            </DropdownMenuContent>
          </DropdownMenu>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button className="flex items-center gap-2 rounded-full pr-1 transition hover:bg-secondary focus:outline-none focus-visible:ring-2 focus-visible:ring-primary" aria-label="Account menu">
                <div className="hidden text-right sm:block">
                  <div className="text-sm font-medium leading-tight">{user.name}</div>
                  <div className="text-xs text-muted-foreground">{user.email}</div>
                </div>
                <Avatar className="h-9 w-9 ring-2 ring-primary/30">
                  <AvatarImage src={user.avatar_url || `https://api.dicebear.com/9.x/notionists/svg?seed=${user.name}`} alt={user.name} />
                  <AvatarFallback>{initials}</AvatarFallback>
                </Avatar>
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuItem asChild>
                <Link to="/profile" className="flex items-center gap-2">
                  <User className="h-4 w-4" /> Profile
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={handleLogout} className="flex items-center gap-2 text-destructive focus:text-destructive">
                <LogOut className="h-4 w-4" /> Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Mobile sign-out fallback */}
          <Button variant="ghost" size="icon" onClick={handleLogout} aria-label="Sign out" className="md:hidden">
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </header>
  );
};

const NavBtn = ({ to, icon, active, children }: { to: string; icon: React.ReactNode; active: boolean; children: React.ReactNode }) => (
  <Link
    to={to}
    className={`inline-flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition-colors ${
      active ? "bg-primary/15 text-primary" : "text-muted-foreground hover:bg-secondary hover:text-foreground"
    }`}
  >
    {icon}
    {children}
  </Link>
);
