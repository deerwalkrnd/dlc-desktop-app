<script lang="ts">
	import EmptyState from './EmptyState.svelte';
	import SubjectCard from './SubjectCard.svelte';
	let { subjectData } = $props();
</script>

{#if subjectData?.subjects?.length}
	<div class="grid grid-cols-3 gap-8">
		{#each subjectData.subjects as subject}
			<SubjectCard {subject} />
		{/each}
	</div>
{:else}
	<EmptyState message="No class Data availabe" />
{/if}
