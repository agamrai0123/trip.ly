---
applyTo: "frontend/src/**/*.test.tsx"
---

# Frontend Test Rules

## Framework
- Use `vitest` as the test runner. Use `@testing-library/react` for rendering and interaction.
- Use `@testing-library/user-event` for all user interactions (typing, clicking). Never use `fireEvent` directly.
- Use `vi.mock()` to mock the axios API client (`src/lib/apiClient.ts`) — never make real HTTP calls in tests.

## What to test
- All form validation: submit with missing/invalid fields and assert the inline error messages appear.
- All `@tanstack/react-query` hooks: wrap the component in a fresh `QueryClientProvider` for each test.
- All route guards: assert that unauthenticated users are redirected to `/login`.
- Drag-and-drop reorder: simulate `@dnd-kit` events and assert the mutation is called with the correct new order.

## Assertions
- Query by accessible role first (`getByRole`, `getByLabelText`). Fall back to `getByTestId` only when semantics are insufficient.
- Assert that loading spinners appear while queries are pending and disappear on resolution.
- Never snapshot-test entire page layouts — snapshot only small, stable UI atoms.

## Naming
- Test files: `<Component>.test.tsx` next to the component file.
- Describe blocks: component name. It blocks: behaviour in plain English (e.g. `it("shows an error when title is empty")`).
