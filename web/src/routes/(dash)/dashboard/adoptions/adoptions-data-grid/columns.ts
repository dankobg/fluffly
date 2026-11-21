import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import DataTableCellStatus from './data-table-cell-status.svelte';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { renderComponent } from '@tanstack/svelte-table';

const columnHelper = createAppColumnHelper<components['schemas']['Adoption']>();

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
				href: `/dashboard/adoptions/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('animal_id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Animal ID'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('user_id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'User ID'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('organization_id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Organization ID'
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
	columnHelper.accessor('note', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Note'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('is_permanent', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Permanent'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		}
	}),
	columnHelper.accessor('adopted_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Adopted time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		}
	}),
	columnHelper.accessor('returned_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Returned time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
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
