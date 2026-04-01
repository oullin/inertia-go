import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  root: ".",
  build: {
    outDir: "../api/cmd/dist",
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      input: "resources/js/app.js",
      output: {
        entryFileNames: "dist/app.js",
        chunkFileNames: "dist/[name]-[hash].js",
        assetFileNames: "dist/app.[ext]",
      },
    },
  },
});
