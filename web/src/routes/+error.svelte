<script lang="ts">
	import { page } from '$app/state';
	import notFoundAnimation from '$lib/assets/animation2.webm';

	const notFound = $derived(page.status === 404);
	const offline = $derived(page.status === 503);
</script>

<svelte:head>
	<title>{notFound ? 'Not found' : 'Something went wrong'} · Deerwalk Learning Center</title>
</svelte:head>

<div class="flex min-h-[70vh] flex-col items-center justify-center gap-5 px-6 text-center">
	{#if notFound}
		<video src={notFoundAnimation} loop autoplay muted playsinline class="max-h-64">
			<track kind="captions" />
		</video>
		<h1 class="text-3xl font-bold text-slate-800">We couldn't find that page</h1>
	{:else if offline}
		<h1 class="text-3xl font-bold text-slate-800">The DLC app is not responding</h1>
	{:else}
		<h1 class="text-3xl font-bold text-slate-800">Something went wrong</h1>
	{/if}

	{#if page.error?.message}
		<p class="max-w-lg text-slate-600">{page.error.message}</p>
	{/if}

	<div class="flex gap-3">
		<a href="/" class="bg-brand-600 hover:bg-brand-700 rounded-md px-4 py-2 font-medium text-white"
			>Go to classes</a
		>
		{#if !notFound}
			<button
				class="rounded-md border border-slate-300 bg-white px-4 py-2 font-medium text-slate-700 hover:bg-slate-50"
				onclick={() => location.reload()}>Try again</button
			>
		{/if}
	</div>
</div>
