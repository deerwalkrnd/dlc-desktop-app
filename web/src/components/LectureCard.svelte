<script lang="ts">
	import LessonItem from './LessonItem.svelte';

	export let lecture: {
		number: number;
		name: string;
		lessons?: Array<any>;
	};
</script>

<div>
	<h3>
		Lecture {lecture.number} : {lecture.name}
	</h3>

	{#if lecture.lessons?.length}
		<ul>
			{#each lecture.lessons as lesson}
				<LessonItem {lesson} />
			{/each}
		</ul>
	{:else}
		<p class="italic text-slate-500">No lessons available</p>
	{/if}
</div>
