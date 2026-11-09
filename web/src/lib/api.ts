// Typed client for the Go server's JSON API (internal/server/api.go).
// All URLs are relative: in production the app and API share an origin, and
// in development Vite proxies /api and /media to the Go server.
import { error } from '@sveltejs/kit';

export type Syllabus = 'old' | 'new';

export interface ClassSummary {
	number: number;
	subjectCount: number;
	lessonCount: number;
	syllabi: Syllabus[];
}

export interface SubjectSummary {
	slug: string;
	name: string;
	syllabus: Syllabus;
	chapterCount: number;
	lessonCount: number;
}

export interface ClassDetail {
	number: number;
	syllabi: Syllabus[];
	subjects: SubjectSummary[];
}

export interface Lesson {
	id: string;
	number: string;
	title: string;
	teacher: string;
	mediaUrl: string;
}

export interface Chapter {
	number: number;
	title: string;
	lessons: Lesson[];
}

export interface SubjectDetail {
	class: number;
	slug: string;
	name: string;
	syllabus: Syllabus;
	chapters: Chapter[];
}

export interface Teacher {
	name: string;
	lessonCount: number;
	subjects: string[];
}

export interface Issue {
	path: string;
	message: string;
}

export interface LibraryStatus {
	path: string;
	exists: boolean;
	scannedAt: string;
	scanMillis: number;
	classCount: number;
	lessonCount: number;
	skipped: Issue[];
	warnings: Issue[];
	version: string;
}

export class ApiError extends Error {
	status: number;
	constructor(status: number, message: string) {
		super(message);
		this.status = status;
	}
}

type Fetch = typeof fetch;

async function request<T>(fetch: Fetch, path: string, init?: RequestInit): Promise<T> {
	let res: Response;
	try {
		res = await fetch(path, { ...init, headers: { Accept: 'application/json', ...init?.headers } });
	} catch {
		throw new ApiError(
			503,
			'Cannot reach the DLC app. Make sure dlc-desktop-app is still running.'
		);
	}
	if (!res.ok) {
		let message = res.statusText || `Request failed (${res.status})`;
		try {
			const body = await res.json();
			if (typeof body?.error === 'string') message = body.error;
		} catch {
			// not JSON; keep statusText
		}
		throw new ApiError(res.status, message);
	}
	return res.json() as Promise<T>;
}

export const api = {
	classes: (f: Fetch) =>
		request<{ classes: ClassSummary[] }>(f, '/api/classes').then((r) => r.classes),
	class: (f: Fetch, number: number | string) =>
		request<ClassDetail>(f, `/api/classes/${encodeURIComponent(number)}`),
	subject: (f: Fetch, number: number | string, slug: string) =>
		request<SubjectDetail>(
			f,
			`/api/classes/${encodeURIComponent(number)}/subjects/${encodeURIComponent(slug)}`
		),
	teachers: (f: Fetch) =>
		request<{ teachers: Teacher[] }>(f, '/api/teachers').then((r) => r.teachers),
	library: (f: Fetch) => request<LibraryStatus>(f, '/api/library'),
	rescan: (f: Fetch) => request<LibraryStatus>(f, '/api/library/rescan', { method: 'POST' })
};

/** Converts an ApiError into a SvelteKit error page; use inside load(). */
export async function orErrorPage<T>(promise: Promise<T>): Promise<T> {
	try {
		return await promise;
	} catch (e) {
		if (e instanceof ApiError) error(e.status, e.message);
		throw e;
	}
}

export function syllabusLabel(s: Syllabus): string {
	return s === 'old' ? 'Old syllabus' : 'New syllabus';
}

/** Path to a subject page. */
export function subjectHref(classNumber: number, slug: string): string {
	return `/class/${classNumber}/subject/${encodeURIComponent(slug)}`;
}
