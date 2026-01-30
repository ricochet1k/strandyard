import { defineConfig } from "vite";
import solid from "vite-plugin-solid";

export default defineConfig({
  plugins: [solid()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8686',
        changeOrigin: true,
      },
    },
  },
  build: {
    target: "esnext",
    outDir: "../../pkg/web/dist",
    emptyOutDir: true,
  },
});
