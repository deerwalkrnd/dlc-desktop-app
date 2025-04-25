// src/lib/data.ts
export type Class = {
	id: number;
	number: number;
	created_at: Date;
	updated_at: Date;
	deleted_at: Date | null;
	subjects?: SubjectWithRelations[];
};

export type Subject = {
	id: number;
	name: string;
	type: number;
	class_id: number;
	created_at: Date;
	updated_at: Date;
	deleted_at: Date | null;
};

export type Lecture = {
	id: number;
	number: number;
	name: string;
	subject_id: number;
	created_at: Date;
	updated_at: Date;
	deleted_at: Date | null;
};

export type Teacher = {
	id: number;
	name: string;
	created_at: Date;
	updated_at: Date;
	deleted_at: Date | null;
};

export type Lesson = {
	id: number;
	name: string;
	number: number;
	video_url: string;
	teacher_id: number;
	lecture_id: number;
	created_at: Date;
	updated_at: Date;
	deleted_at: Date | null;
};

export type LessonWithTeacher = Lesson & { teacher: Teacher | null };
export type LectureWithLessons = Lecture & { lessons: LessonWithTeacher[] };
export type SubjectWithRelations = Subject & { lectures: LectureWithLessons[] };

export const classes: Class[] = [
	{ id: 1, number: 10, created_at: new Date(), updated_at: new Date(), deleted_at: null }
];

export const subjects: Subject[] = [
	{
		id: 1,
		name: 'Mathematics',
		type: 1,
		class_id: 1,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 2,
		name: 'Science',
		type: 1,
		class_id: 1,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 3,
		name: 'English',
		type: 2,
		class_id: 1,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 4,
		name: 'History',
		type: 2,
		class_id: 2,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 5,
		name: 'Computer Science',
		type: 1,
		class_id: 2,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 6,
		name: 'Geography',
		type: 2,
		class_id: 2,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 7,
		name: 'Physics',
		type: 1,
		class_id: 3,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 8,
		name: 'Chemistry',
		type: 1,
		class_id: 3,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	},
	{
		id: 9,
		name: 'Biology',
		type: 1,
		class_id: 3,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	}
];

export const lectures: Lecture[] = [
	{
		id: 1,
		number: 1,
		name: 'Algebra Basics',
		subject_id: 1,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	}
];

export const teachers: Teacher[] = [
	{ id: 1, name: 'Mr. Smith', created_at: new Date(), updated_at: new Date(), deleted_at: null }
];

export const lessons: Lesson[] = [
	{
		id: 1,
		name: 'Introduction to Variables',
		number: 1.1,
		video_url: 'https://example.com/video1',
		teacher_id: 1,
		lecture_id: 1,
		created_at: new Date(),
		updated_at: new Date(),
		deleted_at: null
	}
];

/**
 * Build a nested data structure:
 * Class → [Subject → [Lecture → [Lesson + Teacher]]]
 */
export function getNestedData(): Class[] {
	return classes.map((cls) => {
		const clsSubjects = subjects
			.filter((s) => s.class_id === cls.id)
			.map<SubjectWithRelations>((s) => {
				const subjLectures = lectures
					.filter((l) => l.subject_id === s.id)
					.map<LectureWithLessons>((l) => {
						const lsn = lessons
							.filter((ls) => ls.lecture_id === l.id)
							.map<LessonWithTeacher>((ls) => ({
								...ls,
								teacher: teachers.find((t) => t.id === ls.teacher_id) || null
							}));
						return { ...l, lessons: lsn };
					});
				return { ...s, lectures: subjLectures };
			});

		return { ...cls, subjects: clsSubjects };
	});
}
