import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import { renderComponent } from '@tanstack/svelte-table';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';

const columnHelper = createAppColumnHelper<components['schemas']['AnimalSpecie']>();

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
				href: `/dashboard/animal-species/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('animal_type_id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Animal type ID'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
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
