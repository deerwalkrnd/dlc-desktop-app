import { test } from 'node:test';
import assert from 'node:assert/strict';
import {
	MAX_ENTRIES,
	complete,
	continueWatching,
	isComplete,
	parseProgress,
	record,
	resumePosition,
	watchedFraction,
	type LessonProgress,
	type ProgressUpdate
} from './progress.ts';

const meta = {
	title: 'Lesson 1.2: Sets',
	subtitle: 'Chapter 1',
	href: '/class/8/subject/maths-new?lesson=a'
};
const at = (position: number, duration = 600): ProgressUpdate => ({ ...meta, position, duration });

test('resume only for meaningful, unfinished positions', () => {
	assert.equal(resumePosition(undefined), null);
	const m = record({}, 'a', at(5), 1);
	assert.equal(resumePosition(m.a), null, 'too close to the start');
	assert.equal(resumePosition(record({}, 'a', at(272), 1).a), 272);
	assert.equal(resumePosition(record({}, 'a', at(590), 1).a), null, 'within the end margin');
});

test('completion: 95% or the last 15 seconds', () => {
	assert.equal(isComplete(569, 600), false);
	assert.equal(isComplete(570, 600), true);
	assert.equal(isComplete(25, 40), false, 'short video: the end margin does not apply');
	assert.equal(isComplete(38, 40), true, 'short video at 95%');
	assert.equal(isComplete(61, 75), true, '75 s video: end margin applies');
	assert.equal(isComplete(12.9, 27), false, 'regression: 27 s demo clip at 0:12');
	assert.equal(isComplete(100, 0), false, 'unknown duration');
	assert.equal(record({}, 'a', at(590), 1).a.completed, true);
});

test('completed lessons stay completed when rewatched', () => {
	let m = complete({}, 'a', at(600), 1);
	m = record(m, 'a', at(30), 2);
	assert.equal(m.a.completed, true);
	assert.equal(m.a.position, 30);
	assert.equal(resumePosition(m.a), null);
	assert.equal(watchedFraction(m.a), 1);
});

test('record keeps a known duration when the update has none', () => {
	let m = record({}, 'a', at(100, 600), 1);
	m = record(m, 'a', at(120, NaN), 2);
	assert.equal(m.a.duration, 600);
	assert.equal(watchedFraction(m.a), 0.2);
});

test('record ignores invalid positions and does not mutate its input', () => {
	const m = record({}, 'a', at(100), 1);
	const frozen = Object.freeze({ ...m });
	assert.equal(record(frozen, 'a', at(NaN), 2), frozen);
	assert.equal(record(frozen, 'a', at(-1), 2), frozen);
	const next = record(frozen, 'b', at(50), 3);
	assert.notEqual(next, frozen);
	assert.equal(Object.keys(frozen).length, 1);
});

test('continue watching: unfinished lessons, newest first', () => {
	let m = record({}, 'old', at(100), 1);
	m = record(m, 'new', at(200), 3);
	m = record(m, 'barely', at(3), 4);
	m = complete(m, 'done', at(600), 5);
	assert.deepEqual(
		continueWatching(m, 10).map(([id]) => id),
		['new', 'old']
	);
	assert.equal(continueWatching(m, 1).length, 1);
});

test('storage is pruned to the newest entries', () => {
	let m = {};
	for (let i = 0; i < MAX_ENTRIES + 5; i++) m = record(m, `id${i}`, at(100), i);
	assert.equal(Object.keys(m).length, MAX_ENTRIES);
	assert.ok(!('id0' in m) && `id${MAX_ENTRIES + 4}` in m);
});

test('parseProgress survives anything', () => {
	for (const raw of [
		null,
		'',
		'not json',
		'[]',
		'42',
		'null',
		'{"a": 5}',
		'{"a": {"position": "x"}}'
	]) {
		assert.deepEqual(parseProgress(raw), {}, String(raw));
	}
	const good: LessonProgress = {
		position: 10,
		duration: 60,
		completed: false,
		updatedAt: 1,
		...meta
	};
	assert.deepEqual(parseProgress(JSON.stringify({ a: good })), { a: good });
	// Unsafe or missing fields are sanitised.
	const odd = parseProgress(JSON.stringify({ a: { position: 10, href: 'https://evil.example' } }));
	assert.equal(odd.a.href, '/');
	for (const href of ['//evil.example/x', '/\\evil.example', 'javascript:alert(1)', 42]) {
		const r = parseProgress(JSON.stringify({ a: { position: 10, href } }));
		assert.equal(r.a.href, '/', String(href));
	}
	assert.equal(
		parseProgress(JSON.stringify({ a: { position: 10, href: '/class/8?lesson=x' } })).a.href,
		'/class/8?lesson=x'
	);
	assert.equal(odd.a.title, '');
});
