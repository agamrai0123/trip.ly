import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Compass } from "lucide-react";
import { toast } from "sonner";
import hero from "@/assets/hero-ocean.jpg";

const Signup = () => {
  const { login, user } = useApp();
  const nav = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");

  if (user) {
    nav("/dashboard", { replace: true });
  }

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email || !password) {
      toast.error("Fill in all fields");
      return;
    }
    if (password.length < 6) {
      toast.error("Password must be at least 6 characters");
      return;
    }
    if (password !== confirm) {
      toast.error("Passwords don't match");
      return;
    }
    login(email);
    toast.success(`Welcome to Wandr, ${name.split(" ")[0]}`);
    nav("/dashboard");
  };

  const onGoogle = () => {
    login("traveler@google.com");
    toast.success("Signed up with Google");
    nav("/dashboard");
  };

  return (
    <div className="relative min-h-screen overflow-hidden">
      <div className="absolute inset-0">
        <img src={hero} alt="" width={1600} height={1024} className="h-full w-full object-cover opacity-60" />
        <div className="absolute inset-0 bg-gradient-to-t from-background via-background/85 to-background/40" />
        <div className="absolute inset-0 bg-aurora" />
      </div>

      <div className="relative z-10 grid min-h-screen lg:grid-cols-2">
        <div className="flex flex-col justify-between p-8 md:p-12 lg:p-16">
          <Link to="/" className="flex items-center gap-2">
            <div className="grid h-10 w-10 place-items-center rounded-xl bg-cta shadow-glow">
              <Compass className="h-5 w-5 text-primary-foreground" />
            </div>
            <span className="font-display text-2xl font-semibold tracking-tight">wandr</span>
          </Link>

          <div className="max-w-xl animate-fade-up">
            <h1 className="font-display text-5xl font-semibold leading-[1.05] tracking-tight text-balance md:text-6xl lg:text-7xl">
              Start your<br />
              <span className="bg-cta bg-clip-text text-transparent">next adventure.</span>
            </h1>
            <p className="mt-6 max-w-md text-lg text-muted-foreground">
              Join travelers sharing real itineraries. Plan smarter, spend better, explore further.
            </p>
          </div>

          <div className="text-xs text-muted-foreground">© {new Date().getFullYear()} Wandr · Made for explorers</div>
        </div>

        <div className="flex items-center justify-center p-6 md:p-12">
          <div className="w-full max-w-md animate-fade-up">
            <div className="glass rounded-3xl p-8 shadow-card md:p-10">
              <h2 className="font-display text-3xl font-semibold tracking-tight">Create account</h2>
              <p className="mt-2 text-sm text-muted-foreground">It takes less than a minute.</p>

              <form onSubmit={onSubmit} className="mt-8 space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Full name</Label>
                  <Input
                    id="name"
                    type="text"
                    placeholder="Ada Lovelace"
                    value={name}
                    onChange={e => setName(e.target.value)}
                    autoComplete="name"
                    maxLength={80}
                    className="h-12 bg-input/60"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="you@example.com"
                    value={email}
                    onChange={e => setEmail(e.target.value)}
                    autoComplete="email"
                    maxLength={120}
                    className="h-12 bg-input/60"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    type="password"
                    placeholder="At least 6 characters"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    autoComplete="new-password"
                    className="h-12 bg-input/60"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="confirm">Confirm password</Label>
                  <Input
                    id="confirm"
                    type="password"
                    placeholder="Repeat your password"
                    value={confirm}
                    onChange={e => setConfirm(e.target.value)}
                    autoComplete="new-password"
                    className="h-12 bg-input/60"
                  />
                </div>

                <Button type="submit" size="lg" className="h-12 w-full bg-cta text-primary-foreground hover:opacity-95 shadow-glow">
                  Create account
                </Button>
              </form>

              <div className="my-6 flex items-center gap-3">
                <div className="h-px flex-1 bg-border" />
                <span className="text-xs uppercase tracking-wider text-muted-foreground">or</span>
                <div className="h-px flex-1 bg-border" />
              </div>

              <Button
                type="button"
                variant="outline"
                size="lg"
                onClick={onGoogle}
                className="h-12 w-full border-border/80 bg-background/40 hover:bg-background/70"
              >
                <GoogleIcon className="mr-2 h-5 w-5" />
                Continue with Google
              </Button>

              <p className="mt-6 text-center text-sm text-muted-foreground">
                Already have an account?{" "}
                <Link to="/" className="font-medium text-foreground underline-offset-4 hover:underline">
                  Sign in
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const GoogleIcon = (props: React.SVGProps<SVGSVGElement>) => (
  <svg viewBox="0 0 48 48" {...props}>
    <path fill="#FFC107" d="M43.6 20.5H42V20H24v8h11.3C33.7 32.6 29.3 36 24 36c-6.6 0-12-5.4-12-12s5.4-12 12-12c3.1 0 5.9 1.2 8 3.1l5.7-5.7C34 6.1 29.3 4 24 4 12.9 4 4 12.9 4 24s8.9 20 20 20 20-8.9 20-20c0-1.2-.1-2.4-.4-3.5z"/>
    <path fill="#FF3D00" d="M6.3 14.7l6.6 4.8C14.7 16.1 19 13 24 13c3.1 0 5.9 1.2 8 3.1l5.7-5.7C34 6.1 29.3 4 24 4 16.3 4 9.7 8.3 6.3 14.7z"/>
    <path fill="#4CAF50" d="M24 44c5.2 0 9.9-2 13.4-5.2l-6.2-5.2C29.2 35 26.7 36 24 36c-5.3 0-9.7-3.4-11.3-8.1l-6.5 5C9.6 39.7 16.2 44 24 44z"/>
    <path fill="#1976D2" d="M43.6 20.5H42V20H24v8h11.3c-.8 2.3-2.3 4.2-4.1 5.6l6.2 5.2C41 35.5 44 30.2 44 24c0-1.2-.1-2.4-.4-3.5z"/>
  </svg>
);

export default Signup;
