<script lang="ts">
	import { lectures } from '$lib/data';

	export let lesson: {
		number: number;
		name: string;
		teacher?: {
			name: string;
		};
	};
</script>

{#snippet Lecture(number: any)}
	<li class="flex items-center rounded p-2 transition hover:bg-slate-100">
		<span class="mr-2 font-medium text-indigo-600">{lesson.number}.</span>
		<span class="flex-1">{lesson.name}</span>
		{#if lesson.teacher?.name}
			<span class="text-sm text-slate-600">by {lesson.teacher.name}</span>
		{/if}
	</li>
{/snippet}

<a href={lesson.name}> {@render Lecture(lesson.number)} </a>
