<script lang="ts">
	export let subjects: { name: string; icon?: string; count?: number }[] = [];
	export let activeSubject: string | null = null;
	export let onSelectSubject: (subject: string) => void = () => {};
</script>

<aside class="h-full w-64 border-r border-gray-200 bg-gray-100 shadow-sm">
	<div class="border-b border-gray-200 p-4">
		<h2 class="text-lg font-semibold text-gray-800">Subjects</h2>
	</div>

	<nav class="mt-2">
		<ul>
			{#each subjects as subject}
				<li>
					<button
						class="flex w-full items-center px-4 py-3 text-left transition-colors hover:bg-gray-200 {activeSubject ===
						subject.name
							? 'bg-gray-200 font-medium'
							: ''}"
						on:click={() => onSelectSubject(subject.name)}
					>
						{#if subject.icon}
							<span class="mr-3 text-gray-500">{subject.icon}</span>
						{/if}
						<span class="text-gray-700">{subject.name}</span>
						{#if subject.count !== undefined}
							<span
								class="ml-auto rounded-full bg-gray-200 px-2 py-1 text-xs font-medium text-gray-600"
							>
								{subject.count}
							</span>
						{/if}
					</button>
				</li>
			{/each}

			{#if subjects.length === 0}
				<li class="px-4 py-3 text-sm italic text-gray-500">No subjects available</li>
			{/if}
		</ul>
	</nav>

	<div class="mt-auto border-t border-gray-200 p-4"></div>
</aside>
