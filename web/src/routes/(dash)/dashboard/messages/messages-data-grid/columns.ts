import DataTableCellId from '#lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import DataTableCellStatus from './data-table-cell-status.svelte';
import DataTableCellType from './data-table-cell-type.svelte';
import DataTableCellTemplateType from './data-table-cell-template-type.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { filterFn_inDateRange, renderComponent } from '@tanstack/svelte-table';

const columnHelper = createAppColumnHelper<components['schemas']['Message']>();

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
				href: `/dashboard/sessions/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('subject', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Subject'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('body', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Body'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('channel', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Channel'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('recipient', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Recipient'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
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
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('type', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Type'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellType, {
				value: row.original.type
			});
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('template_type', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Template type'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellTemplateType, {
				value: row.original.template_type
			});
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('send_count', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Send count'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('created_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Created time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		},
		filterFn: filterFn_inDateRange
	}),
	columnHelper.accessor('updated_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Updated time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		},
		filterFn: filterFn_inDateRange
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
