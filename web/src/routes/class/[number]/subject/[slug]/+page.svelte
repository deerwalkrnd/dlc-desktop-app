<script lang="ts">
	import { subjectHref, syllabusLabel, type Lesson } from '$lib/api';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import Icon from '$lib/components/Icon.svelte';
	import VideoPlayer, { type PlayableLesson } from '$lib/components/VideoPlayer.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const siblings = $derived(data.cls.subjects.filter((s) => s.syllabus === data.subject.syllabus));

	// Lessons in playback order, for previous/next.
	const playlist = $derived(
		data.subject.chapters.flatMap((ch) => ch.lessons.map((lesson) => ({ lesson, chapter: ch })))
	);

	let currentId = $state<string | null>(null);
	const index = $derived(playlist.findIndex((p) => p.lesson.id === currentId));
	const playing = $derived.by((): PlayableLesson | null => {
		const item = playlist[index];
		if (!item) return null;
		return {
			id: item.lesson.id,
			mediaUrl: item.lesson.mediaUrl,
			title: `Lesson ${item.lesson.number}: ${item.lesson.title}`,
			subtitle: `Chapter ${item.chapter.number}: ${item.chapter.title} · ${item.lesson.teacher}`
		};
	});

	function play(lesson: Lesson) {
		currentId = lesson.id;
	}

	function step(delta: number) {
		const next = playlist[index + delta];
		if (next) currentId = next.lesson.id;
	}
</script>

<svelte:head>
	<title>{data.subject.name} · Class {data.subject.class} · Deerwalk Learning Center</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-col gap-6 px-6 py-8 md:flex-row">
	<aside class="md:w-64 md:shrink-0">
		<nav
			class="overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-slate-200 md:sticky md:top-24"
		>
			<div class="border-b border-slate-100 bg-blue-50 px-4 py-3">
				<a
					href="/class/{data.cls.number}?syllabus={data.subject.syllabus}"
					class="text-sm text-blue-700 hover:underline">← Class {data.cls.number}</a
				>
				<h2 class="font-semibold text-blue-900">{syllabusLabel(data.subject.syllabus)}</h2>
			</div>
			<ul class="p-2">
				{#each siblings as s (s.slug)}
					{@const active = s.slug === data.subject.slug}
					<li>
						<a
							href={subjectHref(data.cls.number, s.slug)}
							aria-current={active ? 'page' : undefined}
							class="flex items-center gap-3 rounded-md px-3 py-2 transition {active
								? 'bg-blue-100 font-medium text-blue-800'
								: 'text-slate-700 hover:bg-slate-100'}"
						>
							<span
								class="grid size-8 shrink-0 place-items-center rounded-full bg-blue-100 text-sm font-semibold text-blue-700"
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
		<h1 class="mb-6 text-3xl font-bold text-blue-800">{data.subject.name}</h1>

		{#if data.subject.chapters.length === 0}
			<EmptyState title="No lessons in this subject yet" />
		{/if}

		<div class="flex flex-col gap-6">
			{#each data.subject.chapters as chapter (`${chapter.number}:${chapter.title}`)}
				<article class="rounded-xl bg-white p-6 shadow-sm ring-1 ring-slate-200">
					<h2 class="mb-4 border-b border-slate-100 pb-3 text-xl font-bold text-indigo-700">
						Chapter {chapter.number}: {chapter.title}
					</h2>
					<ol class="flex flex-col gap-2">
						{#each chapter.lessons as lesson (lesson.id)}
							<li>
								<button
									type="button"
									class="group flex w-full items-center gap-3 rounded-lg bg-indigo-50 px-4 py-3 text-left transition hover:bg-indigo-100 focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:outline-none"
									onclick={() => play(lesson)}
								>
									<span
										class="grid size-9 shrink-0 place-items-center rounded-full bg-white text-indigo-600 shadow-sm group-hover:bg-indigo-600 group-hover:text-white"
									>
										<Icon name="play" class="size-4 translate-x-px" />
									</span>
									<span class="min-w-0 flex-1">
										<span class="block font-semibold text-slate-800">
											<span class="text-indigo-600">Lesson {lesson.number}:</span>
											{lesson.title}
										</span>
										<span class="block text-sm text-slate-500">{lesson.teacher}</span>
									</span>
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
	onclose={() => (currentId = null)}
	onprevious={() => step(-1)}
	onnext={() => step(1)}
/>
