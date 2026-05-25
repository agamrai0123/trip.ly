import { Link, useParams } from "react-router-dom";
import { cities } from "@/data/mock";
import { useQuery } from "@tanstack/react-query";
import { searchTrips } from "@/lib/api";
import { ArrowLeft, Calendar, DollarSign, MapPin, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { format } from "date-fns";

const CityDetail = () => {
  const { id } = useParams<{ id: string }>();
  const city = cities.find(c => c.id === id);

  const { data: trips = [], isLoading } = useQuery({
    queryKey: ["search-trips", city?.name],
    queryFn: () => searchTrips(city?.name ?? ""),
    enabled: Boolean(city?.name),
  });

  if (!city) {
    return (
      <main className="container py-20 text-center">
        <p className="text-muted-foreground">City not found.</p>
        <Link to="/dashboard" className="mt-4 inline-block text-primary underline">Back to dashboard</Link>
      </main>
    );
  }

  return (
    <main>
      <section className="relative h-[55vh] min-h-[420px] overflow-hidden">
        <img src={city.image} alt={city.name} width={1600} height={900} className="absolute inset-0 h-full w-full object-cover" />
        <div className="absolute inset-0 bg-gradient-to-t from-background via-background/70 to-background/20" />
        <div className="container relative flex h-full flex-col justify-end pb-10">
          <Link to="/dashboard" className="mb-4 inline-flex w-fit items-center gap-2 rounded-full bg-background/40 px-3 py-1.5 text-sm backdrop-blur hover:bg-background/60">
            <ArrowLeft className="h-4 w-4" /> Destinations
          </Link>
          <div className="flex items-center gap-2 text-primary">
            <MapPin className="h-4 w-4" />
            <span className="text-sm uppercase tracking-wider">{city.country}</span>
          </div>
          <h1 className="mt-1 font-display text-5xl font-semibold tracking-tight md:text-7xl">{city.name}</h1>
          <p className="mt-2 max-w-xl text-muted-foreground">{city.blurb}</p>
        </div>
      </section>

      <section className="container -mt-8 pb-20">
        <div className="mb-6 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-border/60 glass p-4">
          <div>
            <div className="font-display text-xl font-semibold">
              {isLoading ? "Loading..." : `${trips.length} trip${trips.length !== 1 ? "s" : ""} to ${city.name}`}
            </div>
            <div className="text-xs text-muted-foreground">Real trips planned by travelers in your community</div>
          </div>
          <Link to="/trips">
            <Button variant="secondary" size="sm">
              <Plus className="mr-1.5 h-4 w-4" /> Start planning
            </Button>
          </Link>
        </div>

        {isLoading ? (
          <div className="flex justify-center py-16">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
          </div>
        ) : trips.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-border p-14 text-center">
            <p className="text-muted-foreground">No trips to {city.name} yet.</p>
            <p className="mt-1 text-sm text-muted-foreground">Be the first to plan one.</p>
            <Link to="/trips" className="mt-4 inline-block">
              <Button className="bg-cta text-primary-foreground shadow-glow">Plan a trip here</Button>
            </Link>
          </div>
        ) : (
          <div className="space-y-3">
            {trips.map(trip => {
              const dateRange = trip.start_date && trip.end_date
                ? `${format(new Date(trip.start_date), "MMM d")} – ${format(new Date(trip.end_date), "MMM d, yyyy")}`
                : null;
              return (
                <article key={trip.id} className="group overflow-hidden rounded-2xl border border-border/60 bg-card-grad p-4 shadow-card transition hover:border-primary/40 hover:shadow-glow">
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex-1">
                      <h3 className="font-display text-xl font-semibold leading-snug tracking-tight group-hover:text-primary">{trip.title}</h3>
                      <div className="mt-2 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
                        {dateRange && (
                          <span className="inline-flex items-center gap-1.5"><Calendar className="h-4 w-4" />{dateRange}</span>
                        )}
                        {trip.budget_total > 0 && (
                          <span className="inline-flex items-center gap-1.5"><DollarSign className="h-4 w-4" />${trip.budget_total.toLocaleString()} {trip.currency}</span>
                        )}
                      </div>
                    </div>
                    <Link to={`/trips/${trip.id}`}>
                      <Button variant="secondary" size="sm">View</Button>
                    </Link>
                  </div>
                </article>
              );
            })}
          </div>
        )}
      </section>
    </main>
  );
};

export default CityDetail;
