<!--
	Modal video player. All playback state lives in the explicit state machine
	in $lib/player.ts; this component only wires <video> events into it and
	renders the current state. Listeners are scoped to the dialog and removed
	automatically when the component is destroyed.
-->
<script lang="ts" module>
	export interface PlayableLesson {
		id: string;
		mediaUrl: string;
		title: string;
		subtitle?: string;
	}
</script>

<script lang="ts">
	import { untrack } from 'svelte';
	import Icon from './Icon.svelte';
	import {
		SKIP_SECONDS,
		SPEEDS,
		clamp,
		formatTime,
		initialPlayerState,
		isActive,
		isSeekable,
		keyAction,
		mediaErrorMessage,
		transition,
		type KeyAction,
		type PlayerEvent
	} from '$lib/player';

	interface Props {
		lesson: PlayableLesson | null;
		hasPrevious?: boolean;
		hasNext?: boolean;
		/** Short description of the previous/next lesson, shown in the footer. */
		previousLabel?: string;
		nextLabel?: string;
		onclose: () => void;
		onprevious?: () => void;
		onnext?: () => void;
	}

	let {
		lesson,
		hasPrevious = false,
		hasNext = false,
		previousLabel,
		nextLabel,
		onclose,
		onprevious,
		onnext
	}: Props = $props();

	let dialog: HTMLDialogElement;
	let video: HTMLVideoElement;
	let stage: HTMLDivElement;

	let machine = $state(initialPlayerState);
	let currentTime = $state(0);
	let duration = $state(0);
	let bufferedEnd = $state(0);
	let volume = $state(loadPref('dlc.volume', 1));
	let muted = $state(false);
	let rate = $state(loadPref('dlc.rate', 1));
	let fullscreen = $state(false);
	let scrubbing = $state(false);
	let pointerIdle = $state(false);

	const active = $derived(isActive(machine.status));
	const seekable = $derived(isSeekable(machine.status));
	const showControls = $derived(!active || !pointerIdle || scrubbing);

	function send(event: PlayerEvent) {
		machine = transition(
			untrack(() => machine),
			event
		);
	}

	// React to the selected lesson: open/close the dialog and swap sources.
	$effect(() => {
		const current = lesson;
		untrack(() => (current ? start(current) : stop()));
	});

	function start(l: PlayableLesson) {
		if (!dialog.open) dialog.showModal();
		// showModal() focuses the first button ("Next lesson"), which would make
		// Space skip ahead. Focus the stage so shortcuts control playback.
		stage.focus({ preventScroll: true });
		currentTime = 0;
		duration = 0;
		bufferedEnd = 0;
		video.src = l.mediaUrl;
		video.playbackRate = rate;
		video.volume = volume;
		send({ type: 'LOAD' });
		video.load();
		// The click that opened the player counts as a user gesture, so
		// autoplay is allowed; if the browser refuses, we simply stay paused.
		video.play().catch(() => {});
	}

	function stop() {
		send({ type: 'UNLOAD' });
		if (video) {
			video.pause();
			video.removeAttribute('src');
			video.load(); // releases the file and network connection
		}
		if (document.fullscreenElement) document.exitFullscreen().catch(() => {});
		if (dialog?.open) dialog.close();
	}

	function retry() {
		if (lesson) start(lesson);
	}

	function toggle() {
		if (!seekable && machine.status !== 'loading') return;
		if (machine.status === 'ended') video.currentTime = 0;
		if (video.paused) video.play().catch(() => {});
		else video.pause();
	}

	function seekTo(t: number) {
		if (!seekable) return;
		video.currentTime = clamp(t, 0, duration || 0);
		currentTime = video.currentTime;
	}

	function setVolume(v: number) {
		video.volume = clamp(v, 0, 1);
		if (video.volume > 0) video.muted = false;
	}

	async function toggleFullscreen() {
		try {
			if (document.fullscreenElement) await document.exitFullscreen();
			else await stage.requestFullscreen();
		} catch {
			// fullscreen unavailable (e.g. blocked by the browser)
		}
	}

	function perform(action: KeyAction) {
		switch (action) {
			case 'toggle':
				return toggle();
			case 'back':
				return seekTo(video.currentTime - SKIP_SECONDS);
			case 'forward':
				return seekTo(video.currentTime + SKIP_SECONDS);
			case 'volumeUp':
				return setVolume(video.volume + 0.1);
			case 'volumeDown':
				return setVolume(video.volume - 0.1);
			case 'mute':
				video.muted = !video.muted;
				return;
			case 'fullscreen':
				return toggleFullscreen();
			case 'next':
				return hasNext && onnext?.();
			case 'previous':
				return hasPrevious && onprevious?.();
		}
	}

	function onkeydown(e: KeyboardEvent) {
		if (e.ctrlKey || e.metaKey || e.altKey) return;
		const target = e.target as HTMLElement;
		// Let focused controls handle their own activation keys.
		if (target instanceof HTMLButtonElement && (e.key === ' ' || e.key === 'Enter')) return;
		if (target instanceof HTMLSelectElement) return;
		if (target instanceof HTMLInputElement && target.dataset.role !== 'seek') {
			if (e.key !== ' ') return; // volume slider keeps its arrow keys
		}
		const action = keyAction(e.key, e.shiftKey);
		if (!action) return;
		e.preventDefault();
		perform(action);
	}

	// Hide the controls after a moment of inactivity while playing.
	let idleTimer: ReturnType<typeof setTimeout> | undefined;
	function wake() {
		pointerIdle = false;
		clearTimeout(idleTimer);
		idleTimer = setTimeout(() => (pointerIdle = true), 2500);
	}

	$effect(() => {
		const onFullscreen = () => (fullscreen = document.fullscreenElement === stage);
		document.addEventListener('fullscreenchange', onFullscreen);
		return () => {
			document.removeEventListener('fullscreenchange', onFullscreen);
			clearTimeout(idleTimer);
		};
	});

	// Remember volume and speed between lessons and visits.
	$effect(() => savePref('dlc.volume', volume));
	$effect(() => savePref('dlc.rate', rate));

	function loadPref(key: string, fallback: number): number {
		try {
			const v = Number(localStorage.getItem(key));
			return localStorage.getItem(key) !== null && Number.isFinite(v) ? v : fallback;
		} catch {
			return fallback;
		}
	}

	function savePref(key: string, value: number) {
		try {
			localStorage.setItem(key, String(value));
		} catch {
			// storage unavailable; preferences just won't persist
		}
	}

	function updateBuffered() {
		const b = video.buffered;
		for (let i = 0; i < b.length; i++) {
			if (b.start(i) <= video.currentTime && video.currentTime <= b.end(i)) {
				bufferedEnd = b.end(i);
				return;
			}
		}
	}

	// On-video controls: white on the navy gradient.
	const control =
		'grid size-10 place-items-center rounded-full text-white transition hover:bg-white/15 disabled:opacity-40';
	// Footer lesson navigation: white cards on the light-blue strip.
	const navButton =
		'group flex min-w-0 items-center gap-3 rounded-lg bg-white px-3 py-2 text-left ring-1 ring-brand-100 transition hover:bg-brand-50 hover:ring-brand-300 disabled:invisible';

	const progress = $derived(duration ? (currentTime / duration) * 100 : 0);
	const buffered = $derived(duration ? (bufferedEnd / duration) * 100 : 0);
</script>

<dialog
	bind:this={dialog}
	class="ring-brand-100 backdrop:bg-brand-950/60 m-auto w-full max-w-5xl overflow-hidden rounded-2xl bg-white p-0 text-slate-900 shadow-2xl ring-1 backdrop:backdrop-blur-sm"
	aria-label={lesson?.title ?? 'Video player'}
	{onkeydown}
	onclose={() => lesson && onclose()}
	onclick={(e) => e.target === dialog && onclose()}
>
	<header class="flex items-start gap-3 px-6 py-4">
		<div class="min-w-0 flex-1">
			<h2 class="text-brand-800 truncate text-lg font-semibold">{lesson?.title}</h2>
			{#if lesson?.subtitle}
				<p class="truncate text-sm text-slate-500">{lesson.subtitle}</p>
			{/if}
		</div>
		<button
			class="hover:bg-brand-50 hover:text-brand-700 -mr-2 rounded-full p-2 text-slate-500 transition"
			onclick={onclose}
			aria-label="Close player"
			title="Close (Esc)"
		>
			<Icon name="close" />
		</button>
	</header>

	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		bind:this={stage}
		tabindex="-1"
		class="group relative bg-black outline-none {showControls ? '' : 'cursor-none'}"
		onpointermove={wake}
		onpointerdown={wake}
	>
		<!-- Click toggles playback; keyboard users have the buttons and shortcuts below. -->
		<video
			bind:this={video}
			class="aspect-video max-h-[80vh] w-full bg-black {fullscreen ? 'h-full max-h-none' : ''}"
			preload="metadata"
			playsinline
			onclick={toggle}
			ondblclick={toggleFullscreen}
			oncanplay={() => send({ type: 'CAN_PLAY' })}
			onplay={() => send({ type: 'PLAY' })}
			onplaying={() => send({ type: 'PLAY' })}
			onpause={() => send({ type: 'PAUSE' })}
			onwaiting={() => send({ type: 'WAITING' })}
			onended={() => send({ type: 'ENDED' })}
			onerror={() =>
				lesson && send({ type: 'ERROR', message: mediaErrorMessage(video.error?.code) })}
			ontimeupdate={() => {
				if (!scrubbing) currentTime = video.currentTime;
				updateBuffered();
			}}
			ondurationchange={() => (duration = Number.isFinite(video.duration) ? video.duration : 0)}
			onprogress={updateBuffered}
			onvolumechange={() => {
				volume = video.volume;
				muted = video.muted;
			}}
			onratechange={() => (rate = video.playbackRate)}
		>
			<track kind="captions" />
		</video>

		<!-- State overlays -->
		{#if machine.status === 'loading' || machine.status === 'buffering'}
			<div class="pointer-events-none absolute inset-0 grid place-items-center" role="status">
				<span
					class="border-t-brand-300 size-14 animate-spin rounded-full border-4 border-white/25"
					aria-label="Loading"
				></span>
			</div>
		{:else if machine.status === 'error'}
			<div
				class="bg-brand-950/85 absolute inset-0 flex flex-col items-center justify-center gap-4 p-6 text-center text-white"
				role="alert"
			>
				<span class="grid size-14 place-items-center rounded-full bg-white/10">
					<Icon name="alert" class="text-brand-200 size-8" />
				</span>
				<p class="max-w-md text-lg">{machine.error}</p>
				<button
					class="bg-brand-600 hover:bg-brand-500 flex items-center gap-2 rounded-lg px-4 py-2 font-medium text-white"
					onclick={retry}><Icon name="refresh" class="size-4" /> Try again</button
				>
			</div>
		{:else if machine.status === 'ended'}
			<div class="bg-brand-950/75 absolute inset-0 flex items-center justify-center gap-3">
				<button
					class="text-brand-800 hover:bg-brand-50 flex items-center gap-2 rounded-lg bg-white px-4 py-2 font-medium"
					onclick={toggle}
				>
					<Icon name="replay" /> Watch again
				</button>
				{#if hasNext}
					<button
						class="bg-brand-600 hover:bg-brand-500 flex items-center gap-2 rounded-lg px-4 py-2 font-medium text-white"
						onclick={onnext}
					>
						Next lesson <Icon name="next" />
					</button>
				{/if}
			</div>
		{:else if machine.status === 'paused'}
			<button
				class="bg-brand-600/90 hover:bg-brand-600 absolute inset-0 m-auto grid size-20 place-items-center rounded-full text-white shadow-lg ring-4 ring-white/30 transition hover:scale-105"
				onclick={toggle}
				aria-label="Play"
			>
				<Icon name="play" class="size-10 translate-x-0.5" />
			</button>
		{/if}

		<!-- Controls -->
		<div
			class="from-brand-950/90 via-brand-950/50 absolute inset-x-0 bottom-0 bg-linear-to-t to-transparent px-4 pt-12 pb-3 text-white transition-opacity duration-200 {showControls
				? 'opacity-100'
				: 'pointer-events-none opacity-0'}"
		>
			<input
				type="range"
				data-role="seek"
				class="seek mb-2 w-full"
				min="0"
				max={duration || 0}
				step="0.1"
				value={currentTime}
				disabled={!seekable}
				aria-label="Seek"
				aria-valuetext="{formatTime(currentTime)} of {formatTime(duration)}"
				style="--progress: {progress}%; --buffered: {buffered}%"
				oninput={(e) => {
					scrubbing = true;
					seekTo(Number(e.currentTarget.value));
				}}
				onchange={() => (scrubbing = false)}
				onpointerup={() => (scrubbing = false)}
			/>

			<div class="flex items-center gap-1 sm:gap-2">
				<button
					class={control}
					onclick={toggle}
					aria-label={active ? 'Pause' : 'Play'}
					title="{active ? 'Pause' : 'Play'} (Space)"
				>
					<Icon name={active ? 'pause' : 'play'} class="size-6" />
				</button>
				<button
					class="{control} relative"
					onclick={() => perform('back')}
					disabled={!seekable}
					aria-label="Back {SKIP_SECONDS} seconds"
					title="Back {SKIP_SECONDS}s (←)"
				>
					<Icon name="back" class="size-6" />
					<span class="absolute inset-0 grid place-items-center pt-0.5 text-[9px] font-bold"
						>{SKIP_SECONDS}</span
					>
				</button>
				<button
					class="{control} relative"
					onclick={() => perform('forward')}
					disabled={!seekable}
					aria-label="Forward {SKIP_SECONDS} seconds"
					title="Forward {SKIP_SECONDS}s (→)"
				>
					<Icon name="forward" class="size-6" />
					<span class="absolute inset-0 grid place-items-center pt-0.5 text-[9px] font-bold"
						>{SKIP_SECONDS}</span
					>
				</button>

				<button
					class={control}
					onclick={() => perform('mute')}
					aria-label={muted || volume === 0 ? 'Unmute' : 'Mute'}
					title="Mute (M)"
				>
					<Icon name={muted || volume === 0 ? 'muted' : 'volume'} class="size-6" />
				</button>
				<input
					type="range"
					class="volume hidden w-24 sm:block"
					min="0"
					max="1"
					step="0.05"
					value={muted ? 0 : volume}
					aria-label="Volume"
					oninput={(e) => setVolume(Number(e.currentTarget.value))}
				/>

				<span class="ml-2 text-sm text-white/90 tabular-nums">
					{formatTime(currentTime)} / {formatTime(duration)}
				</span>

				<span class="flex-1"></span>

				<label class="flex items-center gap-1 text-sm text-white">
					<span class="sr-only">Playback speed</span>
					<select
						class="rounded-full bg-white/15 px-2.5 py-1 text-white hover:bg-white/25"
						value={rate}
						onchange={(e) => (video.playbackRate = Number(e.currentTarget.value))}
					>
						{#each SPEEDS as s (s)}
							<option value={s} class="text-slate-900">{s}×</option>
						{/each}
					</select>
				</label>

				<button
					class={control}
					onclick={toggleFullscreen}
					aria-label={fullscreen ? 'Exit full screen' : 'Full screen'}
					title="Full screen (F)"
				>
					<Icon name={fullscreen ? 'minimize' : 'maximize'} class="size-6" />
				</button>
			</div>
		</div>
	</div>
	<footer
		class="border-brand-100 bg-brand-50 flex items-center justify-between gap-4 border-t px-6 py-3"
	>
		<button
			class={navButton}
			onclick={onprevious}
			disabled={!hasPrevious}
			title="Previous lesson (Shift+P)"
		>
			<Icon name="previous" class="text-brand-600 size-4 shrink-0" />
			<span class="min-w-0">
				<span class="block text-xs font-medium tracking-wide text-slate-500 uppercase"
					>Previous</span
				>
				<span class="text-brand-800 block truncate text-sm font-medium">{previousLabel ?? ''}</span>
			</span>
		</button>

		<p class="hidden text-center text-xs text-slate-500 lg:block">
			<kbd class="kbd">Space</kbd> play/pause · <kbd class="kbd">←</kbd><kbd class="kbd">→</kbd>
			10s · <kbd class="kbd">F</kbd> full screen
		</p>

		<button
			class="{navButton} justify-end text-right"
			onclick={onnext}
			disabled={!hasNext}
			title="Next lesson (Shift+N)"
		>
			<span class="min-w-0">
				<span class="block text-xs font-medium tracking-wide text-slate-500 uppercase">Next</span>
				<span class="text-brand-800 block truncate text-sm font-medium">{nextLabel ?? ''}</span>
			</span>
			<Icon name="next" class="text-brand-600 size-4 shrink-0" />
		</button>
	</footer>
</dialog>

<style>
	.seek {
		appearance: none;
		height: 0.375rem;
		border-radius: 9999px;
		cursor: pointer;
		background: linear-gradient(
			to right,
			var(--color-brand-400) var(--progress),
			rgb(255 255 255 / 0.45) var(--progress),
			rgb(255 255 255 / 0.45) var(--buffered),
			rgb(255 255 255 / 0.2) var(--buffered)
		);
	}
	.seek:disabled {
		cursor: default;
	}
	.seek::-webkit-slider-thumb {
		appearance: none;
		width: 0.9rem;
		height: 0.9rem;
		border-radius: 9999px;
		background: white;
		box-shadow: 0 0 0 3px var(--color-brand-400);
	}
	.seek::-moz-range-thumb {
		width: 0.9rem;
		height: 0.9rem;
		border: 0;
		border-radius: 9999px;
		background: white;
		box-shadow: 0 0 0 3px var(--color-brand-400);
	}

	.volume {
		accent-color: var(--color-brand-300);
	}

	.kbd {
		display: inline-block;
		min-width: 1.4em;
		margin: 0 0.1em;
		padding: 0 0.3em;
		border-radius: 0.25rem;
		background: white;
		box-shadow: inset 0 -1px 0 var(--color-brand-200);
		border: 1px solid var(--color-brand-100);
		color: var(--color-brand-800);
		font-family: inherit;
		font-size: 0.7rem;
	}
</style>
