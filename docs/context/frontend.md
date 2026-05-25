# ctx:frontend | 2026-05-24 | Vite SPA
## state
todo: -
planned: (1)axios API client+JWT refresh interceptor | (2)Auth flow(Google+GitHub OAuth,callback,access-token in-memory,remove mock auth) | (3)Remove mock data+wire Dashboard(GET /trips+GET /users/me/stats→recharts) | (4)Wire Trips+TripDetail(CRUD,dnd-kit→PATCH /trips/:id/items/reorder) | (5)Wire CityDetail(GET /search/places)+PostDetail(GET /trips/:id) | (6)Collaborators panel on TripDetail(list+invite+role-change+remove) | (7)Notifications bell(poll 30s+WebSocket+mark-read) | (8)Settings page /settings(GET+PATCH /users/me) | (9)Place autocomplete(debounced GET /search/places in item form) | (10)Trip search bar in Dashboard(debounced GET /search/trips) | (11)Dark/light mode audit(next-themes+dark: classes+theme toggle in Header)
done: Lovable scaffold+routing+shadcn setup@2026-05-24 | context-doc@2026-05-24
errors: mock data still in use; AppContext has hardcoded trips; Header notifications unimplemented; no Settings page; window.location in api.ts refresh-failure handler (replace with navigate)
_update when: new pages/routes/API hooks/state management/dependencies added_

## stack
React18+TypeScript(strict) | Vite5 | bun | react-router-dom v6 | @tanstack/react-query v5 | react-hook-form+zod | shadcn/ui(Radix) | Tailwind+next-themes | axios | date-fns+react-day-picker | recharts | @dnd-kit | sonner | vitest+@testing-library/react

## files
src/main.tsx: entry — StrictMode+RouterProvider+QueryClientProvider+ThemeProvider
src/App.tsx: root layout; all routes
src/types.ts: Place,TripPost,City interfaces
src/lib/utils.ts: cn() (clsx+tailwind-merge)
src/lib/api.ts: axios; baseURL=VITE_API_BASE_URL; Bearer attach; 401→POST /auth/refresh→retry; refresh-fail→/login
src/store/AppContext.tsx: access token in-memory; ⚠TODO replace mock trips with react-query
src/data/mock.ts: ⚠DELETE ALL USAGES — hardcoded cities+trip posts
src/components/Header.tsx: nav,avatar,logout,notification bell(⚠TODO wire)
src/components/RequireAuth.tsx: route guard; redirect /login if no token
src/components/ui/*: 50+ shadcn/ui components (Tailwind+Radix only)
src/pages/Index.tsx: redirect→/login | src/pages/Login.tsx: Google+GitHub buttons→GET /auth/{provider}/login
src/pages/Signup.tsx: OAuth callback; stores token; redirect→/dashboard
src/pages/Dashboard.tsx: ⚠mock→GET /trips+GET /users/me/stats (recharts charts)
src/pages/Trips.tsx: ⚠mock→GET /trips via react-query
src/pages/TripDetail.tsx: ⚠wire all CRUD; dnd-kit reorder→PATCH /trips/:id/items/reorder
src/pages/CityDetail.tsx: ⚠wire→GET /search/places | src/pages/PostDetail.tsx: ⚠wire→GET /trips/:id
src/pages/NotFound.tsx: 404 | src/test/setup.ts: vitest+RTL global setup

## routes
/ → Index(→/login) | /login → Login | /signup → Signup
/dashboard → Dashboard[RequireAuth] | /trips → Trips[RequireAuth]
/trips/:id → TripDetail[RequireAuth] | /cities/:id → CityDetail[RequireAuth]
/posts/:id → PostDetail[RequireAuth] | * → NotFound

## env: VITE_API_BASE_URL (https://wanderplan-api-gateway.onrender.com in prod)
## rq-keys: ['trips'] ['trips',id] ['notifications'] ['user','me'] ['user','stats']
## patterns: useQuery+useMutation; invalidateQueries after mutations; notifications poll refetchInterval:30000; WS to /ws/notifications (native WebSocket)

