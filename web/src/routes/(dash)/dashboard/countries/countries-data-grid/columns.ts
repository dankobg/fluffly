import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import type { components } from '$lib/gen/fluffly_openapi';

export const columns: ColumnDef<components['schemas']['Country']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: `${row.original.id}`,
				href: `/dashboard/countries/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				title: 'Name',
				column
			});
		},
		cell: ({ row }) => {
			const nameSnippet = createRawSnippet<[{ name: string }]>(getName => {
				const { name } = getName();
				return {
					render: () => `<div>${name}</div>`
				};
			});
			return renderSnippet(nameSnippet, {
				name: row.original.name
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'iso_alpha2',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				title: 'ISO-2',
				column
			});
		},
		cell: ({ row }) => {
			const isoAlpha2Snippet = createRawSnippet<[{ iso_alpha2: string }]>(getIsoAlpha2 => {
				const { iso_alpha2 } = getIsoAlpha2();
				return {
					render: () => `<div>${iso_alpha2}</div>`
				};
			});
			return renderSnippet(isoAlpha2Snippet, {
				iso_alpha2: row.original.iso_alpha2
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'iso_alpha3',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				title: 'ISO-3',
				column
			});
		},
		cell: ({ row }) => {
			const isoAlpha3Snippet = createRawSnippet<[{ iso_alpha3: string | undefined }]>(getIsoAlpha3 => {
				const { iso_alpha3 } = getIsoAlpha3();
				return {
					render: () => `<div>${iso_alpha3}</div>`
				};
			});
			return renderSnippet(isoAlpha3Snippet, {
				iso_alpha3: row.original.iso_alpha3
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'iso_numeric',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				title: 'ISO numeric',
				column
			});
		},
		cell: ({ row }) => {
			const isoNumericSnippet = createRawSnippet<[{ iso_numeric: string }]>(getIsoNumeric => {
				const { iso_numeric } = getIsoNumeric();
				return {
					render: () => `<div>${iso_numeric}</div>`
				};
			});
			return renderSnippet(isoNumericSnippet, {
				iso_numeric: row.original.iso_numeric
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				title: 'Create time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const createdAtSnippet = createRawSnippet<[{ createdAt: string }]>(getCreatedAt => {
				const { createdAt } = getCreatedAt();
				return {
					render: () => `<div>${fmt.format(new Date(createdAt))}</div>`
				};
			});
			return renderSnippet(createdAtSnippet, {
				createdAt: row.original.created_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'updated_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Country'], unknown>, {
				title: 'Update time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const updatedAtSnippet = createRawSnippet<[{ updatedAt: string }]>(getUpdatedAt => {
				const { updatedAt } = getUpdatedAt();
				return {
					render: () => `<div>${fmt.format(new Date(updatedAt))}</div>`
				};
			});
			return renderSnippet(updatedAtSnippet, {
				updatedAt: row.original.updated_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		id: 'actions',
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Country']>, { row })
	}
];
