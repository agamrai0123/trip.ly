import { Link } from "react-router-dom";
import { cities, tripPosts } from "@/data/mock";
import { Search, MapPin, Sparkles } from "lucide-react";
import { Input } from "@/components/ui/input";
import { useMemo, useState } from "react";
import { useApp } from "@/store/AppContext";

const Dashboard = () => {
  const { user, trips } = useApp();
  const [q, setQ] = useState("");

  const filtered = useMemo(() => {
    if (!q.trim()) return cities;
    const s = q.toLowerCase();
    return cities.filter(c => c.name.toLowerCase().includes(s) || c.country.toLowerCase().includes(s));
  }, [q]);

  const totalPosts = tripPosts.length;

  return (
    <main className="container py-10 md:py-14">
      {/* Hero greeting */}
      <section className="animate-fade-up">
        <div className="flex items-center gap-2 text-sm text-primary">
          <Sparkles className="h-4 w-4" />
          <span>Welcome back, {user?.name.split(" ")[0]}</span>
        </div>
        <h1 className="mt-3 font-display text-4xl font-semibold tracking-tight md:text-6xl">
          Where to next?
        </h1>
        <p className="mt-3 max-w-xl text-muted-foreground">
          Browse {totalPosts}+ traveler-built itineraries. Pick a city to see what real trips look like — day by day, dollar by dollar.
        </p>

        <div className="mt-8 max-w-xl">
          <div className="relative">
            <Search className="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
            <Input
              value={q}
              onChange={e => setQ(e.target.value)}
              placeholder="Search Tokyo, Bali, Lisbon…"
              className="h-14 rounded-2xl border-border/80 bg-card/60 pl-12 text-base"
            />
          </div>
        </div>
      </section>

      {trips.length > 0 && (
        <section className="mt-10 rounded-2xl border border-border/60 bg-card-grad p-5">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-xs uppercase tracking-wider text-muted-foreground">In progress</div>
              <div className="mt-1 font-display text-xl">You have {trips.length} active trip{trips.length > 1 ? "s" : ""}</div>
            </div>
            <Link to="/trips" className="text-sm font-medium text-primary hover:underline">View all →</Link>
          </div>
        </section>
      )}

      {/* Cities grid */}
      <section className="mt-12">
        <div className="mb-6 flex items-end justify-between">
          <h2 className="font-display text-2xl font-semibold tracking-tight md:text-3xl">Trending destinations</h2>
          <span className="text-sm text-muted-foreground">{filtered.length} cities</span>
        </div>

        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {filtered.map((city, i) => (
            <Link
              key={city.id}
              to={`/city/${city.id}`}
              className="group relative overflow-hidden rounded-3xl border border-border/60 bg-card shadow-card transition-all hover:-translate-y-1 hover:shadow-glow"
              style={{ animation: `fade-up 0.5s var(--transition-smooth) ${i * 60}ms both` }}
            >
              <div className="relative aspect-[4/5] overflow-hidden">
                <img
                  src={city.image}
                  alt={`${city.name}, ${city.country}`}
                  loading="lazy"
                  width={1024}
                  height={768}
                  className="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-background via-background/40 to-transparent" />
              </div>

              <div className="absolute inset-x-0 bottom-0 p-6">
                <div className="flex items-center gap-1.5 text-xs text-primary">
                  <MapPin className="h-3.5 w-3.5" />
                  <span className="uppercase tracking-wider">{city.country}</span>
                </div>
                <h3 className="mt-1 font-display text-3xl font-semibold leading-tight tracking-tight">
                  {city.name}
                </h3>
                <p className="mt-1 text-sm text-muted-foreground">{city.blurb}</p>
                <div className="mt-3 inline-flex items-center gap-2 text-xs text-foreground/80">
                  <span className="inline-block h-1.5 w-1.5 rounded-full bg-accent" />
                  {city.postsCount} traveler itineraries
                </div>
              </div>
            </Link>
          ))}
        </div>

        {filtered.length === 0 && (
          <div className="rounded-2xl border border-dashed border-border p-10 text-center text-muted-foreground">
            No cities match “{q}”. Try Tokyo, Bali, Kyoto…
          </div>
        )}
      </section>
    </main>
  );
};

export default Dashboard;
