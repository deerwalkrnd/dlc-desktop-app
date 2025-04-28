import { APIURL } from '$lib/constant';
import { error } from '@sveltejs/kit';
export const load = async ({ params }: { params: any }) => {
	const classNumber = parseInt(params.slug, 10);
	const res = await fetch(`${APIURL}/classes`);
	const data = await res.json();
	let type = 'old';

	let getClass = data.classes.find((singleClass: any) => singleClass.Number === classNumber);
	if (getClass) {
		const subjects = await fetch(`${APIURL}/classes/${getClass.ID}/subjects?/type=${type}`);
		const data = await subjects.json();
		console.log(data);
		return {
			subjects: data,
			getClass
		};
	}
	error(404, 'Not found');
};
