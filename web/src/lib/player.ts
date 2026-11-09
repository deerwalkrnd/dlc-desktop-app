// Video player state machine and helpers. Pure TypeScript with no imports,
// so it runs under `node --test` without a bundler.
//
// Transitions (anything not listed is ignored):
//
//   loading   --CAN_PLAY--> paused     --PLAY--> playing   --ERROR--> error
//   paused    --PLAY------> playing    --ENDED-> ended     --ERROR--> error
//   playing   --PAUSE-----> paused     --WAITING-> buffering
//             --ENDED-----> ended      --ERROR--> error
//   buffering --PLAY|CAN_PLAY--> playing   --PAUSE--> paused
//             --ENDED-----> ended      --ERROR--> error
//   ended     --PLAY------> playing
//
// LOAD is accepted from every state (switching lessons, retrying after an
// error). Events that make no sense in a state are ignored.

export type PlayerStatus =
	'idle' | 'loading' | 'paused' | 'playing' | 'buffering' | 'ended' | 'error';

export type PlayerEvent =
	| { type: 'LOAD' }
	| { type: 'CAN_PLAY' }
	| { type: 'PLAY' }
	| { type: 'PAUSE' }
	| { type: 'WAITING' }
	| { type: 'ENDED' }
	| { type: 'ERROR'; message: string }
	| { type: 'UNLOAD' };

export interface PlayerState {
	status: PlayerStatus;
	error: string | null;
}

export const initialPlayerState: PlayerState = { status: 'idle', error: null };

const table: Record<PlayerStatus, Partial<Record<PlayerEvent['type'], PlayerStatus>>> = {
	idle: {},
	loading: { CAN_PLAY: 'paused', PLAY: 'playing', ERROR: 'error' },
	paused: { PLAY: 'playing', ENDED: 'ended', ERROR: 'error' },
	playing: { PAUSE: 'paused', WAITING: 'buffering', ENDED: 'ended', ERROR: 'error' },
	buffering: {
		PLAY: 'playing',
		CAN_PLAY: 'playing',
		PAUSE: 'paused',
		ENDED: 'ended',
		ERROR: 'error'
	},
	ended: { PLAY: 'playing' },
	error: {}
};

export function transition(state: PlayerState, event: PlayerEvent): PlayerState {
	if (event.type === 'LOAD') return { status: 'loading', error: null };
	if (event.type === 'UNLOAD') return initialPlayerState;
	const next = table[state.status][event.type];
	if (!next || next === state.status) return state;
	return { status: next, error: event.type === 'ERROR' ? event.message : null };
}

/** True while the video is (or is trying to be) playing. */
export function isActive(status: PlayerStatus): boolean {
	return status === 'playing' || status === 'buffering';
}

/** True when transport controls (seek, skip) make sense. */
export function isSeekable(status: PlayerStatus): boolean {
	return (
		status === 'paused' || status === 'playing' || status === 'buffering' || status === 'ended'
	);
}

/** Human-readable explanation for HTMLMediaElement.error.code. */
export function mediaErrorMessage(code: number | undefined): string {
	switch (code) {
		case 1:
			return 'Playback was stopped.';
		case 2:
			return 'The video could not be read. Check that the app is still running.';
		case 3:
			return 'This video file is damaged or uses an unsupported format.';
		case 4:
			return 'This video could not be found. It may have been moved or renamed; try rescanning the library.';
		default:
			return 'This video could not be played.';
	}
}

export type KeyAction =
	| 'toggle'
	| 'back'
	| 'forward'
	| 'volumeUp'
	| 'volumeDown'
	| 'mute'
	| 'fullscreen'
	| 'next'
	| 'previous';

/** Maps a KeyboardEvent.key to a player action. */
export function keyAction(key: string, shift = false): KeyAction | null {
	switch (key) {
		case ' ':
		case 'k':
		case 'K':
			return 'toggle';
		case 'ArrowLeft':
		case 'j':
		case 'J':
			return 'back';
		case 'ArrowRight':
		case 'l':
		case 'L':
			return 'forward';
		case 'ArrowUp':
			return 'volumeUp';
		case 'ArrowDown':
			return 'volumeDown';
		case 'm':
		case 'M':
			return 'mute';
		case 'f':
		case 'F':
			return 'fullscreen';
		case 'N':
			return shift ? 'next' : null;
		case 'P':
			return shift ? 'previous' : null;
	}
	return null;
}

export const SKIP_SECONDS = 10;
export const SPEEDS = [0.5, 0.75, 1, 1.25, 1.5, 2] as const;

/** Formats seconds as m:ss or h:mm:ss. */
export function formatTime(seconds: number): string {
	if (!Number.isFinite(seconds) || seconds < 0) seconds = 0;
	const s = Math.floor(seconds % 60);
	const m = Math.floor((seconds / 60) % 60);
	const h = Math.floor(seconds / 3600);
	const ss = String(s).padStart(2, '0');
	return h > 0 ? `${h}:${String(m).padStart(2, '0')}:${ss}` : `${m}:${ss}`;
}

export function clamp(value: number, min: number, max: number): number {
	return Math.min(max, Math.max(min, value));
}
