import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { Toaster } from "@/components/ui/toaster";
import { TooltipProvider } from "@/components/ui/tooltip";
import { AppProvider } from "@/store/AppContext";
import { Header } from "@/components/Header";
import { RequireAuth } from "@/components/RequireAuth";
import Index from "./pages/Index.tsx";
import Signup from "./pages/Signup.tsx";
import NotFound from "./pages/NotFound.tsx";
import Dashboard from "./pages/Dashboard.tsx";
import CityDetail from "./pages/CityDetail.tsx";
import Trips from "./pages/Trips.tsx";
import TripDetail from "./pages/TripDetail.tsx";
import PostDetail from "./pages/PostDetail.tsx";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <AppProvider>
          <Header />
          <Routes>
            <Route path="/" element={<Index />} />
            <Route path="/signup" element={<Signup />} />
            <Route path="/dashboard" element={<RequireAuth><Dashboard /></RequireAuth>} />
            <Route path="/city/:id" element={<RequireAuth><CityDetail /></RequireAuth>} />
            <Route path="/post/:id" element={<RequireAuth><PostDetail /></RequireAuth>} />
            <Route path="/trips" element={<RequireAuth><Trips /></RequireAuth>} />
            <Route path="/trips/:id" element={<RequireAuth><TripDetail /></RequireAuth>} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </AppProvider>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
