import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchMe, updateMe, getAccessToken } from "@/lib/api";
import { useApp } from "@/store/AppContext";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useState, useEffect } from "react";
import { User, Mail, Shield } from "lucide-react";

const Profile = () => {
  const qc = useQueryClient();
  const { login } = useApp();

  const { data: user, isLoading } = useQuery({
    queryKey: ["user", "me"],
    queryFn: fetchMe,
  });

  const [name, setName] = useState("");
  const [avatarUrl, setAvatarUrl] = useState("");

  useEffect(() => {
    if (user) {
      setName(user.name ?? "");
      setAvatarUrl(user.avatar_url ?? "");
    }
  }, [user]);

  const save = useMutation({
    mutationFn: () => updateMe({ name: name.trim(), avatar_url: avatarUrl.trim() }),
    onSuccess: (updated) => {
      qc.setQueryData(["user", "me"], updated);
      const token = getAccessToken();
      if (token) login(token, updated);
      toast.success("Profile updated");
    },
    onError: () => toast.error("Could not save changes"),
  });

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }

  if (!user) {
    return (
      <main className="container py-20 text-center text-muted-foreground">
        Could not load profile.
      </main>
    );
  }

  const initials = user.name
    ? user.name.split(" ").map((n) => n[0]).join("").toUpperCase().slice(0, 2)
    : user.email[0].toUpperCase();

  return (
    <main className="container py-10 md:py-14 max-w-2xl">
      <h1 className="font-display text-4xl font-semibold tracking-tight">Profile</h1>
      <p className="mt-2 text-muted-foreground">Manage your account details.</p>

      {/* Avatar preview */}
      <div className="mt-8 flex items-center gap-5">
        <Avatar className="h-20 w-20 ring-4 ring-primary/20">
          <AvatarImage src={avatarUrl || user.avatar_url} alt={user.name} />
          <AvatarFallback className="text-2xl">{initials}</AvatarFallback>
        </Avatar>
        <div>
          <p className="font-display text-xl font-semibold">{user.name}</p>
          <p className="text-sm text-muted-foreground">{user.email}</p>
          <span className="mt-1 inline-flex items-center gap-1 rounded-full border border-border/60 px-2 py-0.5 text-xs text-muted-foreground capitalize">
            <Shield className="h-3 w-3" /> {user.provider}
          </span>
        </div>
      </div>

      {/* Edit form */}
      <form
        onSubmit={(e) => { e.preventDefault(); save.mutate(); }}
        className="mt-8 space-y-5 rounded-2xl border border-border/60 bg-card-grad p-6"
      >
        <div>
          <Label htmlFor="p-name" className="flex items-center gap-1.5 text-sm">
            <User className="h-4 w-4" /> Display name
          </Label>
          <Input
            id="p-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Your name"
            className="mt-1"
          />
        </div>

        <div>
          <Label htmlFor="p-email" className="flex items-center gap-1.5 text-sm">
            <Mail className="h-4 w-4" /> Email
          </Label>
          <Input
            id="p-email"
            value={user.email}
            disabled
            readOnly
            className="mt-1 cursor-not-allowed opacity-60"
          />
          <p className="mt-1 text-xs text-muted-foreground">Email is managed by your OAuth provider and cannot be changed here.</p>
        </div>

        <div>
          <Label htmlFor="p-avatar" className="text-sm">Avatar URL</Label>
          <Input
            id="p-avatar"
            type="url"
            value={avatarUrl}
            onChange={(e) => setAvatarUrl(e.target.value)}
            placeholder="https://…"
            className="mt-1"
          />
        </div>

        <div className="flex justify-end">
          <Button type="submit" disabled={save.isPending} className="bg-cta text-primary-foreground shadow-glow">
            {save.isPending ? "Saving…" : "Save changes"}
          </Button>
        </div>
      </form>
    </main>
  );
};

export default Profile;
