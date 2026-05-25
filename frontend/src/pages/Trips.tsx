import { Link } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchTrips, deleteTrip, type ApiTrip } from "@/lib/api";
import { Briefcase, Calendar, DollarSign, MapPin, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { format } from "date-fns";

const Trips = () => {
  const qc = useQueryClient();
  const { data: trips = [], isLoading } = useQuery({
    queryKey: ["trips"],
    queryFn: fetchTrips,
  });

  const remove = useMutation({
    mutationFn: deleteTrip,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["trips"] });
      toast.success("Trip deleted");
    },
    onError: () => toast.error("Could not delete trip"),
  });

  return (
    <main className="container py-10 md:py-14">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="font-display text-4xl font-semibold tracking-tight md:text-5xl">My trips</h1>
          <p className="mt-2 text-muted-foreground">All your planned itineraries in one place.</p>
        </div>
        <Link to="/dashboard">
          <Button variant="secondary"><Plus className="mr-1.5 h-4 w-4" /> Plan another</Button>
        </Link>
      </div>

      {isLoading ? (
        <div className="mt-16 flex justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
        </div>
      ) : trips.length === 0 ? (
        <div className="mt-10 grid place-items-center rounded-3xl border border-dashed border-border/80 bg-card-grad p-16 text-center">
          <div className="grid h-14 w-14 place-items-center rounded-2xl bg-primary/15 text-primary">
            <Briefcase className="h-6 w-6" />
          </div>
          <h2 className="mt-5 font-display text-2xl">No trips yet</h2>
          <p className="mt-2 max-w-sm text-muted-foreground">
            Browse a destination to start planning your first trip.
          </p>
          <Link to="/dashboard"><Button className="mt-6 bg-cta text-primary-foreground shadow-glow">Discover cities</Button></Link>
        </div>
      ) : (
        <div className="mt-10 grid grid-cols-1 gap-5 md:grid-cols-2">
          {trips.map(t => <TripCard key={t.id} trip={t} onDelete={() => remove.mutate(t.id)} />)}
        </div>
      )}
    </main>
  );
};

function TripCard({ trip, onDelete }: { trip: ApiTrip; onDelete: () => void }) {
  const dateRange = trip.start_date && trip.end_date
    ? `${format(new Date(trip.start_date), "MMM d")} – ${format(new Date(trip.end_date), "MMM d, yyyy")}`
    : trip.start_date
    ? `From ${format(new Date(trip.start_date), "MMM d, yyyy")}`
    : null;

  return (
    <div className="group overflow-hidden rounded-3xl border border-border/60 bg-card shadow-card transition hover:-translate-y-0.5 hover:shadow-glow">
      {trip.cover_image_url && (
        <div className="relative aspect-[16/9] overflow-hidden">
          <img src={trip.cover_image_url} alt={trip.title} className="h-full w-full object-cover transition-transform duration-700 group-hover:scale-105" />
          <div className="absolute inset-0 bg-gradient-to-t from-background via-background/40 to-transparent" />
          <div className="absolute inset-x-0 bottom-0 p-5">
            <div className="flex items-center gap-1.5 text-xs text-primary">
              <MapPin className="h-3.5 w-3.5" />
              <span className="uppercase tracking-wider">{trip.destination}</span>
            </div>
          </div>
        </div>
      )}
      <div className="p-5">
        <div className="flex items-start justify-between gap-2">
          <div>
            <h3 className="font-display text-xl font-semibold leading-tight">{trip.title}</h3>
            {!trip.cover_image_url && (
              <p className="mt-0.5 text-sm text-muted-foreground">{trip.destination}</p>
            )}
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={e => { e.preventDefault(); onDelete(); }}
            className="shrink-0 text-muted-foreground hover:text-destructive"
            aria-label="Delete trip"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>

        <div className="mt-4 flex flex-wrap gap-3 text-sm text-muted-foreground">
          {dateRange && (
            <span className="inline-flex items-center gap-1.5"><Calendar className="h-4 w-4" />{dateRange}</span>
          )}
          {trip.budget_total > 0 && (
            <span className="inline-flex items-center gap-1.5"><DollarSign className="h-4 w-4" />${trip.budget_total.toLocaleString()} {trip.currency}</span>
          )}
        </div>

        <div className="mt-4">
          <Link to={`/trips/${trip.id}`}>
            <Button size="sm" variant="secondary" className="w-full">View itinerary →</Button>
          </Link>
        </div>
      </div>
    </div>
  );
}

export default Trips;
