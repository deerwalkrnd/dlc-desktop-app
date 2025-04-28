import { APIURL } from '$lib/constant';

export const load = async ({ params }: { params: any }) => {
	const subjectName = params.slug.toLowerCase();
	const subject = subjectName.split('-');

	const lectureRes = await fetch(`${APIURL}/subjects/${subject[1]}/lectures`);
	const lecturesData = await lectureRes.json();

	const lessonPromises = lecturesData.lectures.map(async (lecture: any) => {
		const lessonRes = await fetch(`${APIURL}/lectures/${lecture.ID}/lessons`);
		const lessonData = await lessonRes.json();
		return lessonData.lessons;
	});

	const lessons = await Promise.all(lessonPromises);
	return {
		subjectData: lecturesData,
		lessons
	};
};
