<script lang="ts">
	import LessonItem from './LessonItem.svelte';

	export let lecture: {
		number: number;
		name: string;
		lessons?: Array<any>;
	};
</script>

<div class="rounded-lg bg-slate-50 p-4">
	<h3 class="mb-3 text-xl font-medium text-slate-800">
		Lecture {lecture.number}: {lecture.name}
	</h3>

	{#if lecture.lessons?.length}
		<ul class="space-y-2">
			{#each lecture.lessons as lesson}
				<LessonItem {lesson} />
			{/each}
		</ul>
	{:else}
		<p class="italic text-slate-500">No lessons available</p>
	{/if}
</div>
