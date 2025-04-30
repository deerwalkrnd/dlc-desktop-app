<script lang="ts">
	import Header from '$components/Header.svelte';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let teachers = data.data.teachers;
</script>

{#snippet Teacher(name: string)}
	<div class="group overflow-hidden rounded-lg bg-white shadow-md transition-all hover:shadow-lg">
		<div class="flex aspect-square items-center justify-center bg-blue-100">
			<div
				class="flex h-24 w-24 items-center justify-center rounded-full bg-blue-500 text-3xl font-bold text-white shadow-md"
			>
				{name.charAt(0)}
			</div>
		</div>
		<div class="border-t border-blue-100 bg-white p-4">
			<h3 class="text-center text-lg font-medium text-blue-800">{name}</h3>
		</div>
	</div>
{/snippet}

<div class="flex min-h-screen flex-col bg-blue-50">
	<div class="bg-white px-10 py-10 shadow-md">
		<div class="mx-auto max-w-6xl">
			<h1 class="mb-2 text-center text-3xl font-bold text-blue-800">Our Instructors</h1>
		</div>
	</div>

	<div class="mx-auto w-full max-w-6xl px-6 py-10">
		<div class="grid grid-cols-1 gap-8 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
			{#each teachers as teacher}
				{@render Teacher(teacher.Name)}
			{/each}
		</div>
	</div>
</div>
