import { Navigate } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { ReactNode } from "react";

export const RequireAuth = ({ children }: { children: ReactNode }) => {
  const { user, loading } = useApp();
  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }
  if (!user) return <Navigate to="/" replace />;
  return <>{children}</>;
};
