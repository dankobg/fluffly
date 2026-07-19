import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableCellStatus from './data-table-cell-status.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import type { components } from '$lib/gen/fluffly_openapi';

export const columns: ColumnDef<components['schemas']['Organization']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: `${row.original.id}`,
				href: `/dashboard/organizations/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
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
		accessorKey: 'status',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
				title: 'Status',
				column
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellStatus, {
				value: row.original.status
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'website',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
				title: 'Website',
				column
			});
		},
		cell: ({ row }) => {
			const websiteSnippet = createRawSnippet<[{ website: string | undefined }]>(getWebsite => {
				const { website } = getWebsite();
				return {
					render: () => `<div>${website ?? ''}</div>`
				};
			});
			return renderSnippet(websiteSnippet, {
				website: row.original.website
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'mission_statement',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
				title: 'Mission statement',
				column
			});
		},
		cell: ({ row }) => {
			const missionStatementSnippet = createRawSnippet<[{ mission_statement: string | undefined }]>(
				getMissionStatement => {
					const { mission_statement } = getMissionStatement();
					return {
						render: () => `<div>${mission_statement ?? ''}</div>`
					};
				}
			);
			return renderSnippet(missionStatementSnippet, {
				mission_statement: row.original.mission_statement
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'adoption_policy',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
				title: 'Adoption policy',
				column
			});
		},
		cell: ({ row }) => {
			const adoptionPolicySnippet = createRawSnippet<[{ adoption_policy: string | undefined }]>(getAdoptionPolicy => {
				const { adoption_policy } = getAdoptionPolicy();
				return {
					render: () => `<div>${adoption_policy ?? ''}</div>`
				};
			});
			return renderSnippet(adoptionPolicySnippet, {
				adoption_policy: row.original.adoption_policy
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'adoption_url',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
				title: 'Adoption URL',
				column
			});
		},
		cell: ({ row }) => {
			const adoptionUrlSnippet = createRawSnippet<[{ adoption_url: string | undefined }]>(getAdoptionUrl => {
				const { adoption_url } = getAdoptionUrl();
				return {
					render: () => `<div>${adoption_url ?? ''}</div>`
				};
			});
			return renderSnippet(adoptionUrlSnippet, {
				adoption_url: row.original.adoption_url
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Organization'], unknown>, {
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Organization']>, { row })
	}
];
