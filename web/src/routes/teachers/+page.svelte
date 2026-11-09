<script lang="ts">
	import EmptyState from '$lib/components/EmptyState.svelte';
	import PageHeader from '$lib/components/PageHeader.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Teachers · Deerwalk Learning Center</title>
</svelte:head>

<PageHeader title="Our instructors" />

<div class="mx-auto max-w-7xl px-6 py-10">
	{#if data.teachers.length}
		<ul class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
			{#each data.teachers as teacher (teacher.name)}
				<li class="overflow-hidden rounded-xl bg-white shadow-md">
					<div class="flex h-36 items-center justify-center bg-blue-100">
						<div
							class="grid size-24 place-items-center rounded-full bg-blue-500 text-3xl font-bold text-white shadow-md"
						>
							{teacher.name.charAt(0)}
						</div>
					</div>
					<div class="border-t border-blue-100 p-4 text-center">
						<h2 class="text-lg font-medium text-blue-800">{teacher.name}</h2>
						<p class="mt-1 text-sm text-slate-500">
							{teacher.lessonCount}
							{teacher.lessonCount === 1 ? 'lesson' : 'lessons'}
							{#if teacher.subjects.length}· {teacher.subjects.join(', ')}{/if}
						</p>
					</div>
				</li>
			{/each}
		</ul>
	{:else}
		<EmptyState title="No teachers yet" />
	{/if}
</div>
