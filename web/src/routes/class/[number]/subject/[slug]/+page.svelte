<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		lessonHref,
		lessonSubtitle,
		lessonTitle,
		subjectHref,
		syllabusLabel,
		type Lesson
	} from '#lib/api.ts';
	import EmptyState from '#lib/components/EmptyState.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import VideoPlayer, { type PlayableLesson } from '#lib/components/VideoPlayer.svelte';
	import { formatTime } from '#lib/player.ts';
	import { progress } from '#lib/progress.svelte.ts';
	import { resumePosition, watchedFraction } from '#lib/progress.ts';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const siblings = $derived(data.cls.subjects.filter((s) => s.syllabus === data.subject.syllabus));

	// Lessons in playback order, for previous/next.
	const playlist = $derived(
		data.subject.chapters.flatMap((ch) => ch.lessons.map((lesson) => ({ lesson, chapter: ch })))
	);

	// The playing lesson lives in the URL (?lesson=ID), so search results and
	// "Continue watching" can link straight to it, Back closes the player and
	// a reload keeps it open.
	const currentId = $derived(page.url.searchParams.get('lesson'));
	const index = $derived(playlist.findIndex((p) => p.lesson.id === currentId));
	const playing = $derived.by((): PlayableLesson | null => {
		const item = playlist[index];
		if (!item) return null;
		return {
			id: item.lesson.id,
			mediaUrl: item.lesson.mediaUrl,
			title: lessonTitle(item.lesson),
			subtitle: lessonSubtitle(item.chapter, item.lesson.teacher, data.subject.name),
			href: lessonHref(data.subject.class, data.subject.slug, item.lesson.id)
		};
	});

	// True when this page opened the player (so closing can go Back instead
	// of leaving a duplicate history entry).
	let openedHere = false;

	function urlWithLesson(id: string | null): URL {
		const url = new URL(page.url.href);
		if (id) url.searchParams.set('lesson', id);
		else url.searchParams.delete('lesson');
		return url;
	}

	function play(lesson: Lesson) {
		openedHere = true;
		goto(urlWithLesson(lesson.id), { reset: false });
	}

	function step(delta: number) {
		const next = playlist[index + delta];
		if (next) goto(urlWithLesson(next.lesson.id), { replace: true, reset: false });
	}

	function close() {
		if (openedHere) {
			openedHere = false;
			history.back();
		} else {
			goto(urlWithLesson(null), { replace: true, reset: false });
		}
	}

	const labelAt = (i: number) => {
		const item = playlist[i];
		return item ? `${item.lesson.number} · ${item.lesson.title}` : undefined;
	};

	const watchedIn = (lessons: Lesson[]) =>
		lessons.filter((l) => progress.get(l.id)?.completed).length;
</script>

<svelte:head>
	<title>{data.subject.name} · Class {data.subject.class} · Deerwalk Learning Center</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-col gap-6 px-6 py-8 md:flex-row">
	<aside class="md:w-64 md:shrink-0">
		<nav
			class="ring-brand-100 overflow-hidden rounded-xl bg-white shadow-sm ring-1 md:sticky md:top-24"
		>
			<div class="border-brand-100 bg-brand-50 border-b px-4 py-3">
				<a
					href="/class/{data.cls.number}?syllabus={data.subject.syllabus}"
					class="text-brand-700 text-sm hover:underline">← Class {data.cls.number}</a
				>
				<h2 class="text-brand-900 font-semibold">{syllabusLabel(data.subject.syllabus)}</h2>
			</div>
			<ul class="p-2">
				{#each siblings as s (s.slug)}
					{@const active = s.slug === data.subject.slug}
					<li>
						<a
							href={subjectHref(data.cls.number, s.slug)}
							aria-current={active ? 'page' : undefined}
							class="flex items-center gap-3 rounded-md px-3 py-2 transition {active
								? 'bg-brand-100 text-brand-800 font-medium'
								: 'hover:bg-brand-50 text-slate-700'}"
						>
							<span
								class="bg-brand-100 text-brand-700 grid size-8 shrink-0 place-items-center rounded-full text-sm font-semibold"
								>{s.name.charAt(0)}</span
							>
							<span class="truncate">{s.name}</span>
						</a>
					</li>
				{/each}
			</ul>
		</nav>
	</aside>

	<section class="min-w-0 flex-1">
		<h1 class="text-brand-800 mb-6 text-3xl font-bold">{data.subject.name}</h1>

		{#if data.subject.chapters.length === 0}
			<EmptyState title="No lessons in this subject yet" />
		{/if}

		<div class="flex flex-col gap-6">
			{#each data.subject.chapters as chapter (`${chapter.number}:${chapter.title}`)}
				{@const watched = watchedIn(chapter.lessons)}
				<article class="ring-brand-100 rounded-xl bg-white p-6 shadow-sm ring-1">
					<header
						class="border-brand-100 mb-4 flex items-baseline justify-between gap-4 border-b pb-3"
					>
						<h2 class="text-brand-800 text-xl font-bold">
							Chapter {chapter.number}: {chapter.title}
						</h2>
						{#if watched > 0}
							<span class="shrink-0 text-sm text-slate-500">
								{watched} of {chapter.lessons.length} watched
							</span>
						{/if}
					</header>
					<ol class="flex flex-col gap-2">
						{#each chapter.lessons as lesson (lesson.id)}
							{@const p = progress.get(lesson.id)}
							{@const done = p?.completed ?? false}
							{@const resumeAt = resumePosition(p)}
							<li>
								<button
									type="button"
									class="group bg-brand-50 hover:bg-brand-100 relative flex w-full items-center gap-3 overflow-hidden rounded-lg px-4 py-3 text-left transition"
									onclick={() => play(lesson)}
								>
									<span
										class="grid size-9 shrink-0 place-items-center rounded-full shadow-sm transition {done
											? 'bg-brand-600 text-white'
											: 'text-brand-600 group-hover:bg-brand-600 bg-white group-hover:text-white'}"
									>
										<Icon
											name={done ? 'check' : 'play'}
											class="size-4 {done ? '' : 'translate-x-px'}"
										/>
									</span>
									<span class="min-w-0 flex-1">
										<span class="block font-semibold text-slate-800">
											<span class="text-brand-600">Lesson {lesson.number}:</span>
											{lesson.title}
										</span>
										<span class="block text-sm text-slate-500">{lesson.teacher}</span>
									</span>
									{#if done}
										<span class="text-brand-700 shrink-0 text-sm font-medium">Watched</span>
									{:else if resumeAt !== null}
										<span class="text-brand-700 shrink-0 text-sm font-medium">
											Resume at {formatTime(resumeAt)}
										</span>
									{/if}
									{#if resumeAt !== null}
										<span class="bg-brand-100 absolute inset-x-0 bottom-0 h-1" aria-hidden="true">
											<span
												class="bg-brand-400 block h-full"
												style="width: {watchedFraction(p) * 100}%"
											></span>
										</span>
									{/if}
								</button>
							</li>
						{/each}
					</ol>
				</article>
			{/each}
		</div>
	</section>
</div>

<VideoPlayer
	lesson={playing}
	hasPrevious={index > 0}
	hasNext={index >= 0 && index < playlist.length - 1}
	previousLabel={labelAt(index - 1)}
	nextLabel={labelAt(index + 1)}
	onclose={close}
	onprevious={() => step(-1)}
	onnext={() => step(1)}
/>
