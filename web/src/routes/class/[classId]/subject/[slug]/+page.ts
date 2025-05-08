import { APIURL } from '$lib/constant';

export const load = async ({ params }: { params: any }) => {
	const classId = params.classId;
	const subjectSlug = params.slug;
	const subjectParts = subjectSlug.split('-');
	const subjectId = subjectParts[1];
	const subjectName = subjectParts[0];

	const lectureRes = await fetch(`${APIURL}/subjects/${subjectId}/lectures`);
	const lecturesData = await lectureRes.json();

	async function fetchClassSubjects() {
		try {
			const classesRes = await fetch(`${APIURL}/classes`);
			const classList = await classesRes.json();
			const currentClass = classList.classes.find((cls: any) => cls.Number == classId);

			if (currentClass) {
				const type = 'new';
				const subjectsRes = await fetch(
					`${APIURL}/classes/${currentClass.ID}/subjects?type=${type}`
				);
				return await subjectsRes.json();
			}
			return null;
		} catch (error) {
			console.error('Error fetching class subjects:', error);
			return null;
		}
	}

	const subjectList = await fetchClassSubjects();
	lecturesData.lectures.sort(function (a: any, b: any) {
		return a.number - b.number;
	});
	return {
		subjectName,
		subjectList,
		subjectData: lecturesData
		// lessons
	};
};
