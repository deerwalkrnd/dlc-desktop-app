<script lang="ts">
	import { page } from '$app/state';

	let { subject } = $props();
	let currentPath = page.url.pathname;
</script>

<section>
	<a href={`${currentPath}/subject/${subject.Name}-${subject.ID}`} class="">
		<h2
			class="mb-4 rounded-lg border-slate-200 bg-white px-4 py-6 text-center text-2xl font-semibold text-indigo-800 shadow-lg ring ring-blue-500 hover:text-indigo-600 hover:shadow-xl"
		>
			{subject.Name}
		</h2>
	</a>
</section>
