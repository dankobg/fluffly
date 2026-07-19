import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import type { components } from '$lib/gen/fluffly_openapi';

export const columns: ColumnDef<components['schemas']['Animal']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: `${row.original.id}`,
				href: `/dashboard/animals/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
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
		accessorKey: 'type',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Animal type',
				column
			});
		},
		cell: ({ row }) => {
			const typeSnippet = createRawSnippet<[{ type: string | undefined }]>(getType => {
				const { type } = getType();
				return {
					render: () => `<div>${type ?? ''}</div>`
				};
			});
			return renderSnippet(typeSnippet, {
				type: row.original.type.name
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'specie',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Specie',
				column
			});
		},
		cell: ({ row }) => {
			const speciesSnippet = createRawSnippet<[{ species: string | undefined }]>(getSpecies => {
				const { species } = getSpecies();
				return {
					render: () => `<div>${species}</div>`
				};
			});
			return renderSnippet(speciesSnippet, {
				species: row.original.specie.name
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'gender',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Gender',
				column
			});
		},
		cell: ({ row }) => {
			const genderSnippet = createRawSnippet<[{ gender: string | undefined }]>(getGender => {
				const { gender } = getGender();
				return {
					render: () => `<div>${gender ? (gender === 'm' ? 'male' : 'female') : ''}</div>`
				};
			});
			return renderSnippet(genderSnippet, {
				gender: row.original.gender
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'hermaphrodite',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Hermaphrodite',
				column
			});
		},
		cell: ({ row }) => {
			const hermaphroditeSnippet = createRawSnippet<[{ hermaphrodite: boolean }]>(getHermaphrodite => {
				const { hermaphrodite } = getHermaphrodite();
				return {
					render: () => `<div>${hermaphrodite}</div>`
				};
			});
			return renderSnippet(hermaphroditeSnippet, {
				hermaphrodite: row.original.hermaphrodite
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'age',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'age',
				column
			});
		},
		cell: ({ row }) => {
			const ageSnippet = createRawSnippet<[{ age: string | undefined }]>(getAge => {
				const { age } = getAge();
				return {
					render: () => `<div>${age ?? ''}</div>`
				};
			});
			return renderSnippet(ageSnippet, {
				age: row.original.age
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'size',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Size',
				column
			});
		},
		cell: ({ row }) => {
			const sizeSnippet = createRawSnippet<[{ size: string | undefined }]>(getSize => {
				const { size } = getSize();
				return {
					render: () => `<div>${size ?? ''}</div>`
				};
			});
			return renderSnippet(sizeSnippet, {
				size: row.original.size
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'status',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Status',
				column
			});
		},
		cell: ({ row }) => {
			const statusSnippet = createRawSnippet<[{ status: string | undefined }]>(getStatus => {
				const { status } = getStatus();
				return {
					render: () => `<div>${status ?? ''}</div>`
				};
			});
			return renderSnippet(statusSnippet, {
				status: row.original.status
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'description',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Description',
				column
			});
		},
		cell: ({ row }) => {
			const descriptionSnippet = createRawSnippet<[{ description: string | undefined }]>(getDescription => {
				const { description } = getDescription();
				return {
					render: () => `<div>${description ?? ''}</div>`
				};
			});
			return renderSnippet(descriptionSnippet, {
				description: row.original.description
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'adoption_id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
				title: 'Adoption ID',
				column
			});
		},
		cell: ({ row }) => {
			const sizeSnippet = createRawSnippet<[{ adoptionId: number | undefined }]>(getAdoptionId => {
				const { adoptionId } = getAdoptionId();
				return {
					render: () => `<div>${adoptionId ?? ''}</div>`
				};
			});
			return renderSnippet(sizeSnippet, {
				adoptionId: row.original.adoption_id
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Animal'], unknown>, {
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Animal']>, { row })
	}
];
