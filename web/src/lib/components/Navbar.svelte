<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import logo from '#lib/assets/dlc.png';
	import Icon from './Icon.svelte';

	const links = [
		{ href: '/', label: 'Classes', match: (p: string) => p === '/' || p.startsWith('/class') },
		{ href: '/teachers', label: 'Teachers', match: (p: string) => p.startsWith('/teachers') },
		{ href: '/library', label: 'Library', match: (p: string) => p.startsWith('/library') }
	];

	const onSearchPage = $derived(page.url.pathname === '/search');
	let query = $state('');

	function submit(e: SubmitEvent) {
		e.preventDefault();
		const q = query.trim();
		query = '';
		goto(q ? `/search?q=${encodeURIComponent(q)}` : '/search');
	}

	// Press "/" anywhere to jump to search (unless already typing).
	function onkeydown(e: KeyboardEvent) {
		if (e.key !== '/' || e.ctrlKey || e.metaKey || e.altKey) return;
		const t = e.target as HTMLElement;
		if (t.isContentEditable || ['INPUT', 'TEXTAREA', 'SELECT'].includes(t.tagName)) return;
		if (document.querySelector('dialog[open]')) return; // the video player
		const input = document.querySelector<HTMLInputElement>('[data-search-input]');
		if (!input) return;
		e.preventDefault();
		input.focus();
		input.select();
	}
</script>

<svelte:window {onkeydown} />

<nav class="border-brand-100 sticky top-0 z-10 border-b bg-white/95 backdrop-blur">
	<div class="mx-auto flex max-w-7xl items-center justify-between gap-6 px-6 py-3">
		<a href="/" class="flex shrink-0 items-center">
			<img alt="Deerwalk Learning Center" src={logo} width="160" height="50" class="h-10 w-auto" />
		</a>

		<div class="flex flex-1 items-center justify-end gap-4">
			{#if !onSearchPage}
				<form role="search" class="hidden max-w-xs flex-1 sm:block" onsubmit={submit}>
					<label class="relative block">
						<span class="sr-only">Search lessons</span>
						<Icon
							name="search"
							class="text-brand-400 pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2"
						/>
						<input
							type="search"
							bind:value={query}
							data-search-input
							autocomplete="off"
							maxlength="200"
							placeholder="Search lessons"
							class="border-brand-100 bg-brand-50 focus:border-brand-300 focus:ring-brand-100 w-full rounded-full border py-1.5 pr-10 pl-9 text-sm text-slate-900 placeholder:text-slate-400 focus:bg-white focus:ring-4 focus:outline-none"
						/>
						<kbd
							class="border-brand-100 pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 rounded border bg-white px-1.5 text-xs text-slate-400"
							aria-hidden="true">/</kbd
						>
					</label>
				</form>
			{/if}

			<div class="flex gap-1 font-medium">
				<a
					href="/search"
					class="hover:bg-brand-50 rounded-md p-2 text-slate-600 sm:hidden"
					aria-label="Search"
				>
					<Icon name="search" />
				</a>
				{#each links as link (link.href)}
					{@const active = link.match(page.url.pathname)}
					<a
						href={link.href}
						aria-current={active ? 'page' : undefined}
						class="rounded-md px-3 py-2 transition-colors {active
							? 'bg-brand-50 text-brand-700'
							: 'hover:bg-brand-50 text-slate-600 hover:text-slate-900'}"
					>
						{link.label}
					</a>
				{/each}
			</div>
		</div>
	</div>
</nav>
