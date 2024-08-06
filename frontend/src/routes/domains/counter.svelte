<script lang="ts">
	import client from '$lib/api';
	import { onMount } from 'svelte';

	let {
		domain,
		start,
		compare,
		text
	}: { domain: string; start: Date; compare?: Date; text: string } = $props();

	let currentCount: number = $state(0);
	let previousCount: number = $state(0);
	const diff = $derived(currentCount - previousCount);

	onMount(async () => {
		currentCount = await client.views.count({ domain, start });
		if (compare) {
			previousCount = await client.views.count({ domain, start: compare, end: start });
		}
	});
</script>

<div
	class="flex w-full flex-col gap-1 rounded-xl border border-gray-100 px-4 py-2 shadow-sm sm:px-8 sm:py-4"
>
	<span class="text-xl font-bold">{text}</span>
	<div class="flex gap-2">
		<span class="text-2xl sm:text-3xl">{currentCount} Views</span>
		{#if compare}
			<span class="mb-auto" class:text-red-500={diff < 0} class:text-green-500={diff > 0}>
				{diff > 0 ? '↑' : '↓'}
				{Math.abs(diff)}
			</span>
		{/if}
	</div>
</div>
