import { Link } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchTrips, deleteTrip, createTrip, type ApiTrip, type CreateTripBody } from "@/lib/api";
import { Briefcase, Calendar, DollarSign, MapPin, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { format } from "date-fns";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const Trips = () => {
  const qc = useQueryClient();
  const nav = useNavigate();
  const [creating, setCreating] = useState(false);

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
        <Button onClick={() => setCreating(true)} className="bg-cta text-primary-foreground shadow-glow">
          <Plus className="mr-1.5 h-4 w-4" /> New trip
        </Button>
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
            Create your first trip and start planning.
          </p>
          <Button onClick={() => setCreating(true)} className="mt-6 bg-cta text-primary-foreground shadow-glow">
            <Plus className="mr-1.5 h-4 w-4" /> Create a trip
          </Button>
        </div>
      ) : (
        <div className="mt-10 grid grid-cols-1 gap-5 md:grid-cols-2">
          {trips.map(t => <TripCard key={t.id} trip={t} onDelete={() => remove.mutate(t.id)} />)}
        </div>
      )}

      {creating && (
        <CreateTripDialog
          onClose={() => setCreating(false)}
          onCreated={(id) => { nav(`/trips/${id}`); }}
          qc={qc}
        />
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

function CreateTripDialog({
  onClose,
  onCreated,
  qc,
}: {
  onClose: () => void;
  onCreated: (id: string) => void;
  qc: ReturnType<typeof useQueryClient>;
}) {
  const [title, setTitle] = useState("");
  const [destination, setDestination] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [budget, setBudget] = useState("");
  const [currency, setCurrency] = useState("USD");

  const create = useMutation({
    mutationFn: (body: CreateTripBody) => createTrip(body),
    onSuccess: (trip) => {
      qc.invalidateQueries({ queryKey: ["trips"] });
      toast.success("Trip created");
      onCreated(trip.id);
    },
    onError: () => toast.error("Could not create trip"),
  });

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !destination.trim()) {
      toast.error("Title and destination are required");
      return;
    }
    create.mutate({
      title: title.trim(),
      destination: destination.trim(),
      ...(startDate && { start_date: startDate }),
      ...(endDate && { end_date: endDate }),
      ...(budget && { budget_total: parseFloat(budget) }),
      currency,
    });
  };

  return (
    <Dialog open onOpenChange={(o) => { if (!o) onClose(); }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="font-display text-2xl">New trip</DialogTitle>
          <DialogDescription>Where are you going? You can add the full itinerary after.</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSubmit} className="space-y-4">
          <div>
            <Label htmlFor="ct-title">Trip name *</Label>
            <Input id="ct-title" value={title} onChange={e => setTitle(e.target.value)} placeholder="e.g. Japan Spring 2025" className="mt-1" autoFocus />
          </div>
          <div>
            <Label htmlFor="ct-dest">Destination *</Label>
            <Input id="ct-dest" value={destination} onChange={e => setDestination(e.target.value)} placeholder="e.g. Tokyo, Japan" className="mt-1" />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="ct-start">Start date</Label>
              <Input id="ct-start" type="date" value={startDate} onChange={e => setStartDate(e.target.value)} className="mt-1" />
            </div>
            <div>
              <Label htmlFor="ct-end">End date</Label>
              <Input id="ct-end" type="date" value={endDate} onChange={e => setEndDate(e.target.value)} className="mt-1" />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="ct-budget">Budget</Label>
              <Input id="ct-budget" type="number" min="0" step="1" value={budget} onChange={e => setBudget(e.target.value)} placeholder="0" className="mt-1" />
            </div>
            <div>
              <Label htmlFor="ct-currency">Currency</Label>
              <select id="ct-currency" value={currency} onChange={e => setCurrency(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-input bg-transparent px-3 text-sm">
                <option value="USD">USD</option>
                <option value="EUR">EUR</option>
                <option value="GBP">GBP</option>
                <option value="JPY">JPY</option>
                <option value="INR">INR</option>
                <option value="AUD">AUD</option>
              </select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="ghost" type="button" onClick={onClose}>Cancel</Button>
            <Button type="submit" disabled={create.isPending} className="bg-cta text-primary-foreground shadow-glow">
              {create.isPending ? "Creating..." : "Create trip"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

export default Trips;