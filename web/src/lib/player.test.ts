// Run with: npm test   (uses Node's built-in test runner and type stripping)
import { test } from 'node:test';
import assert from 'node:assert/strict';
import {
	clamp,
	formatTime,
	initialPlayerState,
	isActive,
	keyAction,
	transition,
	type PlayerEvent,
	type PlayerState
} from './player.ts';

function run(...events: PlayerEvent[]): PlayerState {
	return events.reduce(transition, initialPlayerState);
}

test('happy path: load, play, pause, end, replay', () => {
	assert.equal(run({ type: 'LOAD' }).status, 'loading');
	assert.equal(run({ type: 'LOAD' }, { type: 'CAN_PLAY' }).status, 'paused');
	assert.equal(run({ type: 'LOAD' }, { type: 'CAN_PLAY' }, { type: 'PLAY' }).status, 'playing');
	assert.equal(
		run({ type: 'LOAD' }, { type: 'PLAY' }, { type: 'PAUSE' }, { type: 'ENDED' }).status,
		'ended'
	);
	assert.equal(
		run({ type: 'LOAD' }, { type: 'PLAY' }, { type: 'ENDED' }, { type: 'PLAY' }).status,
		'playing'
	);
});

test('autoplay goes straight from loading to playing', () => {
	assert.equal(run({ type: 'LOAD' }, { type: 'PLAY' }).status, 'playing');
});

test('buffering', () => {
	const buffering = run({ type: 'LOAD' }, { type: 'PLAY' }, { type: 'WAITING' });
	assert.equal(buffering.status, 'buffering');
	assert.ok(isActive(buffering.status));
	assert.equal(transition(buffering, { type: 'PLAY' }).status, 'playing');
	assert.equal(transition(buffering, { type: 'CAN_PLAY' }).status, 'playing');
	assert.equal(transition(buffering, { type: 'PAUSE' }).status, 'paused');
});

test('errors carry a message and are cleared by LOAD', () => {
	const failed = run({ type: 'LOAD' }, { type: 'ERROR', message: 'missing' });
	assert.deepEqual(failed, { status: 'error', error: 'missing' });
	// Late media events must not hide the error.
	assert.equal(transition(failed, { type: 'CAN_PLAY' }), failed);
	assert.equal(transition(failed, { type: 'PLAY' }), failed);
	assert.deepEqual(transition(failed, { type: 'LOAD' }), { status: 'loading', error: null });
});

test('UNLOAD always returns to idle', () => {
	for (const s of [
		run({ type: 'LOAD' }),
		run({ type: 'LOAD' }, { type: 'PLAY' }),
		run({ type: 'LOAD' }, { type: 'ERROR', message: 'x' })
	]) {
		assert.deepEqual(transition(s, { type: 'UNLOAD' }), initialPlayerState);
	}
});

test('idle ignores media events', () => {
	for (const type of ['CAN_PLAY', 'PLAY', 'PAUSE', 'WAITING', 'ENDED'] as const) {
		assert.equal(transition(initialPlayerState, { type }), initialPlayerState);
	}
});

test('irrelevant events keep the same object (no re-render)', () => {
	const paused = run({ type: 'LOAD' }, { type: 'CAN_PLAY' });
	assert.equal(transition(paused, { type: 'PAUSE' }), paused);
	assert.equal(transition(paused, { type: 'CAN_PLAY' }), paused);
});

test('keyboard shortcuts', () => {
	assert.equal(keyAction(' '), 'toggle');
	assert.equal(keyAction('k'), 'toggle');
	assert.equal(keyAction('ArrowLeft'), 'back');
	assert.equal(keyAction('ArrowRight'), 'forward');
	assert.equal(keyAction('m'), 'mute');
	assert.equal(keyAction('f'), 'fullscreen');
	assert.equal(keyAction('N', true), 'next');
	assert.equal(keyAction('N', false), null);
	assert.equal(keyAction('Escape'), null);
});

test('formatTime', () => {
	assert.equal(formatTime(0), '0:00');
	assert.equal(formatTime(9.9), '0:09');
	assert.equal(formatTime(75), '1:15');
	assert.equal(formatTime(3725), '1:02:05');
	assert.equal(formatTime(NaN), '0:00');
	assert.equal(formatTime(Infinity), '0:00');
});

test('clamp', () => {
	assert.equal(clamp(-1, 0, 10), 0);
	assert.equal(clamp(11, 0, 10), 10);
	assert.equal(clamp(5, 0, 10), 5);
});
