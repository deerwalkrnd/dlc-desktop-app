<script lang="ts">
	import LectureCard from '$components/LectureCard.svelte';
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<div>
	<LectureCard lectures={data.subjectData.lectures} lessons={data.lessons} />
</div>
