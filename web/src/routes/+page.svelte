<script lang="ts">
	import Figure from '../components/Figure.svelte';

	let classes = [
		{ name: 'Class 1', bgColor: 'bg-red-500', href: '1' },
		{ name: 'Class 2', bgColor: 'bg-blue-500', href: '2' },
		{ name: 'Class 3', bgColor: 'bg-green-500', href: '3' },
		{ name: 'Class 4', bgColor: 'bg-yellow-500', href: '4' },
		{ name: 'Class 5', bgColor: 'bg-purple-500', href: '5' },
		{ name: 'Class 6', bgColor: 'bg-pink-500', href: '6' },
		{ name: 'Class 7', bgColor: 'bg-indigo-500', href: '7' },
		{ name: 'Class 8', bgColor: 'bg-teal-500', href: '8' },
		{ name: 'Class 9', bgColor: 'bg-orange-500', href: '9' },
		{ name: 'Class 10', bgColor: 'bg-cyan-500', href: '10' }
	];
</script>

<div class="flex flex-col p-5">
	<h1>Choose Your Class</h1>
	<div class="flex flex-row flex-wrap items-center justify-center gap-5">
		{#each classes as classItem}
			<a href={`/class/${classItem.href}`}>
				<Figure name={classItem.name} bgColor={classItem.bgColor} />
			</a>
		{/each}
	</div>
</div>
