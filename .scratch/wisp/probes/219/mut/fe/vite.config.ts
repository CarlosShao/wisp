import { fileURLToPath, URL } from "node:url";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { writeFileSync } from "node:fs";
import { defineConfig } from "vite";

// The panel is loaded by WebView2 from embedded bytes, not over HTTP
// (D29: AddWebResourceRequestedFilter, no local server), so the bundle must
// reference its own files relatively. A root-absolute "/assets/..." URL would
// resolve against the virtual host's origin and 404.
// vite empties outDir before writing, which would delete the tracked anchor that
// keeps the go:embed pattern in frontend/embed.go valid in an unbuilt checkout.
function keepDistAnchor() {
  return {
    name: "wisp-keep-dist-anchor",
    closeBundle() {
      writeFileSync(new URL("./dist/.gitkeep", import.meta.url), "");
    },
  };
}

export default defineConfig({
  base: "./",
  plugins: [react(), tailwindcss(), keepDistAnchor()],
  resolve: {
    alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
  },
  build: {
    // AC#1: the artifact is embedded into wisp.exe, so a size jump is a
    // regression the reviewer should see. No CDN is configured anywhere and
    // none may be added: the binary has to work offline.
    reportCompressedSize: true,
    sourcemap: false,
  },
});
