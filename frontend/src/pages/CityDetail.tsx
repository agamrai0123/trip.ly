import { Link, useParams } from "react-router-dom";
import { cities, tripPosts } from "@/data/mock";
import { useApp } from "@/store/AppContext";
import { useMemo, useState } from "react";
import { ArrowLeft, ArrowUp, Calendar, DollarSign, MapPin, MessageCircle, Sun } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

type SortKey = "likes" | "budget" | "days";

const CityDetail = () => {
  const { id } = useParams<{ id: string }>();
  const { likes, toggleLike, userPosts, commentsByPost } = useApp();
  const [sort, setSort] = useState<SortKey>("likes");

  const city = cities.find(c => c.id === id);
  const allPosts = useMemo(() => {
    const list = [...tripPosts.filter(p => p.cityId === id), ...userPosts.filter(p => p.cityId === id)];
    return list.map(p => ({ ...p, likes: likes[p.id]?.count ?? p.likes }));
  }, [id, userPosts, likes]);

  const sorted = useMemo(() => {
    const arr = [...allPosts];
    if (sort === "likes") arr.sort((a, b) => b.likes - a.likes);
    if (sort === "budget") arr.sort((a, b) => a.totalBudget - b.totalBudget);
    if (sort === "days") arr.sort((a, b) => a.days - b.days);
    return arr;
  }, [allPosts, sort]);

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
            <div className="font-display text-xl font-semibold">{sorted.length} traveler itineraries</div>
            <div className="text-xs text-muted-foreground">Click a post to open the full itinerary</div>
          </div>
          <div className="flex items-center gap-2 text-sm">
            <span className="text-muted-foreground">Sort by</span>
            {(["likes", "budget", "days"] as SortKey[]).map(k => (
              <button
                key={k}
                onClick={() => setSort(k)}
                className={`rounded-full px-3 py-1.5 text-xs font-medium capitalize transition ${
                  sort === k ? "bg-primary text-primary-foreground" : "bg-secondary text-muted-foreground hover:text-foreground"
                }`}
              >
                {k === "likes" ? "Most loved" : k === "budget" ? "Lowest budget" : "Shortest"}
              </button>
            ))}
          </div>
        </div>

        <div className="space-y-3">
          {sorted.map((post, idx) => {
            const liked = likes[post.id]?.liked ?? false;
            const commentCount = post.comments.length + (commentsByPost[post.id]?.length ?? 0);
            return (
              <article key={post.id} className="group flex gap-3 overflow-hidden rounded-2xl border border-border/60 bg-card-grad p-3 shadow-card transition hover:border-primary/40 hover:shadow-glow">
                {/* Vote column */}
                <div className="flex w-10 flex-none flex-col items-center gap-1 pt-1">
                  <button
                    onClick={(e) => { e.preventDefault(); toggleLike(post.id, tripPosts.find(p => p.id === post.id)?.likes ?? post.likes); }}
                    className={`grid h-8 w-8 place-items-center rounded-lg transition ${liked ? "bg-primary/15 text-primary" : "text-muted-foreground hover:bg-secondary hover:text-foreground"}`}
                    aria-label="Upvote"
                  >
                    <ArrowUp className={`h-5 w-5 ${liked ? "fill-current" : ""}`} />
                  </button>
                  <span className={`text-xs font-semibold ${liked ? "text-primary" : ""}`}>{post.likes}</span>
                </div>

                {/* Content */}
                <Link to={`/post/${post.id}`} className="flex flex-1 gap-3 min-w-0">
                  <img src={post.cover} alt={post.title} loading="lazy" className="hidden sm:block h-28 w-36 flex-none rounded-xl object-cover" />
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span className="rounded-full bg-secondary px-2 py-0.5">#{idx + 1}</span>
                      <Avatar className="h-5 w-5">
                        <AvatarImage src={post.author.avatar} alt={post.author.name} />
                        <AvatarFallback>{post.author.name[0]}</AvatarFallback>
                      </Avatar>
                      <span>posted by {post.author.name}</span>
                    </div>
                    <h3 className="mt-1.5 line-clamp-2 font-display text-xl font-semibold leading-snug tracking-tight group-hover:text-primary">
                      {post.title}
                    </h3>
                    <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">{post.summary}</p>
                    <div className="mt-2 flex flex-wrap items-center gap-2">
                      <Badge variant="secondary" className="gap-1.5"><Calendar className="h-3 w-3" /> {post.days}d</Badge>
                      <Badge variant="secondary" className="gap-1.5"><DollarSign className="h-3 w-3" /> ${post.totalBudget.toLocaleString()}</Badge>
                      <Badge variant="secondary" className="gap-1.5"><Sun className="h-3 w-3" /> {post.bestTime}</Badge>
                      <Badge variant="secondary" className="gap-1.5"><MapPin className="h-3 w-3" /> {post.places.length} places</Badge>
                      <span className="ml-auto inline-flex items-center gap-1 text-xs text-muted-foreground">
                        <MessageCircle className="h-3.5 w-3.5" /> {commentCount}
                      </span>
                    </div>
                  </div>
                </Link>
              </article>
            );
          })}

          {sorted.length === 0 && (
            <div className="rounded-2xl border border-dashed border-border p-10 text-center text-muted-foreground">
              No itineraries here yet. Be the first to share one.
            </div>
          )}
        </div>
      </section>
    </main>
  );
};

export default CityDetail;
