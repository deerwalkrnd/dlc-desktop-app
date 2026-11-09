<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { api, ApiError, type Issue, type LibraryStatus } from '$lib/api';
	import Icon from '$lib/components/Icon.svelte';
	import PageHeader from '$lib/components/PageHeader.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// Rescan request lifecycle.
	type Rescan =
		| { state: 'idle' }
		| { state: 'scanning' }
		| { state: 'done'; added: number }
		| { state: 'failed'; message: string };
	let rescan = $state<Rescan>({ state: 'idle' });

	const lib: LibraryStatus = $derived(data.library);

	async function runRescan() {
		const before = lib.lessonCount;
		rescan = { state: 'scanning' };
		try {
			const after = await api.rescan(fetch);
			await invalidateAll();
			rescan = { state: 'done', added: after.lessonCount - before };
		} catch (e) {
			rescan = { state: 'failed', message: e instanceof ApiError ? e.message : String(e) };
		}
	}

	const scannedAt = $derived(new Date(lib.scannedAt).toLocaleString());
	const plural = (n: number, word: string) => `${n} ${word}${n === 1 ? '' : 's'}`;
</script>

<svelte:head>
	<title>Library · Deerwalk Learning Center</title>
</svelte:head>

{#snippet issueList(title: string, issues: Issue[], tone: 'red' | 'amber', help: string)}
	<section class="ring-brand-100 rounded-xl bg-white shadow-sm ring-1">
		<header class="border-brand-100 border-b px-6 py-4">
			<h2 class="text-lg font-semibold {tone === 'red' ? 'text-red-700' : 'text-amber-700'}">
				{title} ({issues.length})
			</h2>
			<p class="text-sm text-slate-500">{help}</p>
		</header>
		<ul class="divide-brand-100 divide-y">
			{#each issues as issue (issue.path + issue.message)}
				<li class="px-6 py-3">
					<p class="font-mono text-sm break-all text-slate-800">{issue.path}</p>
					<p class="text-sm {tone === 'red' ? 'text-red-600' : 'text-amber-700'}">
						{issue.message}
					</p>
				</li>
			{/each}
		</ul>
	</section>
{/snippet}

<PageHeader title="Library" subtitle="Where the app finds its videos, and what it found." />

<div class="mx-auto flex max-w-5xl flex-col gap-6 px-6 py-10">
	<section
		class="ring-brand-100 flex flex-col gap-4 rounded-xl bg-white p-6 shadow-sm ring-1 sm:flex-row sm:items-center"
	>
		<div class="min-w-0 flex-1">
			<p class="text-sm font-medium text-slate-500">Video folder</p>
			<p class="font-mono break-all text-slate-900">{lib.path}</p>
			{#if !lib.exists}
				<p class="mt-1 text-sm text-red-600">
					This folder does not exist yet. Create it and copy your videos into it.
				</p>
			{/if}
			<p class="mt-2 text-sm text-slate-500">
				{plural(lib.lessonCount, 'lesson')} in {plural(lib.classCount, 'class')} · last scanned {scannedAt}
				({lib.scanMillis} ms) · version {lib.version}
			</p>
		</div>
		<button
			class="bg-brand-600 hover:bg-brand-700 flex shrink-0 items-center justify-center gap-2 rounded-md px-4 py-2 font-medium text-white disabled:opacity-60"
			onclick={runRescan}
			disabled={rescan.state === 'scanning'}
		>
			<Icon name="refresh" class="size-4 {rescan.state === 'scanning' ? 'animate-spin' : ''}" />
			{rescan.state === 'scanning' ? 'Scanning…' : 'Rescan library'}
		</button>
	</section>

	<div aria-live="polite">
		{#if rescan.state === 'done'}
			<p class="rounded-md bg-green-50 px-4 py-3 text-green-800">
				Scan complete.
				{#if rescan.added > 0}{plural(rescan.added, 'new lesson')} found.
				{:else if rescan.added < 0}{plural(-rescan.added, 'lesson')} removed.
				{:else}No changes.{/if}
			</p>
		{:else if rescan.state === 'failed'}
			<p class="rounded-md bg-red-50 px-4 py-3 text-red-800">Scan failed: {rescan.message}</p>
		{/if}
	</div>

	{#if lib.skipped.length}
		{@render issueList(
			'Skipped files',
			lib.skipped,
			'red',
			'These videos are not shown in the app because their names do not follow the naming format.'
		)}
	{/if}
	{#if lib.warnings.length}
		{@render issueList(
			'Warnings',
			lib.warnings,
			'amber',
			'These videos are shown, but their names look like they contain a mistake.'
		)}
	{/if}

	<section class="bg-brand-50 rounded-xl p-6 text-sm text-slate-700">
		<h2 class="mb-2 font-semibold text-slate-900">Naming format</h2>
		<p class="font-mono break-words">
			TEACHER - LESSON NO - CHAPTER TITLE - LESSON TITLE - SUBJECT - OLD|NEW - CLASS.mp4
		</p>
		<p class="mt-2 font-mono break-words text-slate-500">
			Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics -
			OLD - 8.mp4
		</p>
	</section>
</div>
