<script lang="ts">
	import type { Snippet } from 'svelte';
	import Icon, { type IconName } from './Icon.svelte';

	interface Props {
		title: string;
		icon?: IconName;
		children?: Snippet;
	}
	let { title, icon = 'folder', children }: Props = $props();
</script>

<div
	class="border-brand-200 mx-auto flex max-w-xl flex-col items-center gap-3 rounded-xl border border-dashed bg-white px-6 py-12 text-center"
>
	<Icon name={icon} class="size-10 text-slate-400" />
	<h2 class="text-lg font-semibold text-slate-700">{title}</h2>
	{#if children}
		<div class="text-slate-500">{@render children()}</div>
	{/if}
</div>
