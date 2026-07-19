import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import type { components } from '$lib/gen/fluffly_openapi';
import DataTableCellStatus from './data-table-cell-status.svelte';

export const columns: ColumnDef<components['schemas']['Adoption']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: `${row.original.id}`,
				href: `/dashboard/adoptions/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'animal_id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'Animal ID',
				column
			});
		},
		cell: ({ row }) => {
			const animalIdSnippet = createRawSnippet<[{ animalId: number }]>(getAnimalId => {
				const { animalId } = getAnimalId();
				return {
					render: () => `<div>${animalId}</div>`
				};
			});
			return renderSnippet(animalIdSnippet, {
				animalId: row.original.animal_id
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'user_id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'User ID',
				column
			});
		},
		cell: ({ row }) => {
			const userIdSnippet = createRawSnippet<[{ userId: string }]>(getUserId => {
				const { userId } = getUserId();
				return {
					render: () => `<div>${userId}</div>`
				};
			});
			return renderSnippet(userIdSnippet, {
				userId: row.original.user_id
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'organization_id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'Organization ID',
				column
			});
		},
		cell: ({ row }) => {
			const organizationIdSnippet = createRawSnippet<[{ organizationId: number | undefined }]>(getOrganizationId => {
				const { organizationId } = getOrganizationId();
				return {
					render: () => `<div>${organizationId ?? ''}</div>`
				};
			});
			return renderSnippet(organizationIdSnippet, {
				organizationId: row.original.organization_id
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'status',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
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
		accessorKey: 'note',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'Note',
				column
			});
		},
		cell: ({ row }) => {
			const noteSnippet = createRawSnippet<[{ note: string | undefined }]>(getNote => {
				const { note } = getNote();
				return {
					render: () => `<div>${note ?? ''}</div>`
				};
			});
			return renderSnippet(noteSnippet, {
				note: row.original.note
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'is_permanent',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'Permanent',
				column
			});
		},
		cell: ({ row }) => {
			const isPermanentSnippet = createRawSnippet<[{ is_permanent: boolean }]>(getIsPermanent => {
				const { is_permanent } = getIsPermanent();
				return {
					render: () => `<div>${is_permanent}</div>`
				};
			});
			return renderSnippet(isPermanentSnippet, {
				is_permanent: row.original.is_permanent
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'adopted_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'Adoption time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const adoptedAtSnippet = createRawSnippet<[{ adoptedAt: string | undefined }]>(getAdoptedAt => {
				const { adoptedAt } = getAdoptedAt();
				return {
					render: () => `<div>${adoptedAt ? fmt.format(new Date(adoptedAt)) : ''}</div>`
				};
			});
			return renderSnippet(adoptedAtSnippet, {
				adoptedAt: row.original.adopted_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'returned_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
				title: 'Returned time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const returnedAtSnippet = createRawSnippet<[{ returnedAt: string | undefined }]>(getReturnedAt => {
				const { returnedAt } = getReturnedAt();
				return {
					render: () => `<div>${returnedAt ? fmt.format(new Date(returnedAt)) : ''}</div>`
				};
			});
			return renderSnippet(returnedAtSnippet, {
				returnedAt: row.original.returned_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Adoption'], unknown>, {
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Adoption']>, { row })
	}
];
