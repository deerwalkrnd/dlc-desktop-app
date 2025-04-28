import { APIURL } from '$lib/constant';
import { classes } from '$lib/data';
import { error } from '@sveltejs/kit';
export const load = async ({ params }: { params: any }) => {
	const classNumber = parseInt(params.slug, 10);
	// const allClasses = getNestedData();

	// const singleClass = allClasses.find((cls) => cls.id === classNumber);
	const res = await fetch(`${APIURL}/classes`);
	const data = await res.json();

	let getClass = data.classes.find((singleClass: any) => singleClass.Number === classNumber);
	return {
		getClass
	};
};
