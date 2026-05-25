import { Navigate } from "react-router-dom";
import { useApp } from "@/store/AppContext";
import { ReactNode } from "react";

export const RequireAuth = ({ children }: { children: ReactNode }) => {
  const { user } = useApp();
  if (!user) return <Navigate to="/" replace />;
  return <>{children}</>;
};
