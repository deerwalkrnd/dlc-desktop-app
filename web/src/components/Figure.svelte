<script lang="ts">
	export let name: string;
	export let bgColor: string;
</script>

<figure>
	<div class={`${bgColor} rounded px-10 py-10 text-xl font-bold text-white lg:text-3xl`}>
		<p>
			{name}
		</p>
	</div>
</figure>
