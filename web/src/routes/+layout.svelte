<script lang="ts">
	import Navbar from '$components/Navbar.svelte';
	export const prerender = true;
	import '../app.css';

	let { children } = $props();
</script>

<div class="flex flex-col">
	<Navbar />
	<div class="flex-1">
		{@render children()}
	</div>
</div>
