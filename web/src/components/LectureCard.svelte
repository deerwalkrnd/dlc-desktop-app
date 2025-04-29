<script lang="ts">
	import { onMount } from 'svelte';
	import CloseButton from './CloseButton.svelte';
	import { MEDIAURL } from '$lib/constant';

	let { lectures, lessons } = $props();
	console.log(lessons);
	let dialog: any;
	let videoElement: HTMLVideoElement | null;

	onMount(() => {
		dialog = document.getElementById('confirmation-dialog');
		videoElement = document.querySelector('#video-player');

		if (dialog) {
			dialog.addEventListener('close', () => {
				if (videoElement) {
					videoElement.pause();
					videoElement.currentTime = 0;
				}
			});
		}
	});

	const showDialogClick = (asModal = true) => {
		try {
			dialog[asModal ? 'showModal' : 'show']();
		} catch (e) {
			console.log(e);
		}
	};

	const closeClick = () => {
		if (videoElement) {
			videoElement.pause();
			videoElement.currentTime = 0;
		}
		dialog.close();
	};
</script>

{#snippet Lesson(lessonId: number, lessonName: string, lessonNumber: number, videoUrl: string)}
	<div
		class="flex flex-col rounded-lg bg-white p-6 shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<h1 class="mb-4 border-b pb-2 text-2xl font-bold text-indigo-700">
			Lesson {lessonNumber}: {lessonName}
			{lessonId}
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
			class="flex cursor-pointer items-center text-lg font-semibold text-gray-800"
			onclick={() => showDialogClick(true)}
		>
			<span class="mr-2 text-indigo-600">Lecture {lectureNumber}:</span>
			{lectureName}
		</h3>
		{videoUrl}
		{@render VideoDialog('Lecture ' + lectureNumber + ' : ' + lectureName, videoUrl, closeClick)}
	</div>
{/snippet}

{#snippet VideoDialog(lectureName: string, videoUrl: string, closeClick: () => void)}
	<dialog
		id="confirmation-dialog"
		class="fixed inset-0 m-auto w-full max-w-5xl rounded-lg bg-blue-500 p-6 shadow-2xl"
	>
		<div class="flex flex-col">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-xl font-bold text-white">{lectureName}</h2>
				<CloseButton {closeClick} />
			</div>

			<div class="overflow-hidden rounded-lg bg-black">
				<video id="video-player" width="100%" controls class="aspect-video">
					<track kind="captions" />
					<source src={`${MEDIAURL}/videos/${videoUrl}`} />
					Your browser does not support the video tag.
				</video>
			</div>

			<div class="mt-4 flex justify-end">
				<button
					onclick={closeClick}
					class="rounded-md bg-red-500 px-6 py-2 font-medium text-white transition-colors hover:bg-red-700"
				>
					Close
				</button>
			</div>
		</div>
	</dialog>
{/snippet}

<div class="flex min-h-screen flex-col gap-6 bg-gray-50 p-8">
	{#each lessons as lessonGroup}
		{#each lessonGroup as lesson}
			{@render Lesson(lesson.ID, lesson.Name, lesson.Number, lesson.VideoUrl)}
		{/each}
	{/each}
</div>
