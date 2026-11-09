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

	const labelAt = (i: number) => {
		const item = playlist[i];
		return item ? `${item.lesson.number} · ${item.lesson.title}` : undefined;
	};

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
				<article class="ring-brand-100 rounded-xl bg-white p-6 shadow-sm ring-1">
					<h2 class="border-brand-100 text-brand-800 mb-4 border-b pb-3 text-xl font-bold">
						Chapter {chapter.number}: {chapter.title}
					</h2>
					<ol class="flex flex-col gap-2">
						{#each chapter.lessons as lesson (lesson.id)}
							<li>
								<button
									type="button"
									class="group bg-brand-50 hover:bg-brand-100 focus-visible:ring-brand-500 flex w-full items-center gap-3 rounded-lg px-4 py-3 text-left transition focus-visible:ring-2 focus-visible:outline-none"
									onclick={() => play(lesson)}
								>
									<span
										class="text-brand-600 group-hover:bg-brand-600 grid size-9 shrink-0 place-items-center rounded-full bg-white shadow-sm group-hover:text-white"
									>
										<Icon name="play" class="size-4 translate-x-px" />
									</span>
									<span class="min-w-0 flex-1">
										<span class="block font-semibold text-slate-800">
											<span class="text-brand-600">Lesson {lesson.number}:</span>
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
	previousLabel={labelAt(index - 1)}
	nextLabel={labelAt(index + 1)}
	onclose={() => (currentId = null)}
	onprevious={() => step(-1)}
	onnext={() => step(1)}
/>
