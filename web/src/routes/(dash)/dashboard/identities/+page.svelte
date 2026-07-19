<script lang="ts">
	import DataTable from './identities-data-grid/data-table.svelte';
	import { columns } from './identities-data-grid/columns';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#if data.identitiesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Identities</h2>
				<p class="text-muted-foreground">List of identities</p>
			</div>
		</div>
		<DataTable data={data.identitiesResult?.data ?? []} {columns} />
	</div>
{/if}
