import { Link, useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, useEffect } from "react";
import {
  fetchTrip, fetchTripDays, fetchDayItems, reorderItems, createDayItem, deleteItem,
  type ApiTrip, type ApiDay, type ApiItem, type CreateItemBody,
} from "@/lib/api";
import { ArrowLeft, Calendar, DollarSign, GripVertical, MapPin, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { format } from "date-fns";
import {
  DndContext, type DragEndEvent, PointerSensor, closestCenter, useSensor, useSensors,
} from "@dnd-kit/core";
import {
  SortableContext, arrayMove, useSortable, verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

const typeBadgeClass: Record<string, string> = {
  activity: "bg-sky-500/15 text-sky-400",
  food: "bg-orange-500/15 text-orange-400",
  accommodation: "bg-violet-500/15 text-violet-400",
  transport: "bg-emerald-500/15 text-emerald-400",
  default: "bg-muted text-muted-foreground",
};

const TripDetail = () => {
  const { id } = useParams<{ id: string }>();
  const qc = useQueryClient();

  const { data: trip, isLoading: loadingTrip } = useQuery({
    queryKey: ["trip", id],
    queryFn: () => fetchTrip(id as string),
    enabled: Boolean(id),
  });

  const { data: days = [], isLoading: loadingDays } = useQuery({
    queryKey: ["days", id],
    queryFn: () => fetchTripDays(id as string),
    enabled: Boolean(id),
  });

  const [addingDay, setAddingDay] = useState<ApiDay | null>(null);

  if (loadingTrip || loadingDays) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }

  if (!trip) {
    return (
      <main className="container py-20 text-center">
        <p className="text-muted-foreground">Trip not found.</p>
        <Link to="/trips" className="mt-4 inline-block text-primary underline">Back to trips</Link>
      </main>
    );
  }

  const dateRange = trip.start_date && trip.end_date
    ? `${format(new Date(trip.start_date), "MMM d")} – ${format(new Date(trip.end_date), "MMM d, yyyy")}`
    : null;

  return (
    <main>
      <section className="relative h-[38vh] min-h-[300px] overflow-hidden">
        {trip.cover_image_url && (
          <img src={trip.cover_image_url} alt={trip.title} width={1600} height={900} className="absolute inset-0 h-full w-full object-cover" />
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-background via-background/70 to-background/30" />
        <div className="container relative flex h-full flex-col justify-end pb-10">
          <Link to="/trips" className="mb-3 inline-flex w-fit items-center gap-2 rounded-full bg-background/40 px-3 py-1.5 text-sm backdrop-blur hover:bg-background/60">
            <ArrowLeft className="h-4 w-4" /> My trips
          </Link>
          <div className="flex items-center gap-2 text-primary">
            <MapPin className="h-4 w-4" />
            <span className="text-sm uppercase tracking-wider">{trip.destination}</span>
          </div>
          <h1 className="mt-1 font-display text-5xl font-semibold tracking-tight md:text-6xl">{trip.title}</h1>
          {dateRange && <p className="mt-1 text-muted-foreground">{dateRange}</p>}
        </div>
      </section>

      <section className="container -mt-4 pb-4">
        <div className="glass rounded-2xl border border-border/60 p-5 flex flex-wrap gap-6">
          {trip.budget_total > 0 && (
            <div>
              <div className="text-xs uppercase tracking-wider text-muted-foreground">Budget</div>
              <div className="mt-0.5 flex items-center gap-1 font-display text-2xl font-semibold">
                <DollarSign className="h-5 w-5 text-muted-foreground" />{trip.budget_total.toLocaleString()} {trip.currency}
              </div>
            </div>
          )}
          <div>
            <div className="text-xs uppercase tracking-wider text-muted-foreground">Days</div>
            <div className="mt-0.5 flex items-center gap-1 font-display text-2xl font-semibold">
              <Calendar className="h-5 w-5 text-muted-foreground" />{days.length}
            </div>
          </div>
          <div className="ml-auto self-center">
            <Button onClick={() => setAddingDay(days[0] ?? null)} variant="secondary" size="sm" disabled={days.length === 0}>
              <Plus className="mr-1.5 h-4 w-4" /> Add item
            </Button>
          </div>
        </div>
      </section>

      <section className="container pb-20 space-y-10">
        {days.length === 0 ? (
          <div className="rounded-3xl border border-dashed border-border/80 bg-card-grad p-16 text-center text-muted-foreground">
            No days planned yet.
          </div>
        ) : (
          days.map(day => (
            <DaySection key={day.id} trip={trip} day={day} onAddItem={() => setAddingDay(day)} qc={qc} />
          ))
        )}
      </section>

      {addingDay && (
        <AddItemDialog tripId={trip.id} day={addingDay} days={days} onClose={() => setAddingDay(null)} qc={qc} />
      )}
    </main>
  );
};

function DaySection({ trip, day, onAddItem, qc }: {
  trip: ApiTrip;
  day: ApiDay;
  onAddItem: () => void;
  qc: ReturnType<typeof useQueryClient>;
}) {
  const { data: items = [] } = useQuery({
    queryKey: ["items", day.id],
    queryFn: () => fetchDayItems(trip.id, day.id),
  });

  const [ordered, setOrdered] = useState<ApiItem[]>([]);

  useEffect(() => {
    setOrdered([...items].sort((a, b) => a.order_index - b.order_index));
  }, [items]);

  const reorder = useMutation({
    mutationFn: (ids: string[]) => reorderItems(trip.id, ids),
    onError: () => {
      setOrdered([...items].sort((a, b) => a.order_index - b.order_index));
      toast.error("Reorder failed");
    },
  });

  const remove = useMutation({
    mutationFn: (itemId: string) => deleteItem(trip.id, itemId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["items", day.id] });
      toast.success("Removed");
    },
    onError: () => toast.error("Could not remove item"),
  });

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

  const handleDragEnd = (e: DragEndEvent) => {
    const { active, over } = e;
    if (!over || active.id === over.id) return;
    const oldIdx = ordered.findIndex(i => i.id === active.id);
    const newIdx = ordered.findIndex(i => i.id === over.id);
    if (oldIdx < 0 || newIdx < 0) return;
    const reorderedArr = arrayMove(ordered, oldIdx, newIdx);
    setOrdered(reorderedArr);
    reorder.mutate(reorderedArr.map(i => i.id));
  };

  const dateLabel = day.date ? format(new Date(day.date), "EEE, MMM d") : `Day ${day.day_number}`;

  return (
    <div>
      <div className="mb-4 flex items-end justify-between border-b border-border/60 pb-2">
        <div>
          <div className="text-xs uppercase tracking-wider text-primary">Day {day.day_number}</div>
          <h2 className="font-display text-3xl font-semibold tracking-tight">{dateLabel}</h2>
          {day.notes && <p className="mt-0.5 text-sm text-muted-foreground">{day.notes}</p>}
        </div>
        <Button variant="ghost" size="sm" onClick={onAddItem}><Plus className="mr-1 h-4 w-4" />Add</Button>
      </div>

      {ordered.length === 0 ? (
        <div className="rounded-xl border border-dashed border-border/60 px-4 py-5 text-center text-xs text-muted-foreground">
          No stops planned. Click Add to build out this day.
        </div>
      ) : (
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={ordered.map(i => i.id)} strategy={verticalListSortingStrategy}>
            <div className="space-y-2">
              {ordered.map(item => (
                <SortableItem key={item.id} item={item} onRemove={() => remove.mutate(item.id)} />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      )}
    </div>
  );
}

function SortableItem({ item, onRemove }: { item: ApiItem; onRemove: () => void }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: item.id });
  const style = { transform: CSS.Transform.toString(transform), transition, opacity: isDragging ? 0.55 : 1 };
  const badgeClass = typeBadgeClass[item.type ?? ""] ?? typeBadgeClass.default;
  const timeRange = item.start_time && item.end_time
    ? `${item.start_time} – ${item.end_time}`
    : item.start_time ?? null;

  return (
    <div ref={setNodeRef} style={style} className="flex gap-2 overflow-hidden rounded-2xl border border-border/60 bg-card-grad">
      <button
        {...attributes}
        {...listeners}
        className="flex w-8 cursor-grab items-center justify-center text-muted-foreground hover:bg-secondary active:cursor-grabbing"
        aria-label="Drag to reorder"
      >
        <GripVertical className="h-4 w-4" />
      </button>
      <div className="flex flex-1 items-start justify-between gap-3 py-3 pr-3">
        <div className="flex-1">
          <div className="flex flex-wrap items-center gap-2">
            {item.type && (
              <span className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium capitalize ${badgeClass}`}>
                {item.type}
              </span>
            )}
            {timeRange && <span className="text-xs text-muted-foreground">{timeRange}</span>}
          </div>
          <h4 className="mt-1 font-semibold leading-tight">{item.title}</h4>
          {item.location && <div className="mt-0.5 text-xs text-muted-foreground">{item.location}</div>}
          {item.description && <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">{item.description}</p>}
        </div>
        <div className="flex flex-col items-end gap-2">
          {item.cost > 0 && (
            <span className="font-display text-base font-semibold">${item.cost.toLocaleString()}</span>
          )}
          <Button size="icon" variant="ghost" onClick={onRemove} aria-label="Remove item" className="h-7 w-7 text-muted-foreground hover:text-destructive">
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

function AddItemDialog({ tripId, day, days, onClose, qc }: {
  tripId: string;
  day: ApiDay;
  days: ApiDay[];
  onClose: () => void;
  qc: ReturnType<typeof useQueryClient>;
}) {
  const [selectedDayId, setSelectedDayId] = useState(day.id);
  const [title, setTitle] = useState("");
  const [type, setType] = useState("activity");
  const [location, setLocation] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [cost, setCost] = useState("");

  const add = useMutation({
    mutationFn: (body: CreateItemBody) => createDayItem(tripId, selectedDayId, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["items", selectedDayId] });
      toast.success("Item added");
      onClose();
    },
    onError: () => toast.error("Could not add item"),
  });

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) { toast.error("Enter a title"); return; }
    const body: CreateItemBody = {
      title: title.trim(),
      type,
      ...(location && { location }),
      ...(startTime && { start_time: startTime }),
      ...(endTime && { end_time: endTime }),
      ...(cost && { cost: parseFloat(cost) }),
    };
    add.mutate(body);
  };

  return (
    <Dialog open onOpenChange={(o) => { if (!o) onClose(); }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="font-display text-2xl">Add a stop</DialogTitle>
          <DialogDescription>Fill in the details for this activity, meal, or place.</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSubmit} className="space-y-4">
          <div>
            <Label htmlFor="ai-day">Day</Label>
            <select id="ai-day" value={selectedDayId} onChange={e => setSelectedDayId(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-input bg-transparent px-3 text-sm">
              {days.map(d => (
                <option key={d.id} value={d.id}>
                  Day {d.day_number}{d.date ? ` — ${format(new Date(d.date), "MMM d")}` : ""}
                </option>
              ))}
            </select>
          </div>
          <div>
            <Label htmlFor="ai-title">Title *</Label>
            <Input id="ai-title" value={title} onChange={e => setTitle(e.target.value)} placeholder="e.g. Senso-ji Temple" className="mt-1" autoFocus />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="ai-type">Type</Label>
              <select id="ai-type" value={type} onChange={e => setType(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-input bg-transparent px-3 text-sm">
                <option value="activity">Activity</option>
                <option value="food">Food</option>
                <option value="accommodation">Accommodation</option>
                <option value="transport">Transport</option>
              </select>
            </div>
            <div>
              <Label htmlFor="ai-cost">Cost ($)</Label>
              <Input id="ai-cost" type="number" min="0" step="0.01" value={cost} onChange={e => setCost(e.target.value)} placeholder="0" className="mt-1" />
            </div>
          </div>
          <div>
            <Label htmlFor="ai-location">Location</Label>
            <Input id="ai-location" value={location} onChange={e => setLocation(e.target.value)} placeholder="Address or area" className="mt-1" />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="ai-start">Start time</Label>
              <Input id="ai-start" type="time" value={startTime} onChange={e => setStartTime(e.target.value)} className="mt-1" />
            </div>
            <div>
              <Label htmlFor="ai-end">End time</Label>
              <Input id="ai-end" type="time" value={endTime} onChange={e => setEndTime(e.target.value)} className="mt-1" />
            </div>
          </div>
          <DialogFooter>
            <Button variant="ghost" type="button" onClick={onClose}>Cancel</Button>
            <Button type="submit" disabled={add.isPending} className="bg-cta text-primary-foreground shadow-glow">
              {add.isPending ? "Adding..." : "Add stop"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

export default TripDetail;
