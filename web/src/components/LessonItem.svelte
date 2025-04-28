<script lang="ts">
	import { page } from '$app/state';
	import { ProductionUrl } from '$lib/Consntant';
	import { lectures } from '$lib/data';
	import { onMount } from 'svelte';

	export let lesson: {
		number: string;
		name: string;
		teacher?: {
			name: string;
		};
	};
	$: currentPath = page.url.pathname;
	$: ProductionUrl;
	onMount(() => {
		$: console.log(currentPath, ProductionUrl);
	});
</script>

{#snippet Lecture(number: any)}
	<li class="flex items-center rounded p-2">
		<span class="mr-2 font-medium text-indigo-600">
			{lesson.number}.
		</span>
		<span class="flex-1">{lesson.name}</span>
		{#if lesson.teacher?.name}
			<span class="text-sm text-slate-600">
				{lesson.teacher.name}
			</span>
		{/if}
	</li>
{/snippet}

<a href={`${ProductionUrl}${currentPath}/${lesson.name}`}>
	{@render Lecture(lesson.number)}
</a>
