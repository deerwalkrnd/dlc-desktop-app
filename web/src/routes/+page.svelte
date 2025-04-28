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

<div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 px-6 py-12 lg:px-8">
	<div class="mx-auto max-w-7xl">
		<div class="mb-12 text-center">
			<h1 class="mb-4 text-4xl font-bold text-gray-900">Choose Your Class</h1>
		</div>

		<div class="grid grid-cols-2 gap-6 lg:grid-cols-3 xl:grid-cols-5">
			{#each classes as classItem}
				<a
					href={`/class/${classItem.href}`}
					class="transform transition-all duration-300 hover:-translate-y-2 hover:shadow-xl"
				>
					<div class="flex h-full flex-col overflow-hidden rounded-xl bg-white shadow-md">
						<div class={`${classItem.bgColor} flex h-24 items-center justify-center`}>
							<span class="text-4xl font-bold text-white">{classItem.name}</span>
						</div>
					</div>
				</a>
			{/each}
		</div>
	</div>
</div>
