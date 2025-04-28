import { APIURL } from '$lib/constant';

export const load = async () => {
	const res = await fetch(`${APIURL}/teachers`);
	const data = await res.json();

	return {
		data
	};
};
