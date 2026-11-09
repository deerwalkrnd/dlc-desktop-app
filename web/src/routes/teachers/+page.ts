import { api, orErrorPage } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return { teachers: await orErrorPage(api.teachers(fetch)) };
};
