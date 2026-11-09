import { api, orErrorPage } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return { library: await orErrorPage(api.library(fetch)) };
};
