<script lang="ts">
	import DataTable from './messages-data-grid/data-table.svelte';
	import { columns } from './messages-data-grid/columns';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#if data.messagesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Courier messages</h2>
				<p class="text-muted-foreground">List of courier messages</p>
			</div>
		</div>
		<DataTable data={data.messagesResult?.data ?? []} {columns} />
	</div>
{/if}
