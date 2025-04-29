<script lang="ts">
	import LectureCard from '$components/LectureCard.svelte';
	import Sidebar from '$components/Sidebar.svelte';
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<div class="flex flex-row">
	<Sidebar subjects={data.subjectList} />
	<div class="flex-1">
		<LectureCard lectures={data.subjectData.lectures} lessons={data.lessons} />
	</div>
</div>
