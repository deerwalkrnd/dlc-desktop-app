import adapter from '@sveltejs/adapter-static';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// During `npm run dev`, API and video requests go to the Go server
// (`go run ./cmd/dlc-desktop-app -open=false`). changeOrigin makes the Host
// header match what the server expects.
const backend = process.env.DLC_BACKEND ?? 'http://127.0.0.1:3000';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			// Output goes to build/, which the Go binary embeds (see embed.go).
			// Every unknown path falls back to index.html (SPA mode).
			adapter: adapter({ pages: 'build', assets: 'build', fallback: 'index.html' }),
			// The app and the server always ship together in one binary, so
			// there is never a newer version to poll for.
			version: { pollInterval: 0 }
		})
	],
	server: {
		proxy: {
			'/api': { target: backend, changeOrigin: true },
			'/media': { target: backend, changeOrigin: true }
		}
	}
});
