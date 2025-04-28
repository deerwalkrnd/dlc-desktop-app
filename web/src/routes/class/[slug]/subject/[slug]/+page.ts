import { APIURL } from '$lib/constant';

export const load = async ({ params }: { params: any }) => {
	const subjectName = params.slug.toLowerCase();
	console.log(params);
	const subject = subjectName.split('-');
	console.log(subject[1]);

	const res = await fetch(`${APIURL}/subjects/${subject[1]}/lectures`);
	const data = await res.json();
	return {
		subjectData: data
	};
};
