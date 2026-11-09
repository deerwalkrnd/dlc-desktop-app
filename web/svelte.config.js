import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		// Output goes to build/, which the Go binary embeds (see embed.go).
		// Every unknown path falls back to index.html (SPA mode).
		adapter: adapter({ pages: 'build', assets: 'build', fallback: 'index.html' })
	}
};

export default config;
