import { Link, useNavigate, useParams } from "react-router-dom";
import { useMemo, useState } from "react";
import { tripPosts, cities } from "@/data/mock";
import { useApp } from "@/store/AppContext";
import { ArrowLeft, Calendar, DollarSign, Heart, MapPin, MessageCircle, Plus, Star, Sun } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuLabel,
  ContextMenuSeparator,
  ContextMenuSub,
  ContextMenuSubContent,
  ContextMenuSubTrigger,
  ContextMenuTrigger,
} from "@/components/ui/context-menu";
import { toast } from "sonner";
import type { Place } from "@/types";

const PostDetail = () => {
  const { id } = useParams<{ id: string }>();
  const nav = useNavigate();
  const { likes, toggleLike, addComment, commentsByPost, user, userPosts, trips, addPlaceToTrip, addPostToTrips, createEmptyTrip } = useApp();
  const [comment, setComment] = useState("");

  const post = useMemo(() => [...tripPosts, ...userPosts].find(p => p.id === id), [id, userPosts]);
  const city = post ? cities.find(c => c.id === post.cityId) : undefined;

  if (!post) {
    return (
      <main className="container py-20 text-center">
        <p className="text-muted-foreground">Post not found.</p>
        <Link to="/dashboard" className="mt-4 inline-block text-primary underline">Back to dashboard</Link>
      </main>
    );
  }

  const liked = likes[post.id]?.liked ?? false;
  const likeCount = likes[post.id]?.count ?? post.likes;
  const allComments = [...post.comments, ...(commentsByPost[post.id] ?? [])];
  const cityTrips = trips.filter(t => t.cityId === post.cityId);

  const handleAddPlace = (place: Place, tripId: string) => {
    addPlaceToTrip(tripId, place);
    toast.success(`Added “${place.name}” to your trip`);
  };
  const handleAddPlaceToNewTrip = (place: Place) => {
    const tid = createEmptyTrip(post.cityId, post.days);
    addPlaceToTrip(tid, place);
    toast.success(`New trip created with “${place.name}”`);
  };
  const handleAddWholeItinerary = () => {
    const tid = addPostToTrips(post);
    toast.success("Whole itinerary added to your trips");
    nav(`/trips/${tid}`);
  };

  return (
    <main>
      <section className="relative h-[40vh] min-h-[300px] overflow-hidden">
        <img src={post.cover} alt={post.title} className="absolute inset-0 h-full w-full object-cover" />
        <div className="absolute inset-0 bg-gradient-to-t from-background via-background/70 to-background/20" />
        <div className="container relative flex h-full flex-col justify-end pb-8">
          <Link to={`/city/${post.cityId}`} className="mb-3 inline-flex w-fit items-center gap-2 rounded-full bg-background/40 px-3 py-1.5 text-sm backdrop-blur hover:bg-background/60">
            <ArrowLeft className="h-4 w-4" /> {city?.name ?? "Back"}
          </Link>
          <div className="flex items-center gap-2 text-primary">
            <MapPin className="h-4 w-4" />
            <span className="text-sm uppercase tracking-wider">{city?.country}</span>
          </div>
          <h1 className="mt-1 font-display text-4xl font-semibold leading-tight tracking-tight md:text-5xl">{post.title}</h1>
        </div>
      </section>

      <section className="container py-8">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="flex items-center gap-3">
            <Avatar className="h-10 w-10">
              <AvatarImage src={post.author.avatar} alt={post.author.name} />
              <AvatarFallback>{post.author.name[0]}</AvatarFallback>
            </Avatar>
            <div>
              <div className="text-sm font-medium">{post.author.name}</div>
              <div className="text-xs text-muted-foreground">{post.places.length} places · {allComments.length} comments</div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => toggleLike(post.id, post.likes)}
              className={`inline-flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm transition ${
                liked ? "border-primary bg-primary/15 text-primary" : "border-border text-muted-foreground hover:text-foreground"
              }`}
            >
              <Heart className={`h-4 w-4 ${liked ? "fill-current" : ""}`} /> {likeCount}
            </button>
            <Button onClick={handleAddWholeItinerary} className="bg-cta text-primary-foreground shadow-glow">
              <Plus className="mr-1.5 h-4 w-4" /> Add whole trip
            </Button>
          </div>
        </div>

        <p className="mt-4 max-w-3xl text-muted-foreground">{post.summary}</p>

        <div className="mt-4 flex flex-wrap gap-2">
          <Badge variant="secondary" className="gap-1.5"><Calendar className="h-3.5 w-3.5" /> {post.days} days</Badge>
          <Badge variant="secondary" className="gap-1.5"><DollarSign className="h-3.5 w-3.5" /> ${post.totalBudget.toLocaleString()} total</Badge>
          <Badge variant="secondary" className="gap-1.5"><Sun className="h-3.5 w-3.5" /> Best: {post.bestTime}</Badge>
        </div>

        {/* Itinerary places */}
        <div className="mt-8">
          <div className="mb-3 flex items-end justify-between">
            <h2 className="font-display text-2xl font-semibold">The itinerary</h2>
            <span className="text-xs text-muted-foreground">Right-click any place to add to a trip</span>
          </div>
          <div className="space-y-3">
            {post.places.map(p => (
              <ContextMenu key={p.id}>
                <ContextMenuTrigger asChild>
                  <div className="flex cursor-context-menu gap-3 rounded-2xl border border-border/60 bg-card-grad p-3 transition hover:border-primary/40">
                    <img src={p.image} alt={p.name} loading="lazy" className="h-24 w-28 flex-none rounded-xl object-cover" />
                    <div className="min-w-0 flex-1">
                      <div className="flex items-start justify-between gap-2">
                        <div>
                          <div className="flex items-center gap-2">
                            <span className="text-xs uppercase tracking-wider text-accent">{p.category}</span>
                            <span className="text-xs text-muted-foreground">· {p.slot}</span>
                          </div>
                          <div className="mt-0.5 font-medium">{p.name}</div>
                          <div className="text-xs text-muted-foreground">{p.address}</div>
                        </div>
                        <div className="text-right">
                          <div className="font-display text-lg font-semibold">${p.expense}</div>
                          <div className="inline-flex items-center gap-0.5 text-xs text-muted-foreground">
                            <Star className="h-3 w-3 fill-accent text-accent" /> {p.rating.toFixed(1)}
                          </div>
                        </div>
                      </div>
                      <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">"{p.review}"</p>
                      <div className="mt-1 text-xs text-muted-foreground">
                        📍 {p.lat.toFixed(3)}, {p.lng.toFixed(3)}{p.bestTime ? ` · best: ${p.bestTime}` : ""}
                      </div>
                    </div>
                  </div>
                </ContextMenuTrigger>
                <ContextMenuContent className="w-60">
                  <ContextMenuLabel>Add “{p.name}” to…</ContextMenuLabel>
                  <ContextMenuSeparator />
                  {cityTrips.length > 0 ? (
                    <ContextMenuSub>
                      <ContextMenuSubTrigger>Existing trip</ContextMenuSubTrigger>
                      <ContextMenuSubContent>
                        {cityTrips.map(t => (
                          <ContextMenuItem key={t.id} onSelect={() => handleAddPlace(p, t.id)}>
                            {t.cityName} · {t.days}d · {t.items.length} stops
                          </ContextMenuItem>
                        ))}
                      </ContextMenuSubContent>
                    </ContextMenuSub>
                  ) : (
                    <ContextMenuItem disabled>No existing trips for {city?.name}</ContextMenuItem>
                  )}
                  <ContextMenuItem onSelect={() => handleAddPlaceToNewTrip(p)}>+ New trip</ContextMenuItem>
                </ContextMenuContent>
              </ContextMenu>
            ))}
          </div>
        </div>

        {/* Comments */}
        <div className="mt-10 border-t border-border/60 pt-6">
          <div className="mb-3 inline-flex items-center gap-1.5 text-sm font-medium">
            <MessageCircle className="h-4 w-4" /> Comments ({allComments.length})
          </div>
          <div className="space-y-3">
            {allComments.map(c => (
              <div key={c.id} className="flex gap-2.5">
                <Avatar className="h-7 w-7"><AvatarImage src={c.avatar} /><AvatarFallback>{c.author[0]}</AvatarFallback></Avatar>
                <div className="flex-1 rounded-xl bg-background/50 px-3 py-2 text-sm">
                  <div className="text-xs"><span className="font-medium">{c.author}</span> <span className="text-muted-foreground">· {c.at}</span></div>
                  <div>{c.text}</div>
                </div>
              </div>
            ))}
          </div>
          <form
            className="mt-3 flex gap-2"
            onSubmit={(e) => {
              e.preventDefault();
              if (!comment.trim() || !user) return;
              addComment(post.id, comment.trim());
              setComment("");
            }}
          >
            <Input value={comment} onChange={e => setComment(e.target.value)} placeholder="Add a comment…" className="h-10 bg-background/40" />
            <Button type="submit" size="sm" variant="secondary">Post</Button>
          </form>
        </div>
      </section>
    </main>
  );
};

export default PostDetail;
