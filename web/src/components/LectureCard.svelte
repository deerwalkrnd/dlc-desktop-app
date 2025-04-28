<script lang="ts">
	import LessonItem from './LessonItem.svelte';

	let { lectures } = $props();
	console.log(lectures);
</script>

{#snippet Lecture(lectureNumber: number, lectureName: string, lectureId: number)}
	<h3>
		Lecture {lectureNumber} : {lectureName}
	</h3>
{/snippet}

<div>
	{#each lectures as lecture}
		{@render Lecture(lecture.Number, lecture.Name, lecture.ID)}
	{/each}
	<!-- {#if lecture.lessons?.length}
		<ul>
			{#each lecture.lessons as lesson}
				<LessonItem {lesson} />
			{/each}
		</ul>
	{:else}
		<p class="italic text-slate-500">No lessons available</p>
	{/if} -->
</div>
