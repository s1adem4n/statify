<script lang="ts">
	import client from '$lib/api';
</script>

<h1 class="text-2xl font-bold">Domains</h1>

{#await client.views.domains()}
	<p>Loading...</p>
{:then domains}
	{#each Object.entries(domains) as [domain, count]}
		<div class="flex justify-between rounded-xl border border-gray-100 px-4 py-2 shadow-sm">
			<a href="/domains?domain={domain}">{domain}</a>
			<span>{count} Views (total)</span>
		</div>
	{/each}
{/await}
