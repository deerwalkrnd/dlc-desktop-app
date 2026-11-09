<script lang="ts">
	import { page } from '$app/state';
	import logo from '$lib/assets/dlc.png';

	const links = [
		{ href: '/', label: 'Classes', match: (p: string) => p === '/' || p.startsWith('/class') },
		{ href: '/teachers', label: 'Teachers', match: (p: string) => p.startsWith('/teachers') },
		{ href: '/library', label: 'Library', match: (p: string) => p.startsWith('/library') }
	];
</script>

<nav class="sticky top-0 z-10 border-b border-slate-200 bg-white/95 backdrop-blur">
	<div class="mx-auto flex max-w-7xl items-center justify-between gap-6 px-6 py-3">
		<a href="/" class="flex shrink-0 items-center">
			<img alt="Deerwalk Learning Center" src={logo} width="160" height="50" class="h-10 w-auto" />
		</a>
		<div class="flex gap-1 font-medium">
			{#each links as link (link.href)}
				{@const active = link.match(page.url.pathname)}
				<a
					href={link.href}
					aria-current={active ? 'page' : undefined}
					class="rounded-md px-3 py-2 transition-colors {active
						? 'bg-blue-50 text-blue-700'
						: 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}"
				>
					{link.label}
				</a>
			{/each}
		</div>
	</div>
</nav>
