<script lang="ts">
	import DataTable from './countries-data-grid/data-table.svelte';
	import { columns } from './countries-data-grid/columns';
	import type { PageProps } from './$types';

	let { data, params }: PageProps = $props();
</script>

{#if data.countriesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Countries</h2>
				<p class="text-muted-foreground">List of Countries</p>
			</div>
		</div>
		<DataTable data={data.countriesResult?.data.data ?? []} {columns} />
	</div>
{/if}

{data.countriesResult?.data}
