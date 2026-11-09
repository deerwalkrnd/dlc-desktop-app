<script lang="ts">
	import EmptyState from '#lib/components/EmptyState.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import { formatTime } from '#lib/player.ts';
	import { progress } from '#lib/progress.svelte.ts';
	import { continueWatching, watchedFraction } from '#lib/progress.ts';
	import PageHeader from '#lib/components/PageHeader.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const resume = $derived(continueWatching(progress.map, 4));
</script>

<PageHeader title="Choose your class" />

<div class="mx-auto max-w-7xl px-6 py-10">
	{#if resume.length && data.classes.length}
		<section class="mb-10" aria-labelledby="continue-heading">
			<h2 id="continue-heading" class="text-brand-800 mb-4 text-xl font-semibold">
				Continue watching
			</h2>
			<ul class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
				{#each resume as [id, p] (id)}
					<li
						class="group ring-brand-100 hover:ring-brand-300 relative overflow-hidden rounded-xl bg-white shadow-sm ring-1 transition hover:shadow-md"
					>
						<a href={p.href} class="flex items-start gap-3 p-4 pb-5">
							<span
								class="bg-brand-600 grid size-9 shrink-0 place-items-center rounded-full text-white shadow-sm"
							>
								<Icon name="play" class="size-4 translate-x-px" />
							</span>
							<span class="min-w-0 flex-1">
								<span class="line-clamp-2 font-semibold text-slate-800">{p.title}</span>
								<span class="mt-0.5 block truncate text-sm text-slate-500">{p.subtitle}</span>
								<span class="text-brand-700 mt-1 block text-sm font-medium">
									Resume at {formatTime(p.position)}
								</span>
							</span>
						</a>
						<button
							class="hover:bg-brand-50 hover:text-brand-700 absolute top-2 right-2 rounded-full p-1 text-slate-400 opacity-0 transition group-hover:opacity-100 focus-visible:opacity-100"
							onclick={() => progress.forget(id)}
							aria-label="Remove {p.title} from Continue watching"
							title="Remove"
						>
							<Icon name="close" class="size-4" />
						</button>
						<span class="bg-brand-100 absolute inset-x-0 bottom-0 h-1" aria-hidden="true">
							<span class="bg-brand-400 block h-full" style="width: {watchedFraction(p) * 100}%"
							></span>
						</span>
					</li>
				{/each}
			</ul>
		</section>
		<h2 class="text-brand-800 mb-4 text-xl font-semibold">All classes</h2>
	{/if}

	{#if data.classes.length}
		<ul class="grid grid-cols-2 gap-6 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
			{#each data.classes as c (c.number)}
				<li>
					<a
						href="/class/{c.number}"
						class="ring-brand-100 hover:ring-brand-300 block overflow-hidden rounded-xl bg-white shadow-sm ring-1 transition hover:-translate-y-1 hover:shadow-lg focus-visible:-translate-y-1"
					>
						<div
							class="from-brand-600 to-brand-400 flex h-28 items-center justify-center bg-linear-to-b"
						>
							<span class="text-4xl font-bold text-white drop-shadow-sm">Class {c.number}</span>
						</div>
						<div class="px-4 py-3 text-sm text-slate-600">
							{c.subjectCount}
							{c.subjectCount === 1 ? 'subject' : 'subjects'} · {c.lessonCount}
							{c.lessonCount === 1 ? 'lesson' : 'lessons'}
						</div>
					</a>
				</li>
			{/each}
		</ul>
	{:else}
		<EmptyState title="No lessons yet">
			<p>
				Copy your lesson videos into
				<code class="bg-brand-50 rounded px-1.5 py-0.5 text-sm break-all text-slate-700"
					>{data.library.path}</code
				>
				and then rescan the library.
			</p>
			<a
				href="/library"
				class="bg-brand-600 hover:bg-brand-700 mt-4 inline-block rounded-md px-4 py-2 font-medium text-white"
				>Open library settings</a
			>
		</EmptyState>
	{/if}
</div>
