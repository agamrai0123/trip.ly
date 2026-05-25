import { Link, useLocation } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { Compass, MapPin, Briefcase, LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export const Header = () => {
  const { user, logout, trips } = useApp();
  const loc = useLocation();
  const onAuth = loc.pathname === "/";

  if (onAuth || !user) return null;

  const initials = user.name.split(" ").map(s => s[0]).slice(0, 2).join("").toUpperCase();

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
        </nav>

        <div className="flex items-center gap-3">
          <div className="hidden text-right sm:block">
            <div className="text-sm font-medium">{user.name}</div>
            <div className="text-xs text-muted-foreground">{user.email}</div>
          </div>
          <Avatar className="h-9 w-9 ring-2 ring-primary/30">
            <AvatarImage src={`https://api.dicebear.com/9.x/notionists/svg?seed=${user.name}`} alt={user.name} />
            <AvatarFallback>{initials}</AvatarFallback>
          </Avatar>
          <Button variant="ghost" size="icon" onClick={logout} aria-label="Sign out">
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
