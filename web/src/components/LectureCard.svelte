<script lang="ts">
	import { Controls } from '$lib/controls';
	import { onMount } from 'svelte';
	import CloseButton from './CloseButton.svelte';
	import { MEDIAURL } from '$lib/constant';

	let { lectures } = $props();
	console.log(lectures);
	let dialog: any;
	let videoElement: HTMLVideoElement | null;

	onMount(() => {
		dialog = document.getElementById('confirmation-dialog');
		videoElement = document.querySelector('#video-player');
		Controls();

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

{#snippet Lecture(lectureId: number, lectureName: string, lectureNumber: number, lessons: any)}
	<div
		class="flex flex-col rounded-lg bg-white p-6 shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<h1 class="mb-4 border-b pb-2 text-2xl font-bold text-indigo-700">
			Lecture {lectureNumber}: {lectureName}
		</h1>
		<div class="ml-4 mt-2 space-y-3">
			{#each lessons as lesson}
				{#if lesson.lectureId == lectureId}
					<div>
						{@render Lesson(lesson.id, lesson.name, lesson.number, lesson.videoUrl)}
					</div>
				{/if}
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
		class="fixed inset-0 m-auto w-full max-w-5xl rounded-lg bg-blue-300 p-6 shadow-2xl"
	>
		<div class="flex flex-col">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-xl font-bold text-white">{lessonName}</h2>
				<CloseButton {closeClick} />
			</div>

			<div class="group relative overflow-hidden rounded-lg bg-black" id="video-container">
				<figure>
					<!-- CONTROLS -->
					<div
						id="controls"
						class="absolute bottom-0 left-0 w-full p-5 opacity-0 transition-opacity duration-300 ease-linear group-hover:opacity-100"
					>
						<!-- PROGRESS BAR -->
						<div id="progress-bar" class="mb-4 h-1 w-full cursor-pointer bg-white">
							<div
								id="progress-indicator"
								class="h-full w-9 bg-indigo-800 transition-all duration-500 ease-in-out"
							></div>
						</div>
						<div class="flex items-center justify-between">
							<div class="flex items-center justify-center gap-10">
								<!-- REWIND BUTTON -->
								<button id="rewind" class="transition-all duration-100 ease-linear hover:scale-125">
									<i class="material-icons w-12 text-3xl text-white">replay_10 </i>
								</button>

								<!-- PLAY BUTTON -->
								<button
									id="play-pause"
									class=" transition-all duration-100 ease-linear hover:scale-125"
								>
									<i class="material-icons inline-block text-center text-5xl text-white"
										>play_arrow</i
									>
								</button>

								<!-- FAST FORWARD BUTTON -->
								<button
									id="fast-forward"
									class="transition-all duration-100 ease-linear hover:scale-125"
								>
									<i class="material-icons w-12 text-3xl text-white">forward_10 </i>
								</button>
								<!-- VOLUME BUTTON -->

								<button id="volume" class="transition-all duration-100 ease-linear hover:scale-125">
									<i class="material-icons text-3xl text-white">volume_up</i>
								</button>
								<input
									class="my-auto flex cursor-pointer appearance-none items-center justify-center rounded-lg rounded-xl bg-gray-200"
									type="range"
									id="volumeSlider"
									min="0"
									max="1"
									step="0.01"
									value="1"
								/>
							</div>

							<div>
								<button
									id="fullscreen"
									class="transition-all duration-100 ease-linear hover:scale-125"
								>
									<i class="material-icons text-3xl text-white">fullscreen</i>
								</button>
							</div>
						</div>
					</div>

					<video id="video-player" width="100%" class="aspect-video rounded-lg">
						<track kind="captions" />
						<source src={`${MEDIAURL}/videos/${videoUrl}`} />
						Your browser does not support the video tag.
					</video>
				</figure>
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
		{@render Lecture(lecture.id, lecture.name, lecture.number, lecture.lessons)}
	{/each}
</div>
