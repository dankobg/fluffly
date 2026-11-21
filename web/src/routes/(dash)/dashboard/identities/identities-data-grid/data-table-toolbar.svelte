<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableFacetedFilter from '#lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import { states } from './data';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import FilterInput from '#lib/components/data-grid-shared/filter-input.svelte';

	const table = useTableContext();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	const emailCol = $derived(table.getColumn('email'));
	const stateCol = $derived(table.getColumn('state'));
</script>

<div class="flex gap-4">
	<Button href="/dashboard/identities/create">Create</Button>
</div>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if emailCol}
			<FilterInput column={emailCol} />
		{/if}

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
	<DataTableViewOptions />
</div>
