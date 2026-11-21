<script lang="ts">
	import {
		createTableState,
		type ColumnDef,
		type ColumnFiltersState,
		type ColumnVisibilityState,
		type PaginationState,
		type RowSelectionState,
		type SortingState
	} from '@tanstack/svelte-table';
	import DataTableToolbar from './data-table-toolbar.svelte';
	import * as Table from '#lib/components/ui/table/index.ts';
	import { cn } from '#lib/utils.ts';
	import { createAppTable, type AppFeatures } from '#lib/components/data-grid-shared/app-table.ts';
	import type { components } from '#lib/gen/fluffly_openapi.ts';
	import { getCommonPinningStyles } from '#lib/components/data-grid-shared/common-pinning-styles.ts';

	let {
		columns,
		data
	}: {
		columns: ColumnDef<AppFeatures, components['schemas']['Identity']>[];
		data: components['schemas']['Identity'][];
	} = $props();

	const [rowSelection, setRowSelection] = createTableState<RowSelectionState>({});
	const [columnVisibility, setColumnVisibility] = createTableState<ColumnVisibilityState>({});
	const [columnFilters, setColumnFilters] = createTableState<ColumnFiltersState>([]);
	const [sorting, setSorting] = createTableState<SortingState>([]);
	const [pagination, setPagination] = createTableState<PaginationState>({ pageIndex: 0, pageSize: 100 });

	const table = createAppTable({
		get data() {
			return data;
		},
		get columns() {
			return columns;
		},
		initialState: {
			columnPinning: { start: ['select'], end: ['actions'] }
		},
		state: {
			get sorting() {
				return sorting();
			},
			get columnVisibility() {
				return columnVisibility();
			},
			get rowSelection() {
				return rowSelection();
			},
			get columnFilters() {
				return columnFilters();
			},
			get pagination() {
				return pagination();
			}
		},
		onSortingChange: setSorting,
		onColumnVisibilityChange: setColumnVisibility,
		onRowSelectionChange: setRowSelection,
		onColumnFiltersChange: setColumnFilters,
		onPaginationChange: setPagination
	});
</script>

<table.AppTable>
	<div class="space-y-4">
		<DataTableToolbar />
		<div class="rounded-md border">
			<Table.Root>
				<Table.Header>
					{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
						<Table.Row>
							{#each headerGroup.headers.filter(header => header.column.getIsVisible()) as header (header.id)}
								<table.AppHeader {header}>
									{#snippet children(h)}
										<Table.Head
											colspan={h.colSpan}
											class={cn('relative', {
												'border-r': h.id !== 'actions',
												'text-center [&>[role=checkbox]]:mx-auto': h.column.id === 'select'
											})}
											style={`flex-grow: ${h.getSize()}; width: ${h.getSize()}px; ${getCommonPinningStyles(h.column)}`}
										>
											{#if !h.isPlaceholder}
												<h.FlexRender header={h} />
											{/if}
											<h.ResizeHandle />
										</Table.Head>
									{/snippet}
								</table.AppHeader>
							{/each}
						</Table.Row>
					{/each}
				</Table.Header>
				<Table.Body>
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row data-state={row.getIsSelected() ? 'selected' : undefined} aria-selected={row.getIsSelected()}>
							{#each row.getVisibleCells() as cell (cell.id)}
								<table.AppCell {cell}>
									{#snippet children(c)}
										<Table.Cell
											class={cn(
												c.column.id === 'actions' ? '' : 'border-r',
												c.column.id === 'select' && 'text-center [&>[role=checkbox]]:mx-auto'
											)}
											style={`flex-grow: ${c.column.getSize()}; width: ${c.column.getSize()}px; ${getCommonPinningStyles(c.column, row.getIsSelected())}`}
										>
											<!-- {#if c.getIsGrouped()}
												<c.GroupedCell />
											{:else}
												<c.FlexRender cell={c} />
											{/if} -->
											<c.FlexRender cell={c} />
										</Table.Cell>
									{/snippet}
								</table.AppCell>
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
		<table.DataTablePagination />
	</div>
</table.AppTable>
