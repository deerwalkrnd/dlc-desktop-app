<script lang="ts">
	import ClassContent from '../../../components/ClassContent.svelte';
	import Header from '../../../components/Header.svelte';
	import Sidebar from '../../../components/Sidebar.svelte';

	export let data: {
		classData?: import('$lib/data').Class;
	};
</script>

<div class="flex min-h-screen flex-col bg-slate-50">
	<Header classNumber={data.classData?.number} />

	<div class="flex flex-1 flex-col md:flex-row">
		<aside class="w-full bg-white p-4 shadow-md md:w-64">
			<Sidebar subjects={data.classData?.subjects} />
		</aside>

		<main class="flex-1 p-6">
			<ClassContent classData={data.classData} />
		</main>
	</div>
</div>
