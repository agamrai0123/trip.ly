---
applyTo: "frontend/src/**/*.tsx"
---

# Frontend Component Rules

## Components
- Every component is a named function export — never default-export anonymous arrow functions.
- Use shadcn/ui (Radix UI) primitives for all UI elements. Do not install or import other component libraries.
- All Tailwind class strings must be static where possible so PurgeCSS can detect them. Use `cn()` from `lib/utils` to merge conditional classes.
- Dark/light mode: always pair a base class with a `dark:` variant (e.g. `bg-white dark:bg-zinc-900`). Never hardcode colours outside Tailwind tokens.

## Data fetching
- All server state lives in `@tanstack/react-query` v5 queries and mutations. No `useEffect` + `fetch`/`axios` for data loading.
- Define query keys as constants in `src/lib/queryKeys.ts`. Never inline raw string arrays as query keys.
- On mutation success, invalidate the affected query keys — do not manually update the cache unless performance requires it.
- Show loading and error states for every query. Never render `undefined` data to the DOM.

## Forms
- All forms use `react-hook-form` + `zod`. Define the zod schema above the component; derive the TypeScript type with `z.infer<typeof schema>`.
- Validation errors must be shown inline next to the field using shadcn/ui `FormMessage`. Never use `alert()`.
- Disable the submit button while the mutation is `isPending`.

## Date handling
- Use `date-fns` for all date arithmetic and formatting. Never import `moment` or `dayjs`.
- Use `react-day-picker` for all date picker UI. Never write a custom calendar component.

## Drag-and-drop
- All reorder interactions use `@dnd-kit`. On every `onDragEnd` event, immediately call the `PATCH /trips/:id/items/reorder` mutation — do not wait for a save button.
- Maintain optimistic order in local state so the UI does not jump while the mutation is in-flight.

## API calls
- All HTTP calls go through the axios instance exported from `src/lib/apiClient.ts`. Never import axios directly.
- The base URL always comes from `import.meta.env.VITE_API_BASE_URL`. Never hardcode `localhost` or any port.
- The axios instance has a response interceptor that retries on 401 by calling the refresh endpoint once, then redirecting to `/login` on a second 401.

## Routing
- Use `react-router-dom` v6 `<Link>` and `useNavigate`. Never use `<a href>` for internal navigation or `window.location.href` for redirects.
- Protect authenticated routes with the existing `<RequireAuth>` wrapper. Do not duplicate auth checks in components.

## Environment variables
- Prefix all custom env vars with `VITE_`. Never access `process.env` — use `import.meta.env`.
