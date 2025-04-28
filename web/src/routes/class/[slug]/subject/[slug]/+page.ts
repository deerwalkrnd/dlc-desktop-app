export const load = async ({ params }: { params: any }) => {
	const subjectName = params.slug.toLowerCase();

	return {
		subjectName: subjectName
	};
};
