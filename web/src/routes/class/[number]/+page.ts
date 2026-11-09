import { api, orErrorPage } from '#lib/api.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	return { cls: await orErrorPage(api.class(fetch, params.number)) };
};
