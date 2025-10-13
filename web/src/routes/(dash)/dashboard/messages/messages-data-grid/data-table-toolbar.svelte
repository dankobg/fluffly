<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconX from '@lucide/svelte/icons/x';
	import type { Table } from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Input } from '$lib/components/ui/input/index';
	import DataTableFacetedFilter from '$lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '$lib/components/data-grid-shared/data-table-view-options.svelte';
	import { statuses, templateTypes, types } from './data';

	let { table }: { table: Table<TData> } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	const statusCol = $derived(table.getColumn('status'));
	const typeCol = $derived(table.getColumn('type'));
	const templateTypeCol = $derived(table.getColumn('templateType'));
</script>

<div class="flex items-center justify-between">
	<div class="flex flex-1 items-center space-x-2">
		<Input
			placeholder="Filter by recipient..."
			value={(table.getColumn('recipient')?.getFilterValue() as string) ?? ''}
			oninput={e => {
				table.getColumn('recipient')?.setFilterValue(e.currentTarget.value);
			}}
			onchange={e => {
				table.getColumn('recipient')?.setFilterValue(e.currentTarget.value);
			}}
			class="h-8 w-[150px] lg:w-[250px]"
		/>

		{#if statusCol}
			<DataTableFacetedFilter column={statusCol} title="Status" options={statuses} />
		{/if}
		{#if typeCol}
			<DataTableFacetedFilter column={typeCol} title="Type" options={types} />
		{/if}
		{#if templateTypeCol}
			<DataTableFacetedFilter column={templateTypeCol} title="Template type" options={templateTypes} />
		{/if}

		{#if isFiltered}
			<Button variant="ghost" onclick={() => table.resetColumnFilters()} class="h-8 px-2 lg:px-3">
				Reset
				<IconX />
			</Button>
		{/if}
	</div>
	<DataTableViewOptions {table} />
</div>
