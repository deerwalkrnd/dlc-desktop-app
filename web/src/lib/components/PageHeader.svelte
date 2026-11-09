<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		subtitle?: string;
		children?: Snippet;
	}
	let { title, subtitle, children }: Props = $props();
</script>

<header class="border-brand-100 border-b bg-white">
	<div class="mx-auto max-w-7xl px-6 py-8">
		<h1 class="text-brand-800 text-3xl font-bold">{title}</h1>
		{#if subtitle}
			<p class="mt-1 text-slate-500">{subtitle}</p>
		{/if}
		{#if children}
			<div class="mt-4">{@render children()}</div>
		{/if}
	</div>
</header>
