import { api, orErrorPage } from '#lib/api.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	const [cls, subject] = await orErrorPage(
		Promise.all([api.class(fetch, params.number), api.subject(fetch, params.number, params.slug)])
	);
	return { cls, subject };
};
