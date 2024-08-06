<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<div class="mx-auto flex w-full max-w-xl flex-col gap-2 p-4">
	{@render children()}
</div>
