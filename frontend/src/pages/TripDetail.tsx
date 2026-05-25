import { Link, useParams } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { useMemo, useState } from "react";
import { cities } from "@/data/mock";
import type { ItineraryItem, Slot } from "@/types";
import { ArrowLeft, Check, GripVertical, MapPin, Pencil, Star, Trash2, CreditCard, Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { toast } from "sonner";
import {
  DndContext,
  DragEndEvent,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

const slotMeta: Record<Slot, { label: string; emoji: string; time: string }> = {
  morning: { label: "Morning", emoji: "🌅", time: "8:00 – 11:30" },
  lunch: { label: "Afternoon", emoji: "🍜", time: "12:00 – 16:00" },
  evening: { label: "Evening", emoji: "🌇", time: "16:00 – 19:00" },
  night: { label: "Night", emoji: "🌙", time: "20:00 – late" },
};
const slotOrder: Slot[] = ["morning", "lunch", "evening", "night"];

const TripDetail = () => {
  const { id } = useParams<{ id: string }>();
  const { trips, toggleVisited, removeItem, updateItemExpense, reviewPlace, reorderTrip } = useApp();
  const trip = trips.find(t => t.id === id);
  const city = cities.find(c => c.id === trip?.cityId);

  const [reviewing, setReviewing] = useState<ItineraryItem | null>(null);
  const [reviewText, setReviewText] = useState("");
  const [reviewRating, setReviewRating] = useState(5);
  const [editingExp, setEditingExp] = useState<string | null>(null);
  const [expDraft, setExpDraft] = useState("");

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

  const total = useMemo(() => trip?.items.reduce((s, i) => s + i.place.expense, 0) ?? 0, [trip]);

  const grouped = useMemo(() => {
    if (!trip) return [] as { day: number; slots: { slot: Slot; items: ItineraryItem[] }[] }[];
    const dayMap = new Map<number, ItineraryItem[]>();
    trip.items.forEach(it => {
      if (!dayMap.has(it.dayIndex)) dayMap.set(it.dayIndex, []);
      dayMap.get(it.dayIndex)!.push(it);
    });
    const days = Array.from({ length: trip.days }, (_, d) => d);
    dayMap.forEach((_v, k) => { if (!days.includes(k)) days.push(k); });
    return days.sort((a, b) => a - b).map(day => ({
      day,
      slots: slotOrder.map(slot => ({
        slot,
        items: (dayMap.get(day) ?? []).filter(it => it.slot === slot),
      })),
    }));
  }, [trip]);

  if (!trip) {
    return (
      <main className="container py-20 text-center">
        <p className="text-muted-foreground">Trip not found.</p>
        <Link to="/trips" className="mt-4 inline-block text-primary underline">Back to trips</Link>
      </main>
    );
  }

  const visitedCount = trip.items.filter(i => i.visited).length;
  const progress = Math.round((visitedCount / Math.max(trip.items.length, 1)) * 100);

  const openReview = (it: ItineraryItem) => {
    setReviewing(it);
    setReviewText("");
    setReviewRating(5);
  };
  const submitReview = () => {
    if (!reviewing || !reviewText.trim()) {
      toast.error("Add a few words about the place");
      return;
    }
    reviewPlace(trip.id, reviewing.id, reviewText.trim(), reviewRating);
    toast.success("Review saved & shared as a post 🎉");
    setReviewing(null);
  };

  const handleDragEnd = (day: number, slot: Slot) => (e: DragEndEvent) => {
    const { active, over } = e;
    if (!over || active.id === over.id) return;
    // Get current items array; reorder within this slot/day
    const currentSlotItems = trip.items.filter(it => it.dayIndex === day && it.slot === slot);
    const oldIndex = currentSlotItems.findIndex(i => i.id === active.id);
    const newIndex = currentSlotItems.findIndex(i => i.id === over.id);
    if (oldIndex < 0 || newIndex < 0) return;
    const reorderedSlot = arrayMove(currentSlotItems, oldIndex, newIndex);
    // Rebuild full items array
    const others = trip.items.filter(it => !(it.dayIndex === day && it.slot === slot));
    reorderTrip(trip.id, [...others, ...reorderedSlot]);
  };

  return (
    <main>
      <section className="relative h-[40vh] min-h-[320px] overflow-hidden">
        {city && <img src={city.image} alt={city.name} width={1600} height={900} className="absolute inset-0 h-full w-full object-cover" />}
        <div className="absolute inset-0 bg-gradient-to-t from-background via-background/70 to-background/30" />
        <div className="container relative flex h-full flex-col justify-end pb-10">
          <Link to="/trips" className="mb-3 inline-flex w-fit items-center gap-2 rounded-full bg-background/40 px-3 py-1.5 text-sm backdrop-blur hover:bg-background/60">
            <ArrowLeft className="h-4 w-4" /> My trips
          </Link>
          <div className="flex items-center gap-2 text-primary">
            <MapPin className="h-4 w-4" />
            <span className="text-sm uppercase tracking-wider">{city?.country}</span>
          </div>
          <h1 className="mt-1 font-display text-5xl font-semibold tracking-tight md:text-6xl">{city?.name}</h1>
          <p className="mt-1 text-muted-foreground">{trip.days} days · {trip.items.length} stops</p>
        </div>
      </section>

      <section className="container -mt-6 pb-20">
        <div className="grid gap-3 rounded-3xl border border-border/60 glass p-5 md:grid-cols-[1fr_1fr_1fr_auto]">
          <Summary label="Days" value={trip.days} />
          <Summary label="Total budget" value={`$${total.toLocaleString()}`} />
          <div>
            <div className="text-xs uppercase tracking-wider text-muted-foreground">Progress</div>
            <div className="mt-1 font-display text-2xl font-semibold">{progress}%</div>
            <div className="mt-1.5 h-1.5 overflow-hidden rounded-full bg-secondary">
              <div className="h-full bg-cta transition-all" style={{ width: `${progress}%` }} />
            </div>
          </div>
          <Button size="lg" className="bg-cta text-primary-foreground shadow-glow hover:opacity-95" onClick={() => toast.message("Checkout coming soon", { description: `We'll be able to book your $${total.toLocaleString()} itinerary in one click.` })}>
            <CreditCard className="mr-1.5 h-4 w-4" /> Book whole trip · ${total.toLocaleString()}
          </Button>
        </div>

        <div className="mt-3 text-xs text-muted-foreground">
          Tip: drag <GripVertical className="inline h-3 w-3" /> to reorder items within a time block.
        </div>

        <div className="mt-8 space-y-10">
          {grouped.map(({ day, slots }) => {
            const dayItems = slots.flatMap(s => s.items);
            const dayTotal = dayItems.reduce((s, i) => s + i.place.expense, 0);
            return (
              <div key={day}>
                <div className="mb-4 flex items-end justify-between border-b border-border/60 pb-2">
                  <div>
                    <div className="text-xs uppercase tracking-wider text-primary">Day {day + 1}</div>
                    <h2 className="font-display text-3xl font-semibold tracking-tight">A day in {city?.name}</h2>
                  </div>
                  <div className="text-right text-sm text-muted-foreground">
                    Day total<br /><span className="font-display text-xl text-foreground">${dayTotal}</span>
                  </div>
                </div>

                <div className="space-y-6">
                  {slots.map(({ slot, items }) => {
                    const meta = slotMeta[slot];
                    return (
                      <div key={slot}>
                        <div className="mb-2 flex items-center gap-2">
                          <span className="text-lg">{meta.emoji}</span>
                          <h3 className="font-display text-lg font-semibold">{meta.label}</h3>
                          <span className="text-xs text-muted-foreground">{meta.time}</span>
                        </div>
                        {items.length === 0 ? (
                          <div className="rounded-xl border border-dashed border-border/60 px-4 py-3 text-xs text-muted-foreground">
                            Nothing planned. Right-click a place from a post to drop one here.
                          </div>
                        ) : (
                          <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd(day, slot)}>
                            <SortableContext items={items.map(i => i.id)} strategy={verticalListSortingStrategy}>
                              <div className="space-y-2">
                                {items.map(it => (
                                  <SortableRow
                                    key={it.id}
                                    item={it}
                                    isEditingExp={editingExp === it.id}
                                    expDraft={expDraft}
                                    setExpDraft={setExpDraft}
                                    onStartEdit={() => { setEditingExp(it.id); setExpDraft(String(it.place.expense)); }}
                                    onSaveExp={(n) => { updateItemExpense(trip.id, it.id, n); setEditingExp(null); toast.success("Updated"); }}
                                    onToggleVisited={() => toggleVisited(trip.id, it.id)}
                                    onRemove={() => { removeItem(trip.id, it.id); toast.success("Removed"); }}
                                    onReview={() => openReview(it)}
                                  />
                                ))}
                              </div>
                            </SortableContext>
                          </DndContext>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
      </section>

      <Dialog open={!!reviewing} onOpenChange={(o) => !o && setReviewing(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="font-display text-2xl">Review {reviewing?.place.name}</DialogTitle>
            <DialogDescription>Your review becomes a tiny post other travelers can see.</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <div className="mb-2 text-sm font-medium">Rating</div>
              <div className="flex gap-1">
                {[1, 2, 3, 4, 5].map(n => (
                  <button key={n} onClick={() => setReviewRating(n)} aria-label={`${n} stars`}>
                    <Star className={`h-7 w-7 ${n <= reviewRating ? "fill-accent text-accent" : "text-muted-foreground"}`} />
                  </button>
                ))}
              </div>
            </div>
            <div>
              <div className="mb-2 text-sm font-medium">How was it?</div>
              <Textarea value={reviewText} onChange={e => setReviewText(e.target.value)} placeholder="A line or two — what you loved, what to know…" rows={4} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="ghost" onClick={() => setReviewing(null)}>Cancel</Button>
            <Button className="bg-cta text-primary-foreground" onClick={submitReview}>Save & share</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </main>
  );
};

const SortableRow = ({
  item, isEditingExp, expDraft, setExpDraft, onStartEdit, onSaveExp, onToggleVisited, onRemove, onReview,
}: {
  item: ItineraryItem;
  isEditingExp: boolean;
  expDraft: string;
  setExpDraft: (v: string) => void;
  onStartEdit: () => void;
  onSaveExp: (n: number) => void;
  onToggleVisited: () => void;
  onRemove: () => void;
  onReview: () => void;
}) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: item.id });
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.6 : 1,
  };
  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`flex gap-2 overflow-hidden rounded-2xl border bg-card-grad ${item.visited ? "border-primary/40" : "border-border/60"}`}
    >
      <button
        {...attributes}
        {...listeners}
        className="flex w-7 cursor-grab items-center justify-center text-muted-foreground hover:bg-secondary active:cursor-grabbing"
        aria-label="Drag to reorder"
      >
        <GripVertical className="h-4 w-4" />
      </button>
      <div className="grid flex-1 sm:grid-cols-[140px_1fr]">
        <div className="relative aspect-[4/3] sm:aspect-auto">
          <img src={item.place.image} alt={item.place.name} loading="lazy" className="h-full w-full object-cover" />
          {item.visited && (
            <span className="absolute left-2 top-2 grid h-6 w-6 place-items-center rounded-full bg-cta text-primary-foreground">
              <Check className="h-3 w-3" />
            </span>
          )}
        </div>
        <div className="p-3">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <span>{item.place.category}</span>
          </div>
          <h4 className="mt-0.5 font-display text-lg font-semibold leading-tight">{item.place.name}</h4>
          <div className="text-xs text-muted-foreground">{item.place.address}</div>
          {item.place.review && <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">"{item.place.review}"</p>}

          <div className="mt-2 flex flex-wrap items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              {isEditingExp ? (
                <form
                  className="flex items-center gap-1"
                  onSubmit={e => {
                    e.preventDefault();
                    const n = Number(expDraft);
                    if (!Number.isFinite(n) || n < 0) return;
                    onSaveExp(n);
                  }}
                >
                  <span>$</span>
                  <Input className="h-8 w-20" value={expDraft} onChange={e => setExpDraft(e.target.value)} autoFocus />
                  <Button size="sm" type="submit" variant="secondary">Save</Button>
                </form>
              ) : (
                <button
                  onClick={onStartEdit}
                  className="inline-flex items-center gap-1 rounded-full bg-secondary px-2.5 py-1 text-sm hover:bg-secondary/70"
                >
                  <span className="font-display font-semibold">${item.place.expense}</span>
                  <Pencil className="h-3 w-3 text-muted-foreground" />
                </button>
              )}
              <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
                <Star className="h-3 w-3 fill-accent text-accent" /> {item.place.rating.toFixed(1)}
              </span>
            </div>
            <div className="flex items-center gap-1">
              {item.visited ? (
                <Button size="sm" variant="ghost" onClick={onToggleVisited}>Mark unvisited</Button>
              ) : (
                <Button size="sm" className="bg-primary text-primary-foreground" onClick={onReview}>
                  <Sparkles className="mr-1.5 h-3.5 w-3.5" /> Visited & review
                </Button>
              )}
              <Button size="icon" variant="ghost" onClick={onRemove} aria-label="Remove">
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const Summary = ({ label, value }: { label: string; value: React.ReactNode }) => (
  <div>
    <div className="text-xs uppercase tracking-wider text-muted-foreground">{label}</div>
    <div className="mt-1 font-display text-2xl font-semibold">{value}</div>
  </div>
);

export default TripDetail;
