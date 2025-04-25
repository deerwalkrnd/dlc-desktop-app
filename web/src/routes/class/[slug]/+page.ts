import { classes, getNestedData } from '$lib/data';
import { error } from '@sveltejs/kit';
export const load = async ({ params }: { params: any }) => {
	const classNumber = parseInt(params.slug, 10);
	const allClasses = getNestedData();

	const singleClass = allClasses.find((cls) => cls.id === classNumber);
	if (!singleClass) {
		return {
			status: 404,
			error: new Error(`Class ${classNumber} not found `)
		};
	}

	return {
		classData: singleClass
	};
};
