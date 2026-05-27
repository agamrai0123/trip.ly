import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { Button } from "@/components/ui/button";
import { Compass } from "lucide-react";
import hero from "@/assets/hero-ocean.jpg";
import { googleLoginUrl, githubLoginUrl } from "@/lib/api";

const Login = () => {
  const { user, loading } = useApp();
  const nav = useNavigate();

  useEffect(() => {
    if (!loading && user) nav("/dashboard", { replace: true });
  }, [user, loading, nav]);

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="relative min-h-screen overflow-hidden">
      {/* Hero image */}
      <div className="absolute inset-0">
        <img src={hero} alt="" width={1600} height={1024} className="h-full w-full object-cover opacity-60" />
        <div className="absolute inset-0 bg-gradient-to-t from-background via-background/85 to-background/40" />
        <div className="absolute inset-0 bg-aurora" />
      </div>

      <div className="relative z-10 grid min-h-screen lg:grid-cols-2">
        {/* Left — branding */}
        <div className="flex flex-col justify-between p-8 md:p-12 lg:p-16">
          <div className="flex items-center gap-2">
            <div className="grid h-10 w-10 place-items-center rounded-xl bg-cta shadow-glow">
              <Compass className="h-5 w-5 text-primary-foreground" />
            </div>
            <span className="font-display text-2xl font-semibold tracking-tight">wandr</span>
          </div>

          <div className="max-w-xl animate-fade-up">
            <h1 className="font-display text-5xl font-semibold leading-[1.05] tracking-tight text-balance md:text-6xl lg:text-7xl">
              Travel by the<br />
              <span className="bg-cta bg-clip-text text-transparent">routes others love.</span>
            </h1>
            <p className="mt-6 max-w-md text-lg text-muted-foreground">
              Pick a destination. Borrow the best itineraries from real travelers. Make them yours.
            </p>
          </div>

          <div className="text-xs text-muted-foreground">© {new Date().getFullYear()} Wandr · Made for explorers</div>
        </div>

        {/* Right — auth card */}
        <div className="flex items-center justify-center p-6 md:p-12">
          <div className="w-full max-w-md animate-fade-up">
            <div className="glass rounded-3xl p-8 shadow-card md:p-10">
              <h2 className="font-display text-3xl font-semibold tracking-tight">Sign in</h2>
              <p className="mt-2 text-sm text-muted-foreground">Continue planning your next trip.</p>

              <div className="mt-8 space-y-3">
                <Button
                  type="button"
                  variant="outline"
                  size="lg"
                  onClick={() => { window.location.href = googleLoginUrl(); }}
                  className="h-12 w-full border-border/80 bg-background/40 hover:bg-background/70"
                >
                  <GoogleIcon className="mr-2 h-5 w-5" />
                  Continue with Google
                </Button>

                <Button
                  type="button"
                  variant="outline"
                  size="lg"
                  onClick={() => { window.location.href = githubLoginUrl(); }}
                  className="h-12 w-full border-border/80 bg-background/40 hover:bg-background/70"
                >
                  <GithubIcon className="mr-2 h-5 w-5" />
                  Continue with GitHub
                </Button>
              </div>

              <p className="mt-6 text-center text-xs text-muted-foreground">
                By continuing you agree to our Terms &amp; Privacy.
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

const GithubIcon = (props: React.SVGProps<SVGSVGElement>) => (
  <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
    <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
  </svg>
);

export default Login;
