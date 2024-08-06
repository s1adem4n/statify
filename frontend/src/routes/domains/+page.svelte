<script lang="ts">
	import { page } from '$app/stores';
	import client from '$lib/api';
	import Counter from './counter.svelte';
	import TopTable from './top-table.svelte';

	const todayMidnight = new Date();
	todayMidnight.setHours(0, 0, 0, 0);
	const yesterdayMidnight = new Date();
	yesterdayMidnight.setDate(yesterdayMidnight.getDate() - 1);
	yesterdayMidnight.setHours(0, 0, 0, 0);

	// gets the currents week monday
	const thisWeek = new Date();
	thisWeek.setDate(thisWeek.getDate() - thisWeek.getDay() + 1);
	thisWeek.setHours(0, 0, 0, 0);
	const lastWeek = new Date(thisWeek);
	lastWeek.setDate(lastWeek.getDate() - 7);

	const thisMonth = new Date();
	thisMonth.setDate(1);
	thisMonth.setHours(0, 0, 0, 0);
	const lastMonth = new Date(thisMonth);
	lastMonth.setMonth(lastMonth.getMonth() - 1);

	const beginningOfTime = new Date(0);

	let domain = $page.url.searchParams.get('domain');
</script>

{#if domain}
	<h1 class="text-2xl font-bold">{domain}</h1>

	<h2 class="mt-2 text-xl font-bold">Unique Views</h2>
	<div class="grid grid-cols-2 gap-2">
		<Counter {domain} start={todayMidnight} compare={yesterdayMidnight} text="Today" />
		<Counter {domain} start={thisWeek} compare={lastWeek} text="This Week" />
		<Counter {domain} start={thisMonth} compare={lastMonth} text="This Month" />
		<Counter {domain} start={beginningOfTime} text="All Time" />
	</div>

	<TopTable {domain} type="paths" />
	<TopTable {domain} type="devices" />
	<TopTable {domain} type="sessions" />
{:else}
	Please specify a domain.
{/if}
