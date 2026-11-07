// vite config for the Wisp jarvis-orb single-file export (ticket 65).
// No @vitejs/plugin-react: the entry uses no JSX runtime beyond what esbuild
// already transpiles, and one fewer dependency keeps the offline build path
// honest. The single-file shape (everything inlined into index.html) is
// produced by inline.mjs after vite's emit, NOT by assetsInlineLimit tricks,
// so the output document is deterministic regardless of chunk boundaries.
import { defineConfig } from "vite";

export default defineConfig({
  resolve: {
    // The file: dependency is symlinked, so its `import "react"` would walk up
    // into the upstream clone's own devDep react - a SECOND React instance and
    // the classic "null (reading 'useRef')" hook-dispatcher crash. Dedupe to
    // this scaffold's single copy (the vite docs' own fix for duplicate react).
    dedupe: ["react", "react-dom"],
  },
  esbuild: {
    // Automatic JSX runtime: the entry imports createRoot, not React itself.
    jsx: "automatic",
  },
  build: {
    target: "es2020",
    minify: true,
    sourcemap: false,
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
      },
    },
  },
});
