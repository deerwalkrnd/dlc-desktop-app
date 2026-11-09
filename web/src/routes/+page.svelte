<script lang="ts">
	import EmptyState from '$lib/components/EmptyState.svelte';
	import PageHeader from '$lib/components/PageHeader.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<PageHeader title="Choose your class" />

<div class="mx-auto max-w-7xl px-6 py-10">
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
