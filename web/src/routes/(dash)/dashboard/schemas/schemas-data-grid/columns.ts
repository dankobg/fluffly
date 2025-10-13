import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '$lib/gen/fluffly_openapi';

export const columns: ColumnDef<components['schemas']['IdentitySchemaContainer']>[] = [
	{
		id: 'select',
		header: ({ table }) =>
			renderComponent(DataTableCheckbox, {
				checked: table.getIsAllPageRowsSelected(),
				onCheckedChange: value => table.toggleAllPageRowsSelected(!!value),
				'aria-label': 'Select all',
				class: 'translate-y-[2px]'
			}),
		cell: ({ row }) =>
			renderComponent(DataTableCheckbox, {
				checked: row.getIsSelected(),
				onCheckedChange: value => row.toggleSelected(!!value),
				'aria-label': 'Select row',
				class: 'translate-y-[2px]'
			}),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['IdentitySchemaContainer'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/schemas/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		id: 'actions',
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['IdentitySchemaContainer']>, { row })
	}
];
