<script lang="ts">
	import { page } from '$app/state';
	import { subjectHref, syllabusLabel, type Syllabus } from '$lib/api';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import PageHeader from '$lib/components/PageHeader.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// The selected syllabus lives in the URL (?syllabus=new) so it survives
	// reloads and the back button. Default to the first one with content.
	const syllabus = $derived.by((): Syllabus | undefined => {
		const requested = page.url.searchParams.get('syllabus');
		return data.cls.syllabi.find((s) => s === requested) ?? data.cls.syllabi[0];
	});
	const subjects = $derived(data.cls.subjects.filter((s) => s.syllabus === syllabus));
</script>

<svelte:head>
	<title>Class {data.cls.number} · Deerwalk Learning Center</title>
</svelte:head>

<PageHeader title="Class {data.cls.number}">
	{#if data.cls.syllabi.length > 1}
		<div class="bg-brand-50 inline-flex rounded-lg p-1" role="tablist" aria-label="Syllabus">
			{#each data.cls.syllabi as s (s)}
				<a
					href="?syllabus={s}"
					role="tab"
					aria-selected={s === syllabus}
					data-sveltekit-replacestate
					data-sveltekit-noscroll
					class="rounded-md px-4 py-1.5 text-sm font-medium transition {s === syllabus
						? 'text-brand-700 bg-white shadow-sm'
						: 'text-slate-600 hover:text-slate-900'}">{syllabusLabel(s)}</a
				>
			{/each}
		</div>
	{:else if syllabus}
		<p class="text-sm text-slate-500">{syllabusLabel(syllabus)}</p>
	{/if}
</PageHeader>

<div class="mx-auto max-w-7xl px-6 py-10">
	{#if subjects.length}
		<ul class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each subjects as subject (subject.slug)}
				<li>
					<a
						href={subjectHref(data.cls.number, subject.slug)}
						class="ring-brand-100 hover:ring-brand-300 flex h-full flex-col gap-2 rounded-xl bg-white p-6 shadow-md ring-1 transition hover:-translate-y-1 hover:shadow-xl"
					>
						<h2 class="text-brand-800 text-2xl font-semibold">{subject.name}</h2>
						<p class="text-sm text-slate-500">
							{subject.chapterCount}
							{subject.chapterCount === 1 ? 'chapter' : 'chapters'} · {subject.lessonCount}
							{subject.lessonCount === 1 ? 'lesson' : 'lessons'}
						</p>
					</a>
				</li>
			{/each}
		</ul>
	{:else}
		<EmptyState title="No subjects for this class yet" />
	{/if}
</div>
