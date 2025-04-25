<script lang="ts">
	import LectureCard from './LectureCard.svelte';

	export let subject: {
		name: string;
		lectures: Array<any>;
	};
</script>

<section class="rounded-lg bg-white p-6 shadow-md">
	<h2 class="mb-4 border-b border-slate-200 pb-2 text-2xl font-semibold text-indigo-800">
		{subject.name}
	</h2>

	<div class="space-y-6">
		{#each subject.lectures as lecture}
			<LectureCard {lecture} />
		{/each}
	</div>
</section>
