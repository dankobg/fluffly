import DataTableCellStatus from './data-table-cell-status.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import { renderComponent } from '@tanstack/svelte-table';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';

const columnHelper = createAppColumnHelper<components['schemas']['Organization']>();

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
				href: `/dashboard/organizations/${row.original.id}`
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
	columnHelper.accessor('status', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Status'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellStatus, {
				value: row.original.status
			});
		}
	}),
	columnHelper.accessor('website', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Website'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('mission_statement', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Mission statement'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('adoption_policy', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Adoption policy'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('adoption_url', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Adoption URL'
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
