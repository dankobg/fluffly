import DataTableCellId from '#lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import DataTableCellState from './data-table-cell-state.svelte';
import DataTableCellSchemaId from './data-table-cell-schema-id.svelte';
import DataTableCellMainVerifiableAddress from './data-table-cell-main-verifiable-address.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import type { CustomTraits } from '#lib/kratos/service.ts';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { filterFn_inDateRange, renderComponent } from '@tanstack/svelte-table';

const columnHelper = createAppColumnHelper<components['schemas']['Identity']>();

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
				href: `/dashboard/identities/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('schema_id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Schema ID'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellSchemaId, {
				value: row.original.schema_id
			});
		}
	}),
	columnHelper.accessor(row => (row?.traits as CustomTraits)?.['email'], {
		id: 'email',
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'E-Mail'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor(row => (row?.traits as CustomTraits)?.['first_name'], {
		id: 'first_name',
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'First name'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor(row => (row?.traits as CustomTraits)?.['last_name'], {
		id: 'last_name',
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Last name'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('state', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'State'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellState, {
				value: row.original.state
			});
		}
	}),
	columnHelper.accessor('state_changed_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'State changed time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		},
		filterFn: filterFn_inDateRange
	}),
	columnHelper.accessor(row => row.recovery_addresses?.[0]?.value ?? '', {
		id: 'main_recovery_address',
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Main recovery address'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: filterFn_inDateRange
	}),
	columnHelper.accessor(row => row.verifiable_addresses?.[0]?.value, {
		id: 'main_verifiable_address',
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Main verifiable address'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellMainVerifiableAddress, {
				address: row.original.verifiable_addresses?.[0]?.value,
				verified: Boolean(row.original.verifiable_addresses?.[0]?.verified)
			});
		},
		filterFn: filterFn_inDateRange
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
