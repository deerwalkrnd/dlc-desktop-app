<script lang="ts">
	import LectureCard from '$components/LectureCard.svelte';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
	let subjectName = data.subjectName;
</script>

<div>
	tesz
	<!-- <LectureCard /> -->
</div>
