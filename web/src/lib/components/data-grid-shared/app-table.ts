import {
	createTableHook,
	columnFacetingFeature,
	columnFilteringFeature,
	columnVisibilityFeature,
	columnOrderingFeature,
	columnPinningFeature,
	columnSizingFeature,
	columnResizingFeature,
	createFacetedRowModel,
	createFacetedUniqueValues,
	createFilteredRowModel,
	createPaginatedRowModel,
	createSortedRowModel,
	filterFn_includesString,
	rowPaginationFeature,
	rowSelectionFeature,
	rowSortingFeature,
	sortFn_alphanumeric,
	sortFn_text,
	metaHelper
} from '@tanstack/svelte-table';
import DataTableColumnHeader from './data-table-column-header.svelte';
import DataTableCellId from './data-table-cell-id.svelte';
import DataTableCell from './data-table-cell.svelte';
import DataTablePagination from './data-table-pagination.svelte';
import DataTableViewOptions from './data-table-view-options.svelte';
import SelectCell from './select-cell.svelte';
import SelectHeader from './select-header.svelte';
import TextCell from './text-cell.svelte';
import DateCell from './date-cell.svelte';
import ResizeHandle from './resize-handle.svelte';

export interface MyColumnMeta {
	label?: string;
	variant?: 'text' | 'number' | 'date' | 'boolean' | 'select' | 'multi-select';
	options?: Array<{ label: string; value: string; count?: number }>;
}

export const { createAppTable, createAppColumnHelper, appFeatures, useCellContext, useHeaderContext, useTableContext } =
	createTableHook({
		features: {
			columnFacetingFeature,
			columnFilteringFeature,
			columnVisibilityFeature,
			columnOrderingFeature,
			columnPinningFeature,
			columnSizingFeature,
			columnResizingFeature,
			rowPaginationFeature,
			rowSelectionFeature,
			rowSortingFeature,
			facetedRowModel: createFacetedRowModel(),
			facetedUniqueValues: createFacetedUniqueValues(),
			filteredRowModel: createFilteredRowModel(),
			paginatedRowModel: createPaginatedRowModel(),
			sortedRowModel: createSortedRowModel(),
			filterFns: { includesString: filterFn_includesString },
			sortFns: { alphanumeric: sortFn_alphanumeric, text: sortFn_text },
			columnMeta: metaHelper<MyColumnMeta>()
		},
		headerComponents: {
			DataTableColumnHeader,
			SelectHeader,
			ResizeHandle
		},
		cellComponents: {
			DataTableCell,
			DataTableCellId,
			SelectCell,
			TextCell,
			DateCell
		},
		tableComponents: {
			DataTablePagination
			// DataTableViewOptions
		},
		enableRowSelection: true
	});

export type AppFeatures = typeof appFeatures;
