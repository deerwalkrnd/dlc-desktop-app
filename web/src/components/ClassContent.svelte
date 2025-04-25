<script lang="ts">
	import SubjectCard from './SubjectCard.svelte';
	import EmptyState from './EmptyState.svelte';

	export let classData:
		| {
				subjects?: Array<any>;
		  }
		| undefined;
</script>

{#if classData?.subjects?.length}
	<div class="grid gap-8">
		{#each classData.subjects as subject}
			<SubjectCard {subject} />
		{/each}
	</div>
{:else}
	<EmptyState message="No class data available" />
{/if}
