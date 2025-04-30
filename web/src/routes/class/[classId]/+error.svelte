<script lang="ts">
	import Cat404 from '$lib/assets/animation2.webm';
</script>

<div class="flex h-[80vh] flex-col items-center justify-center">
	<div class="flex flex-col items-center justify-center gap-5 text-red-900">
		<div>
			<video src={Cat404} loop autoplay> <track kind="captions" /></video>
			<!-- <video src={Dog404} loop autoplay> <track kind="captions" /></video> -->
		</div>
		<div class="text-4xl font-[700]">Class not found : /</div>
		<div class="cursor-pointer text-xl">
			Go back <span class="underline"><a aria-label="error" href="/"> Home</a> </span>
		</div>
	</div>
</div>
