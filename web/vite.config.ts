import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// During `npm run dev`, API and video requests go to the Go server
// (`go run ./cmd/dlc-desktop-app -open=false`). changeOrigin makes the Host
// header match what the server expects.
const backend = process.env.DLC_BACKEND ?? 'http://127.0.0.1:3000';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': { target: backend, changeOrigin: true },
			'/media': { target: backend, changeOrigin: true }
		}
	}
});
