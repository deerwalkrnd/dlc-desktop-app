<script lang="ts">
	import logo from '$lib/assets/dlc.png';
</script>

<nav
	class="flex flex-row items-center justify-between border-b border-blue-200 bg-white px-6 py-3 shadow-sm"
>
	<a href="/" class="flex items-center">
		<img
			alt="Deerwalk Learning Center"
			src={logo}
			width="160"
			height="50"
			class="h-auto max-h-12 w-auto"
		/>
	</a>

	<div class="flex items-center gap-8">
		<div class="flex gap-6 font-medium text-blue-800">
			<a
				href="/"
				class="relative px-2 py-1 transition-all after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-full after:origin-bottom-right after:scale-x-0 after:bg-blue-500 after:transition-transform hover:text-blue-600 hover:after:origin-bottom-left hover:after:scale-x-100"
			>
				Classes
			</a>
			<a
				href="/teachers"
				class="relative px-2 py-1 transition-all after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-full after:origin-bottom-right after:scale-x-0 after:bg-blue-500 after:transition-transform hover:text-blue-600 hover:after:origin-bottom-left hover:after:scale-x-100"
			>
				Teachers
			</a>
		</div>
	</div>
</nav>
