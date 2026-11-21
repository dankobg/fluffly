<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import FilterInput from '#lib/components/data-grid-shared/filter-input.svelte';

	const table = useTableContext();

	const idCol = $derived(table.getColumn('id'));
	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
</script>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if idCol}
			<FilterInput column={idCol} />
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
