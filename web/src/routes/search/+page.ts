import { api, orErrorPage } from '#lib/api.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url }) => {
	const q = url.searchParams.get('q') ?? '';
	if (!q.trim()) return { q, search: null };
	return { q, search: await orErrorPage(api.search(fetch, q, 100)) };
};
