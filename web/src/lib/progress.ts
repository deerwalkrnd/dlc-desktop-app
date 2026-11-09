// Watch-progress rules: when to resume, when a lesson counts as watched, and
// how progress is stored. Pure TypeScript with no imports, so it runs under
// `node --test`. The reactive, localStorage-backed store is in
// progress.svelte.ts.

export interface LessonProgress {
	/** Last playback position, in seconds. */
	position: number;
	/** Video length in seconds (0 if unknown). */
	duration: number;
	/** True once the lesson has been watched to (nearly) the end. */
	completed: boolean;
	/** Milliseconds since the epoch; used for "Continue watching" order. */
	updatedAt: number;
	/** Display details, so "Continue watching" works without the API. */
	title: string;
	subtitle: string;
	href: string;
}

export type ProgressMap = Record<string, LessonProgress>;

export const STORAGE_KEY = 'dlc.progress.v1';

/** Positions closer than this to the start are not worth resuming. */
export const MIN_RESUME_SECONDS = 10;
/** A lesson counts as watched within this many seconds of the end (for
 * videos at least 5x this long)... */
export const END_MARGIN_SECONDS = 15;
/** ...or after this fraction of it has played. */
export const COMPLETE_FRACTION = 0.95;
/** Oldest entries beyond this are dropped to keep storage small. */
export const MAX_ENTRIES = 500;

export interface ProgressUpdate {
	position: number;
	duration: number;
	title: string;
	subtitle: string;
	href: string;
}

export function isComplete(position: number, duration: number): boolean {
	if (!(duration > 0)) return false;
	if (position / duration >= COMPLETE_FRACTION) return true;
	// The end margin (credits, "thank you") only makes sense for videos long
	// enough that it is a small part of them; short clips need 95%.
	return duration >= 5 * END_MARGIN_SECONDS && duration - position <= END_MARGIN_SECONDS;
}

/** Where to resume a lesson, or null to start from the beginning. */
export function resumePosition(p: LessonProgress | undefined): number | null {
	if (!p || p.completed) return null;
	if (p.position < MIN_RESUME_SECONDS) return null;
	if (p.duration > 0 && isComplete(p.position, p.duration)) return null;
	return p.position;
}

/** How much of the lesson has been watched, from 0 to 1. */
export function watchedFraction(p: LessonProgress | undefined): number {
	if (!p) return 0;
	if (p.completed) return 1;
	if (!(p.duration > 0)) return 0;
	return Math.min(1, Math.max(0, p.position / p.duration));
}

/**
 * Records a playback position. Returns a new map; the input is not modified.
 * Once a lesson is completed it stays completed, even if rewatched.
 */
export function record(map: ProgressMap, id: string, u: ProgressUpdate, now: number): ProgressMap {
	if (!Number.isFinite(u.position) || u.position < 0) return map;
	const previous = map[id];
	const duration =
		Number.isFinite(u.duration) && u.duration > 0 ? u.duration : (previous?.duration ?? 0);
	const entry: LessonProgress = {
		position: u.position,
		duration,
		completed: (previous?.completed ?? false) || isComplete(u.position, duration),
		updatedAt: now,
		title: u.title,
		subtitle: u.subtitle,
		href: u.href
	};
	return prune({ ...map, [id]: entry });
}

/** Marks a lesson as watched (e.g. when the video ends). */
export function complete(
	map: ProgressMap,
	id: string,
	u: ProgressUpdate,
	now: number
): ProgressMap {
	const next = record(map, id, u, now);
	return { ...next, [id]: { ...next[id], completed: true } };
}

/** Lessons started but not finished, most recent first. */
export function continueWatching(map: ProgressMap, limit: number): [string, LessonProgress][] {
	return Object.entries(map)
		.filter(([, p]) => resumePosition(p) !== null)
		.sort(([, a], [, b]) => b.updatedAt - a.updatedAt)
		.slice(0, limit);
}

function prune(map: ProgressMap): ProgressMap {
	const ids = Object.keys(map);
	if (ids.length <= MAX_ENTRIES) return map;
	const keep = ids.sort((a, b) => map[b].updatedAt - map[a].updatedAt).slice(0, MAX_ENTRIES);
	return Object.fromEntries(keep.map((id) => [id, map[id]]));
}

/** True for same-origin paths like "/class/8"; rejects "//host" and "/\\host". */
function isAppPath(href: unknown): href is string {
	return typeof href === 'string' && /^\/(?![/\\])/.test(href);
}

/** Parses stored progress, dropping anything malformed. Never throws. */
export function parseProgress(raw: string | null): ProgressMap {
	if (!raw) return {};
	let data: unknown;
	try {
		data = JSON.parse(raw);
	} catch {
		return {};
	}
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	const out: ProgressMap = {};
	for (const [id, v] of Object.entries(data as Record<string, unknown>)) {
		if (!v || typeof v !== 'object') continue;
		const p = v as Partial<LessonProgress>;
		if (typeof p.position !== 'number' || !Number.isFinite(p.position)) continue;
		out[id] = {
			position: Math.max(0, p.position),
			duration: typeof p.duration === 'number' && p.duration > 0 ? p.duration : 0,
			completed: p.completed === true,
			updatedAt: typeof p.updatedAt === 'number' ? p.updatedAt : 0,
			title: typeof p.title === 'string' ? p.title : '',
			subtitle: typeof p.subtitle === 'string' ? p.subtitle : '',
			href: isAppPath(p.href) ? p.href : '/'
		};
	}
	return out;
}
