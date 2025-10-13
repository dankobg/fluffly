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
	import { states } from './data';

	let { table }: { table: Table<TData> } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	const stateCol = $derived(table.getColumn('state'));
</script>

<div class="flex gap-4">
	<Button href="/dashboard/identities/create">Create</Button>
</div>

<div class="flex items-center justify-between">
	<div class="flex flex-1 items-center space-x-2">
		<Input
			placeholder="Filter by email..."
			value={(table.getColumn('email')?.getFilterValue() as string) ?? ''}
			oninput={e => {
				table.getColumn('email')?.setFilterValue(e.currentTarget.value);
			}}
			onchange={e => {
				table.getColumn('email')?.setFilterValue(e.currentTarget.value);
			}}
			class="h-8 w-[150px] lg:w-[250px]"
		/>
		<Input
			placeholder="Filter by username..."
			value={(table.getColumn('username')?.getFilterValue() as string) ?? ''}
			oninput={e => {
				table.getColumn('username')?.setFilterValue(e.currentTarget.value);
			}}
			onchange={e => {
				table.getColumn('username')?.setFilterValue(e.currentTarget.value);
			}}
			class="h-8 w-[150px] lg:w-[250px]"
		/>

		{#if stateCol}
			<DataTableFacetedFilter column={stateCol} title="State" options={states} />
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
