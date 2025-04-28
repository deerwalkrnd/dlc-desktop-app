<script lang="ts">
	import { onMount } from 'svelte';
	import VideoDialog from './VideoDialog.svelte';

	let { lectures, lessons } = $props();
	console.log(lessons);
	// import logo from '../../../DLC'
	let dialog: any;
	onMount(() => {
		dialog = document.getElementById('confirmation-dialog');
	});

	const showDialogClick = (asModal = true) => {
		try {
			dialog[asModal ? 'showModal' : 'show']();
		} catch (e) {
			console.log(e);
		}
	};
	const closeClick = () => {
		dialog.close();
	};
</script>

{#snippet Lesson(lessonId: number, lessonName: string, lessonNumber: number, videoUrl: string)}
	<div
		class="flex flex-col rounded-lg bg-white p-6 shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<h1 class="mb-4 border-b pb-2 text-2xl font-bold text-indigo-700">
			Lesson {lessonNumber}: {lessonName}
		</h1>
		<div class="ml-4 mt-2 space-y-3">
			{#each lectures as lecture}
				{#if lecture.ID == lessonId}
					<div>
						{@render Lecture(lecture.ID, lecture.Name, lecture.Number, videoUrl)}
					</div>
				{/if}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet Lecture(lectureId: number, lectureName: string, lectureNumber: number, videoUrl: string)}
	<div class="rounded-md bg-indigo-50 p-3 transition-colors duration-200 hover:bg-indigo-100">
		<h3
			class="flex items-center text-lg font-semibold text-gray-800"
			on:click={() => showDialogClick(true)}
		>
			<span class="mr-2 text-indigo-600">Lecture {lectureNumber}:</span>
			{lectureName}
		</h3>
		<VideoDialog {videoUrl} />
	</div>
{/snippet}

<div class="flex min-h-screen flex-col gap-6 bg-gray-50 p-8">
	{#each lessons as lessonGroup}
		{#each lessonGroup as lesson}
			{@render Lesson(lesson.ID, lesson.Name, lesson.Number, lesson.VideoUrl)}
		{/each}
	{/each}
</div>
