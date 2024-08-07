<script lang="ts">
	import client from '$lib/api';

	let { domain, type }: { domain: string; type: 'paths' | 'devices' | 'sessions' } = $props();

	const todayMidnight = new Date();
	todayMidnight.setHours(0, 0, 0, 0);
	const thisWeek = new Date();
	thisWeek.setDate(thisWeek.getDate() - thisWeek.getDay() + 1);
	thisWeek.setHours(0, 0, 0, 0);
	const thisMonth = new Date();
	thisMonth.setDate(1);
	thisMonth.setHours(0, 0, 0, 0);
	const beginningOfTime = new Date(0);

	const times = {
		Today: todayMidnight,
		'This Week': thisWeek,
		'This Month': thisMonth,
		'All Time': beginningOfTime
	};

	let selectedTime: keyof typeof times = $state('Today');

	let data: Record<string, number> = $state({});

	$effect(() => {
		let promise: Promise<Record<string, number>>;

		switch (type) {
			case 'paths':
				promise = client.views.paths({ domain, start: times[selectedTime] });
				break;
			case 'devices':
				promise = client.views.devices({ domain, start: times[selectedTime] });
				break;
			case 'sessions':
				promise = client.views.sessions({ domain, start: times[selectedTime] });
				break;
		}

		promise.then((result) => {
			data = Object.fromEntries(Object.entries(result).sort((a, b) => b[1] - a[1]));
		});
	});
</script>

<h2 class="mt-4 text-xl font-bold capitalize">Top {type}</h2>
<select
	bind:value={selectedTime}
	class="rounded-xl border border-gray-100 px-4 py-2 shadow-sm focus:border-blue-500 focus:outline-none"
>
	{#each Object.keys(times) as time}
		<option value={time}>{time}</option>
	{/each}
</select>
<div class="-mt-2 grid grid-cols-[auto_min-content]">
	<span class="border-b border-gray-100 py-2 font-bold capitalize">{type}</span>
	<span class="border-b border-gray-100 py-2 font-bold">Views</span>
	{#each Object.entries(data) as [item, count]}
		{#if type === 'paths'}
			<a
				href="https://{domain}{item}"
				target="_blank"
				rel="noopener noreferrer"
				class="py-1 hover:underline">{item}</a
			>
		{:else}
			<span class="py-1">{item}</span>
		{/if}
		<span class="py-1">{count}</span>
	{/each}
	{#if Object.keys(data).length === 0}
		<span class="py-1">No data for this timespan</span>
	{/if}
</div>
