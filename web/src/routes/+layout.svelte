<script lang="ts">
	import { navigating } from '$app/state';
	import Navbar from '#lib/components/Navbar.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Deerwalk Learning Center</title>
</svelte:head>

{#if navigating.to}
	<div
		class="bg-brand-100 fixed inset-x-0 top-0 z-20 h-1 overflow-hidden"
		role="progressbar"
		aria-label="Loading"
	>
		<div class="bg-brand-600 h-full w-1/3 animate-[loading_1s_ease-in-out_infinite]"></div>
	</div>
{/if}

<div class="flex min-h-screen flex-col">
	<Navbar />
	<main class="flex-1">
		{@render children()}
	</main>
</div>

<style>
	@keyframes -global-loading {
		from {
			transform: translateX(-100%);
		}
		to {
			transform: translateX(300%);
		}
	}
</style>
