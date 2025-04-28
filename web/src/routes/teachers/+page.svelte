<script lang="ts">
	import Header from '$components/Header.svelte';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let teachers = data.data.teachers;
	onMount(() => {
		console.log(teachers);
	});
</script>

{#snippet Teacher(name: string)}
	<li
		class="flex items-center justify-center rounded-md bg-blue-400 p-10 py-20 text-2xl font-semibold text-white ring ring-2 ring-blue-600"
	>
		<span>
			{name}
		</span>
	</li>
{/snippet}

<div class="flex flex-col gap-10">
	<div class="bg-blue-700 p-10">
		<h1 class="text-3xl font-semibold text-white">Instructors</h1>
	</div>
	<div class="xlgrid-cols-4 grid grid-cols-3 gap-10 p-10">
		{#each teachers as teacher}
			{@render Teacher(teacher.Name)}
		{/each}
	</div>
</div>
