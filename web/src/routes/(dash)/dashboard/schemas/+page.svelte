<script lang="ts">
	import DataTable from './schemas-data-grid/data-table.svelte';
	import { columns } from './schemas-data-grid/columns';
	import type { PageProps } from './$types';

	let { data, params }: PageProps = $props();
</script>

{#if data.schemasResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Identity schemas</h2>
				<p class="text-muted-foreground">List of identity schemas</p>
			</div>
		</div>
		<DataTable data={data.schemasResult?.data ?? []} {columns} />
	</div>
{/if}
