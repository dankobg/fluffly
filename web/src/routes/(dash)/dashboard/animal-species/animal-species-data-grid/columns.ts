import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import type { components } from '$lib/gen/fluffly_openapi';

export const columns: ColumnDef<components['schemas']['AnimalSpecie']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['AnimalSpecie'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: `${row.original.id}`,
				href: `/dashboard/animal-species/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'animal_type_id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['AnimalSpecie'], unknown>, {
				title: 'Animal type id',
				column
			});
		},
		cell: ({ row }) => {
			const animalTypeIdSnippet = createRawSnippet<[{ animal_type_id: number }]>(getAnimalTypeId => {
				const { animal_type_id } = getAnimalTypeId();
				return {
					render: () => `<div>${animal_type_id}</div>`
				};
			});
			return renderSnippet(animalTypeIdSnippet, {
				animal_type_id: row.original.animal_type_id
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['AnimalSpecie'], unknown>, {
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
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['AnimalSpecie'], unknown>, {
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
			return renderComponent(DataTableColumnHeader<components['schemas']['AnimalSpecie'], unknown>, {
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['AnimalSpecie']>, { row })
	}
];
