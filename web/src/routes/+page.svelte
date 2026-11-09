<script lang="ts">
	import EmptyState from '$lib/components/EmptyState.svelte';
	import PageHeader from '$lib/components/PageHeader.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// Full class names so Tailwind includes them in the build.
	const colors = [
		'from-red-500 to-red-600',
		'from-blue-500 to-blue-600',
		'from-green-500 to-green-600',
		'from-amber-500 to-amber-600',
		'from-purple-500 to-purple-600',
		'from-pink-500 to-pink-600',
		'from-indigo-500 to-indigo-600',
		'from-teal-500 to-teal-600',
		'from-orange-500 to-orange-600',
		'from-cyan-500 to-cyan-600',
		'from-lime-600 to-lime-700',
		'from-rose-500 to-rose-600'
	];
	const colorFor = (n: number) =>
		colors[(((n - 1) % colors.length) + colors.length) % colors.length];
</script>

<PageHeader title="Choose your class" />

<div class="mx-auto max-w-7xl px-6 py-10">
	{#if data.classes.length}
		<ul class="grid grid-cols-2 gap-6 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
			{#each data.classes as c (c.number)}
				<li>
					<a
						href="/class/{c.number}"
						class="block overflow-hidden rounded-xl bg-white shadow-md transition hover:-translate-y-1 hover:shadow-xl focus-visible:-translate-y-1"
					>
						<div
							class="flex h-28 items-center justify-center bg-gradient-to-br {colorFor(c.number)}"
						>
							<span class="text-4xl font-bold text-white">Class {c.number}</span>
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
				<code class="rounded bg-slate-100 px-1.5 py-0.5 text-sm break-all text-slate-700"
					>{data.library.path}</code
				>
				and then rescan the library.
			</p>
			<a
				href="/library"
				class="mt-4 inline-block rounded-md bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700"
				>Open library settings</a
			>
		</EmptyState>
	{/if}
</div>
