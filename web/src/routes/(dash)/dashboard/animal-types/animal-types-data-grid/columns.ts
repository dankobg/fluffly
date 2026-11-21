import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { renderComponent } from '@tanstack/svelte-table';

const columnHelper = createAppColumnHelper<components['schemas']['AnimalType']>();

export const columns = columnHelper.columns([
	columnHelper.display({
		id: 'select',
		header: ({ header }) => renderComponent(header.SelectHeader),
		cell: ({ cell }) => renderComponent(cell.SelectCell),
		enableSorting: false,
		enableHiding: false,
		enableResizing: false,
		maxSize: 40
	}),
	columnHelper.accessor('id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'ID'
			});
		},
		cell: ({ cell, row }) => {
			return renderComponent(cell.DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/animal-types/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('name', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Name'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('created_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Create time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		}
	}),
	columnHelper.accessor('updated_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Update time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		}
	}),
	columnHelper.display({
		id: 'actions',
		cell: () => renderComponent(DataTableRowActions),
		enableSorting: false,
		enableHiding: false,
		enableResizing: false,
		maxSize: 40
	})
]);
