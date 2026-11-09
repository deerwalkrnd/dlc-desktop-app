<script lang="ts">
	import { goto } from '$app/navigation';
	import { lessonHref, syllabusLabel } from '#lib/api.ts';
	import EmptyState from '#lib/components/EmptyState.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import { highlight } from '#lib/highlight.ts';
	import { formatTime } from '#lib/player.ts';
	import { progress } from '#lib/progress.svelte.ts';
	import { resumePosition } from '#lib/progress.ts';
	import { untrack } from 'svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let query = $state('');
	// Follow the URL (first load, back/forward) without clobbering typing.
	$effect(() => {
		const q = data.q;
		untrack(() => {
			if (q !== query) query = q;
		});
	});

	// Search as you type, keeping the URL shareable and Back-friendly.
	let timer: ReturnType<typeof setTimeout> | undefined;
	function onInput() {
		clearTimeout(timer);
		timer = setTimeout(() => {
			const url = query.trim() ? `/search?q=${encodeURIComponent(query)}` : '/search';
			goto(url, { replace: true, reset: false });
		}, 200);
	}
	$effect(() => () => clearTimeout(timer));
</script>

<svelte:head>
	<title>{data.q ? `${data.q} · Search` : 'Search'} · Deerwalk Learning Center</title>
</svelte:head>

{#snippet marked(text: string)}
	{#each highlight(text, data.q) as seg, i (i)}
		{#if seg.match}<mark class="bg-brand-100 text-brand-900 rounded-sm px-0.5">{seg.text}</mark
			>{:else}{seg.text}{/if}
	{/each}
{/snippet}

<header class="border-brand-100 border-b bg-white">
	<div class="mx-auto max-w-4xl px-6 py-8">
		<h1 class="text-brand-800 mb-4 text-3xl font-bold">Search lessons</h1>
		<form
			role="search"
			onsubmit={(e) => {
				e.preventDefault();
				clearTimeout(timer);
				goto(`/search?q=${encodeURIComponent(query)}`, { replace: true, reset: false });
			}}
		>
			<label class="relative block">
				<span class="sr-only">Search by lesson, chapter, subject, teacher or class</span>
				<Icon
					name="search"
					class="text-brand-400 pointer-events-none absolute top-1/2 left-4 size-5 -translate-y-1/2"
				/>
				<!-- svelte-ignore a11y_autofocus -->
				<input
					type="search"
					bind:value={query}
					oninput={onInput}
					data-search-input
					autofocus
					autocomplete="off"
					maxlength="200"
					placeholder="Try “polynomials”, “class 8 science” or a teacher's name"
					class="border-brand-200 bg-brand-50 focus:border-brand-400 focus:ring-brand-100 w-full rounded-xl border py-3 pr-4 pl-12 text-lg text-slate-900 placeholder:text-slate-400 focus:bg-white focus:ring-4 focus:outline-none"
				/>
			</label>
		</form>
	</div>
</header>

<div class="mx-auto max-w-4xl px-6 py-8">
	{#if data.search === null}
		<p class="text-center text-slate-500">
			Search across every class by lesson title, chapter, subject, teacher, class or lesson number.
		</p>
	{:else if data.search.results.length === 0}
		<EmptyState title="No lessons match “{data.q}”" icon="search">
			<p>Check the spelling, or try fewer or shorter words.</p>
		</EmptyState>
	{:else}
		<p class="mb-4 text-sm text-slate-500" aria-live="polite">
			{data.search.total}
			{data.search.total === 1 ? 'lesson' : 'lessons'}
			{#if data.search.total > data.search.results.length}
				· showing the best {data.search.results.length}
			{/if}
		</p>
		<ul class="flex flex-col gap-2">
			{#each data.search.results as r (r.lesson.id)}
				{@const p = progress.get(r.lesson.id)}
				{@const resumeAt = resumePosition(p)}
				<li>
					<a
						href={lessonHref(r.class, r.subjectSlug, r.lesson.id)}
						class="group ring-brand-100 hover:bg-brand-50 hover:ring-brand-300 flex items-center gap-4 rounded-xl bg-white px-4 py-3 shadow-sm ring-1 transition"
					>
						<span
							class="grid size-10 shrink-0 place-items-center rounded-full transition {p?.completed
								? 'bg-brand-600 text-white'
								: 'bg-brand-50 text-brand-600 group-hover:bg-brand-600 group-hover:text-white'}"
						>
							<Icon
								name={p?.completed ? 'check' : 'play'}
								class="size-4 {p?.completed ? '' : 'translate-x-px'}"
							/>
						</span>
						<span class="min-w-0 flex-1">
							<span class="block font-semibold text-slate-800">
								<span class="text-brand-600">Lesson {r.lesson.number}:</span>
								{@render marked(r.lesson.title)}
							</span>
							<span class="block truncate text-sm text-slate-500">
								Class {r.class} · {@render marked(r.subjectName)} ({syllabusLabel(
									r.syllabus
								).toLowerCase()}) · Chapter {r.chapterNumber}: {@render marked(r.chapterTitle)}
							</span>
							<span class="block text-sm text-slate-500">{@render marked(r.lesson.teacher)}</span>
						</span>
						{#if p?.completed}
							<span class="text-brand-700 shrink-0 text-sm font-medium">Watched</span>
						{:else if resumeAt !== null}
							<span class="text-brand-700 shrink-0 text-sm font-medium"
								>Resume at {formatTime(resumeAt)}</span
							>
						{/if}
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</div>
