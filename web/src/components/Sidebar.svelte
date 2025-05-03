<script lang="ts">
	let { subjects } = $props();
	let activeSubject: string | null = $state(null);

	function setActiveSubject(subjectName: string) {
		activeSubject = subjectName;
	}
</script>

<aside class="h-full w-64 border-r border-blue-100 bg-white shadow-md">
	<div class="border-b border-blue-100 bg-blue-50 p-4">
		<h2 class="text-lg font-semibold text-blue-800">Subjects</h2>
	</div>

	<nav class="py-2">
		<ul class="space-y-1">
			{#if subjects && subjects.subjects && subjects.subjects.length > 0}
				{#each subjects.subjects as subject}
					<li class="mx-2">
						<a
							class="flex w-full cursor-pointer items-center rounded-md px-4 py-3 text-left transition-all hover:bg-blue-50 hover:shadow-sm {activeSubject ===
							subject.Name
								? 'bg-blue-100 font-medium text-blue-700'
								: 'text-blue-800'}"
							href={`${subject.Name}-${subject.ID}`}
						>
							<div
								class="mr-3 flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-blue-600"
							>
								{subject.Name.charAt(0)}
							</div>
							<span>{subject.Name}</span>
						</a>
					</li>
				{/each}
			{:else}
				<li class="px-6 py-4 text-center">
					<div class="rounded-md bg-blue-50 p-3 text-sm text-blue-500">No subjects available</div>
				</li>
			{/if}
		</ul>
	</nav>
</aside>
