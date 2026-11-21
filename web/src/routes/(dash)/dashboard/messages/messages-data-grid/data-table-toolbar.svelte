<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableFacetedFilter from '#lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import { statuses, templateTypes, types } from './data';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import FilterInput from '#lib/components/data-grid-shared/filter-input.svelte';

	const table = useTableContext();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	const recipientCol = $derived(table.getColumn('recipient'));
	const statusCol = $derived(table.getColumn('status'));
	const typeCol = $derived(table.getColumn('type'));
	const templateTypeCol = $derived(table.getColumn('template_type'));
</script>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if recipientCol}
			<FilterInput column={recipientCol} />
		{/if}

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
	<DataTableViewOptions />
</div>
