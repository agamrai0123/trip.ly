import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react-swc";
import path from "path";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test/setup.ts"],
    include: ["src/**/*.{test,spec}.{ts,tsx}"],
    // forks pool: each file gets its own child process (no shared jsdom state)
    // avoids the ENOSPC issue caused by the worker_threads pool writing to C: temp
    pool: "forks",
  },
  cacheDir: "./node_modules/.vitest-cache",
  resolve: {
    alias: { "@": path.resolve(__dirname, "./src") },
  },
});
