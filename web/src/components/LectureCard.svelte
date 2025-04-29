<script lang="ts">
	import { onMount } from 'svelte';
	import CloseButton from './CloseButton.svelte';
	import { MEDIAURL } from '$lib/constant';

	let { lectures, lessons } = $props();
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

{#snippet Lecture(lectureId: number, lectureName: string, lectureNumber: number)}
	<div
		class="flex flex-col rounded-lg bg-white p-6 shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<h1 class="mb-4 border-b pb-2 text-2xl font-bold text-indigo-700">
			Lecture {lectureNumber}: {lectureName}
		</h1>
		<div class="ml-4 mt-2 space-y-3">
			{#each lessons as lessonGroup}
				{#each lessonGroup as lesson}
					{#if lesson.ID == lectureId}
						<div>
							{@render Lesson(lesson.ID, lesson.Name, lesson.Number, lesson.VideoUrl)}
						</div>
					{/if}
				{/each}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet Lesson(lessonId: number, lessonName: string, lessonNumber: number, videoUrl: string)}
	<div class="rounded-md bg-indigo-50 p-3 transition-colors duration-200 hover:bg-indigo-100">
		<button
			type="button"
			class="flex w-full cursor-pointer items-center text-left text-lg font-semibold text-gray-800 focus:outline-none"
			onclick={() => showDialogClick(true)}
			onkeydown={(e) => e.key === 'Enter' && showDialogClick(true)}
		>
			<span class="mr-2 text-indigo-600">Lesson {lessonNumber}:</span>
			{lessonName}
		</button>
		{@render VideoDialog('Lesson ' + lessonNumber + ' : ' + lessonName, videoUrl, closeClick)}
	</div>
{/snippet}

{#snippet VideoDialog(lessonName: string, videoUrl: string, closeClick: () => void)}
	<dialog
		id="confirmation-dialog"
		class="fixed inset-0 m-auto w-full max-w-5xl rounded-lg bg-blue-500 p-6 shadow-2xl"
	>
		<div class="flex flex-col">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-xl font-bold text-white">{lessonName}</h2>
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
	{#each lectures as lecture}
		{@render Lecture(lecture.ID, lecture.Name, lecture.Number)}
	{/each}
</div>
