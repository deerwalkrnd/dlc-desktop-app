<script lang="ts">
	import ClassContent from '$components/ClassContent.svelte';
	import { onMount } from 'svelte';
	import Header from '../../../components/Header.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="flex min-h-screen flex-col bg-slate-50">
	<Header classNumber={data.getClass.Number} />

	<div class="flex flex-1 flex-col md:flex-row">
		<main class="flex-1 p-6">
			<ClassContent subjectData={data.subjects} />
		</main>
	</div>
</div>
