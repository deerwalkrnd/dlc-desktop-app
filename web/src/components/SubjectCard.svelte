<script lang="ts">
	import { navigating, page } from '$app/state';
	import { onMount } from 'svelte';
	import LectureCard from './LectureCard.svelte';
	import { ProductionUrl } from '$lib/constant';

	let { subject } = $props();
	let currentPath = page.url.pathname;
</script>

<section class="rounded-lg bg-white p-6 shadow-lg">
	<a href={`${ProductionUrl}${currentPath}/subject/${subject.Name}-${subject.ID}`}>
		<h2 class="mb-4 border-slate-200 pb-2 text-2xl font-semibold text-indigo-800">
			{subject.Name}
		</h2>
	</a>

	<!-- <div class="space-y-6">
		{#each subject.lectures as lecture}
			<LectureCard {lecture} />
		{/each}
	</div> -->
</section>
