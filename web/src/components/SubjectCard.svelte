<script lang="ts">
	import { navigating, page } from '$app/state';
	import { onMount } from 'svelte';
	import LectureCard from './LectureCard.svelte';
	import { ProductionUrl } from '$lib/constant';

	export let subject: {
		name: string;
		lectures: Array<any>;
	};
	let currentPath = page.url.pathname;
	onMount(() => {
		console.log(ProductionUrl);
	});
</script>

<section class="rounded-lg bg-white p-6 shadow-lg">
	<a href={`${ProductionUrl}${currentPath}/subject/${subject.name}`}>
		<h2 class="mb-4 border-slate-200 pb-2 text-2xl font-semibold text-indigo-800">
			{subject.name}
		</h2>
	</a>

	<!-- <div class="space-y-6">
		{#each subject.lectures as lecture}
			<LectureCard {lecture} />
		{/each}
	</div> -->
</section>
