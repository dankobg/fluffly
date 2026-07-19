<script lang="ts">
	import DataTable from './adoptions-data-grid/data-table.svelte';
	import { columns } from './adoptions-data-grid/columns';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#if data.adoptionsResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Adoptions</h2>
				<p class="text-muted-foreground">List of Adoptions</p>
			</div>
		</div>
		<DataTable data={data.adoptionsResult?.data.data ?? []} meta={data.adoptionsResult?.data.meta} {columns} />
	</div>
{/if}
