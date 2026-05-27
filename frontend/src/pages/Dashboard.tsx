import { Link } from "react-router-dom";
import { cities } from "@/data/mock";
import { Search, MapPin, Sparkles, Briefcase, Calendar } from "lucide-react";
import { Input } from "@/components/ui/input";
import { useMemo, useState, useEffect } from "react";
import { useApp } from "@/store/AppContext";
import { useQuery } from "@tanstack/react-query";
import { fetchTrips, searchTrips, type ApiTrip } from "@/lib/api";
import { format } from "date-fns";

const Dashboard = () => {
  const { user } = useApp();
  const [q, setQ] = useState("");
  const [debounced, setDebounced] = useState("");

  // 400 ms debounce for trip search
  useEffect(() => {
    const t = setTimeout(() => setDebounced(q.trim()), 400);
    return () => clearTimeout(t);
  }, [q]);

  const { data: trips = [] } = useQuery({
    queryKey: ["trips"],
    queryFn: fetchTrips,
  });

  const { data: tripResults = [], isFetching: searchingTrips } = useQuery({
    queryKey: ["search-trips", debounced],
    queryFn: () => searchTrips(debounced),
    enabled: debounced.length >= 2,
    staleTime: 30_000,
  });

  const filtered = useMemo(() => {
    if (!q.trim()) return cities;
    const s = q.toLowerCase();
    return cities.filter(c => c.name.toLowerCase().includes(s) || c.country.toLowerCase().includes(s));
  }, [q]);

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
          Browse curated destinations and plan your perfect trip day by day.
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

      {trips.length > 0 && !q.trim() && (
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

      {/* Trip search results */}
      {debounced.length >= 2 && (
        <section className="mt-10">
          <h2 className="mb-4 font-display text-2xl font-semibold tracking-tight">
            Your trips matching "{debounced}"
          </h2>
          {searchingTrips && (
            <div className="text-sm text-muted-foreground">Searching…</div>
          )}
          {!searchingTrips && tripResults.length === 0 && (
            <div className="rounded-3xl border border-dashed border-border/80 bg-card-grad p-10 text-center text-muted-foreground">
              No matching trips found.
            </div>
          )}
          {tripResults.length > 0 && (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {tripResults.map((trip: ApiTrip) => (
                <Link
                  key={trip.id}
                  to={`/trips/${trip.id}`}
                  className="group rounded-2xl border border-border/60 bg-card-grad p-5 transition-all hover:-translate-y-1 hover:shadow-glow"
                >
                  {trip.cover_image_url && (
                    <div className="mb-4 aspect-video overflow-hidden rounded-xl">
                      <img
                        src={trip.cover_image_url}
                        alt={trip.title}
                        className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                      />
                    </div>
                  )}
                  <div className="flex items-start justify-between gap-2">
                    <div>
                      <h3 className="font-display text-lg font-semibold leading-tight">{trip.title}</h3>
                      <div className="mt-0.5 flex items-center gap-1 text-xs text-muted-foreground">
                        <MapPin className="h-3.5 w-3.5" />{trip.destination}
                      </div>
                    </div>
                    <span className={`shrink-0 rounded-full px-2 py-0.5 text-xs capitalize ${
                      trip.status === "active" ? "bg-primary/15 text-primary" : "bg-muted text-muted-foreground"
                    }`}>{trip.status}</span>
                  </div>
                  {(trip.start_date || trip.end_date) && (
                    <div className="mt-2 flex items-center gap-1 text-xs text-muted-foreground">
                      <Calendar className="h-3.5 w-3.5" />
                      {trip.start_date && format(new Date(trip.start_date), "MMM d")}
                      {trip.start_date && trip.end_date && " – "}
                      {trip.end_date && format(new Date(trip.end_date), "MMM d, yyyy")}
                    </div>
                  )}
                </Link>
              ))}
            </div>
          )}
        </section>
      )}

      {/* Cities grid */}
      <section className="mt-12">
        <div className="mb-6 flex items-end justify-between">
          <h2 className="font-display text-2xl font-semibold tracking-tight md:text-3xl">Trending destinations</h2>
          <span className="text-sm text-muted-foreground">{filtered.length} cities</span>
        </div>

        {filtered.length === 0 ? (
          <div className="grid place-items-center rounded-3xl border border-dashed border-border/80 bg-card-grad p-16 text-center">
            <Briefcase className="h-10 w-10 text-muted-foreground" />
            <p className="mt-4 text-muted-foreground">No destinations match your search.</p>
          </div>
        ) : (
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
                </div>
              </Link>
            ))}
          </div>
        )}
      </section>
    </main>
  );
};

export default Dashboard;
