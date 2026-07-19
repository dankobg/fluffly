<script lang="ts" module>
	type TData = unknown;
	type TValue = unknown;
</script>

<script lang="ts" generics="TData, TValue">
	import {
		type ColumnDef,
		type ColumnFiltersState,
		type PaginationState,
		type RowSelectionState,
		type SortingState,
		type VisibilityState,
		getCoreRowModel,
		getFacetedRowModel,
		getFacetedUniqueValues,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel
	} from '@tanstack/table-core';
	import DataTableToolbar from './data-table-toolbar.svelte';
	import DataTablePagination from '$lib/components/data-grid-shared/data-table-pagination.svelte';
	import { createSvelteTable } from '$lib/components/ui/data-table/data-table.svelte';
	import FlexRender from '$lib/components/ui/data-table/flex-render.svelte';
	import * as Table from '$lib/components/ui/table/index';
	import type { components } from '$lib/gen/fluffly_openapi';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let {
		columns,
		data,
		meta,
		animalSpecies
	}: {
		columns: ColumnDef<TData, TValue>[];
		data: TData[];
		meta: components['schemas']['PaginationMeta'];
		animalSpecies: components['schemas']['AnimalSpecie'][];
	} = $props();

	let rowSelection = $state<RowSelectionState>({});
	let columnVisibility = $state<VisibilityState>({});
	let columnFilters = $state<ColumnFiltersState>([]);
	let sorting = $state<SortingState>([]);
	let pagination = $state<PaginationState>({ pageIndex: meta.page - 1, pageSize: meta.page_size });

	function gotoWithFilters(params: URLSearchParams) {
		goto(page.url.pathname + params.size ? `?${params}` : '', { keepFocus: true });
	}

	const table = createSvelteTable({
		get data() {
			return data;
		},
		get rowCount() {
			return meta.total;
		},
		manualPagination: true,
		manualFiltering: true,
		state: {
			get sorting() {
				return sorting;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get rowSelection() {
				return rowSelection;
			},
			get columnFilters() {
				return columnFilters;
			},
			get pagination() {
				return pagination;
			}
		},
		columns,
		enableRowSelection: true,
		onRowSelectionChange: updater => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}
		},
		onSortingChange: updater => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
		},
		onColumnFiltersChange: updater => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}
		},
		onColumnVisibilityChange: updater => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}
		},
		onPaginationChange: updater => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
			const sp = new URLSearchParams(page.url.searchParams);
			sp.set('page', `${pagination.pageIndex + 1}`);
			sp.set('page_size', `${pagination.pageSize}`);
			gotoWithFilters(sp);
		},
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues()
	});

	function getStickyColumnClasses(id: string): string {
		if (id === 'actions') return `sticky right-0`;
		if (id === 'select') return `sticky left-0`;
		return '';
	}
</script>

<div class="space-y-4">
	<DataTableToolbar {table} {animalSpecies} />
	<div class="rounded-md border">
		<Table.Root>
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row class="custom-row">
						{#each headerGroup.headers as header (header.id)}
							<Table.Head colspan={header.colSpan} class={getStickyColumnClasses(header.column.id)}>
								{#if !header.isPlaceholder}
									<FlexRender content={header.column.columnDef.header} context={header.getContext()} />
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#each table.getRowModel().rows as row (row.id)}
					<Table.Row data-state={row.getIsSelected() && 'selected'} class="custom-row">
						{#each row.getVisibleCells() as cell (cell.id)}
							<Table.Cell class={getStickyColumnClasses(cell.column.id)}>
								<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
							</Table.Cell>
						{/each}
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
	<DataTablePagination {table} />
</div>
