import { Link } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { Briefcase, Calendar, MapPin, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cities } from "@/data/mock";

const Trips = () => {
  const { trips } = useApp();

  return (
    <main className="container py-10 md:py-14">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="font-display text-4xl font-semibold tracking-tight md:text-5xl">My trips</h1>
          <p className="mt-2 text-muted-foreground">Itineraries you've borrowed and made yours.</p>
        </div>
        <Link to="/dashboard">
          <Button variant="secondary"><Plus className="mr-1.5 h-4 w-4" /> Plan another</Button>
        </Link>
      </div>

      {trips.length === 0 ? (
        <div className="mt-10 grid place-items-center rounded-3xl border border-dashed border-border/80 bg-card-grad p-16 text-center">
          <div className="grid h-14 w-14 place-items-center rounded-2xl bg-primary/15 text-primary">
            <Briefcase className="h-6 w-6" />
          </div>
          <h2 className="mt-5 font-display text-2xl">No trips yet</h2>
          <p className="mt-2 max-w-sm text-muted-foreground">
            Browse a city and add a traveler's itinerary to get started.
          </p>
          <Link to="/dashboard"><Button className="mt-6 bg-cta text-primary-foreground shadow-glow">Discover cities</Button></Link>
        </div>
      ) : (
        <div className="mt-10 grid grid-cols-1 gap-5 md:grid-cols-2">
          {trips.map(t => {
            const city = cities.find(c => c.id === t.cityId);
            const total = t.items.reduce((s, i) => s + i.place.expense, 0);
            const visited = t.items.filter(i => i.visited).length;
            return (
              <Link key={t.id} to={`/trips/${t.id}`} className="group overflow-hidden rounded-3xl border border-border/60 bg-card shadow-card transition hover:-translate-y-0.5 hover:shadow-glow">
                <div className="relative aspect-[16/9]">
                  {city && <img src={city.image} alt={city.name} className="h-full w-full object-cover transition-transform duration-700 group-hover:scale-105" />}
                  <div className="absolute inset-0 bg-gradient-to-t from-background via-background/40 to-transparent" />
                  <div className="absolute inset-x-0 bottom-0 p-5">
                    <div className="flex items-center gap-1.5 text-xs text-primary">
                      <MapPin className="h-3.5 w-3.5" />
                      <span className="uppercase tracking-wider">{city?.country}</span>
                    </div>
                    <h3 className="mt-1 font-display text-2xl font-semibold">{city?.name ?? t.cityName}</h3>
                  </div>
                </div>
                <div className="grid grid-cols-3 divide-x divide-border/60 border-t border-border/60 text-center">
                  <Stat icon={<Calendar className="h-4 w-4" />} label="Days" value={t.days} />
                  <Stat icon={<span className="font-display">$</span>} label="Budget" value={total.toLocaleString()} />
                  <Stat label="Visited" value={`${visited}/${t.items.length}`} />
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </main>
  );
};

const Stat = ({ icon, label, value }: { icon?: React.ReactNode; label: string; value: React.ReactNode }) => (
  <div className="px-3 py-4">
    <div className="inline-flex items-center gap-1.5 text-xs uppercase tracking-wider text-muted-foreground">
      {icon}{label}
    </div>
    <div className="mt-0.5 font-display text-xl font-semibold">{value}</div>
  </div>
);

export default Trips;
