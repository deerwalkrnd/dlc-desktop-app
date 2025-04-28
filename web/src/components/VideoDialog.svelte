<script lang="ts">
	let videoUrl = $props();
	console.log(videoUrl);
	let dialog: any;
	const closeClick = () => {
		dialog.close();
	};
</script>

<dialog id="confirmation-dialog">
	<video width="640" height="360" controls>
		<source
			src={'/video/Aayushman%20Joshi%20-%205.1%20-%20Statistics%20-%201%20-%20Mean,%20Median,%20and%20Mode%20-%20Mathematics%20-%20OLD%20-%208%20.mp4'}
		/>
		Your browser does not support the video tag. Your browser does not support the video tag.
	</video>
	<button on:click={closeClick}> Close </button>
</dialog>
