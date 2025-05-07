<script lang="ts">
	import LectureCard from '$components/LectureCard.svelte';
	import Sidebar from '$components/Sidebar.svelte';
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<div class="flex flex-row">
	<Sidebar subjects={data.subjectList} />
	<div class="flex flex-1 flex-col">
		<h1 class="mb-2 p-10 text-center text-3xl font-bold text-blue-800">
			{data.subjectName.charAt(0).toUpperCase() + data.subjectName.slice(1)}
		</h1>
		<LectureCard lectures={data.subjectData.lectures} />
	</div>
</div>
