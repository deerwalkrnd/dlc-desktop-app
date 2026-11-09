import { api, orErrorPage } from '#lib/api.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const [classes, library] = await orErrorPage(
		Promise.all([api.classes(fetch), api.library(fetch)])
	);
	return { classes, library };
};
