import { useTheme } from "next-themes";
import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useApp } from "@/store/AppContext";

const Settings = () => {
  const { user } = useApp();
  const { theme, setTheme } = useTheme();

  return (
    <main className="container py-10 md:py-14">
      <h1 className="font-display text-4xl font-semibold tracking-tight md:text-5xl">Settings</h1>
      <p className="mt-2 text-muted-foreground">Manage your preferences and account details.</p>

      <div className="mt-10 max-w-lg space-y-6">
        {/* Appearance */}
        <section className="rounded-2xl border border-border/60 bg-card-grad p-6">
          <h2 className="font-display text-xl font-semibold">Appearance</h2>
          <p className="mt-1 text-sm text-muted-foreground">Choose your preferred colour theme.</p>
          <div className="mt-4 flex gap-3">
            <Button
              variant={theme === "dark" ? "default" : "outline"}
              onClick={() => setTheme("dark")}
              className="flex-1 gap-2"
            >
              <Moon className="h-4 w-4" />
              Dark
            </Button>
            <Button
              variant={theme === "light" ? "default" : "outline"}
              onClick={() => setTheme("light")}
              className="flex-1 gap-2"
            >
              <Sun className="h-4 w-4" />
              Light
            </Button>
            <Button
              variant={theme === "system" ? "default" : "outline"}
              onClick={() => setTheme("system")}
              className="flex-1"
            >
              System
            </Button>
          </div>
        </section>

        {/* Account info */}
        {user && (
          <section className="rounded-2xl border border-border/60 bg-card-grad p-6">
            <h2 className="font-display text-xl font-semibold">Account</h2>
            <div className="mt-4 space-y-4">
              <div>
                <div className="text-xs uppercase tracking-wider text-muted-foreground">Name</div>
                <div className="mt-0.5 text-sm font-medium">{user.name}</div>
              </div>
              <div>
                <div className="text-xs uppercase tracking-wider text-muted-foreground">Email</div>
                <div className="mt-0.5 text-sm">{user.email}</div>
              </div>
              <div>
                <div className="text-xs uppercase tracking-wider text-muted-foreground">Sign-in provider</div>
                <div className="mt-0.5 text-sm capitalize">{user.provider}</div>
              </div>
            </div>
          </section>
        )}
      </div>
    </main>
  );
};

export default Settings;
