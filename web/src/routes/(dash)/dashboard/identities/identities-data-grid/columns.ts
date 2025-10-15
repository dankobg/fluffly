import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import DataTableCellState from './data-table-cell-state.svelte';
import DataTableCellSchemaId from './data-table-cell-schema-id.svelte';
import DataTableCellMainVerifiableAddress from './data-table-cell-main-verifiable-address.svelte';
import type { components } from '$lib/gen/fluffly_openapi';
import type { CustomTraits } from '$lib/kratos/service';

export const columns: ColumnDef<components['schemas']['Identity']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/identities/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'schema_id',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'Schema id',
				column
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellSchemaId, {
				value: row.original.schema_id
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorFn: row => (row.traits as CustomTraits)?.['email'],
		id: 'email',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'E-Mail',
				column
			}),
		cell: ({ row }) => {
			const emailSnippet = createRawSnippet<[{ email: string }]>(getEmail => {
				const { email } = getEmail();
				return {
					render: () => `<div>${email}</div>`
				};
			});
			return renderSnippet(emailSnippet, {
				email: row.getValue('email') as string
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorFn: row => (row.traits as CustomTraits)?.['username'],
		id: 'username',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'Username',
				column
			}),
		cell: ({ row }) => {
			const usernameSnippet = createRawSnippet<[{ username: string }]>(getUsername => {
				const { username } = getUsername();
				return {
					render: () => `<div>${username}</div>`
				};
			});
			return renderSnippet(usernameSnippet, {
				username: row.getValue('username') as string
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorFn: row => (row.traits as CustomTraits)?.['first_name'],
		id: 'first_name',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'First name',
				column
			}),
		cell: ({ row }) => {
			const firstNameSnippet = createRawSnippet<[{ firstName: string }]>(getFirstName => {
				const { firstName } = getFirstName();
				return {
					render: () => `<div>${firstName}</div>`
				};
			});
			return renderSnippet(firstNameSnippet, {
				firstName: row.getValue('first_name') as string
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorFn: row => (row.traits as CustomTraits)?.['last_name'],
		id: 'last_name',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'Last name',
				column
			}),
		cell: ({ row }) => {
			const lastNameSnippet = createRawSnippet<[{ lastName: string }]>(getFirstName => {
				const { lastName } = getFirstName();
				return {
					render: () => `<div>${lastName}</div>`
				};
			});
			return renderSnippet(lastNameSnippet, {
				lastName: row.getValue('last_name') as string
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'state',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'State',
				column
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellState, {
				value: row.original.state
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'state_changed_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'State changed time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const stateChangedAtSnippet = createRawSnippet<[{ stateChangedAt: string | null | undefined }]>(
				getStateChangedAt => {
					const { stateChangedAt } = getStateChangedAt();
					return {
						render: () => `<div>${stateChangedAt && fmt.format(new Date(stateChangedAt))}</div>`
					};
				}
			);
			return renderSnippet(stateChangedAtSnippet, {
				stateChangedAt: row.original.state_changed_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorFn: row => row.recovery_addresses?.[0]?.value ?? '',
		id: 'mainRecoveryAddress',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'Main recovery address',
				column
			}),
		cell: ({ row }) => {
			const mainRecoveryAddressSnippet = createRawSnippet<[{ mainRecoveryAddress: string | undefined }]>(
				getMainRecoveryAddress => {
					const { mainRecoveryAddress } = getMainRecoveryAddress();
					return {
						render: () => `<div>${mainRecoveryAddress}</div>`
					};
				}
			);
			return renderSnippet(mainRecoveryAddressSnippet, {
				mainRecoveryAddress: row.getValue('mainRecoveryAddress') as string | undefined
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorFn: row => row.verifiable_addresses?.[0],
		id: 'mainVerifiableAddress',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'Main verifiable address',
				column
			}),
		cell: ({ row }) => {
			return renderComponent(DataTableCellMainVerifiableAddress, {
				address: (row.getValue('mainVerifiableAddress') as components['schemas']['VerifiableIdentityAddress']).value,
				verified: (row.getValue('mainVerifiableAddress') as components['schemas']['VerifiableIdentityAddress']).verified
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
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
			const createdAtSnippet = createRawSnippet<[{ createdAt: string | undefined }]>(getCreatedAt => {
				const { createdAt } = getCreatedAt();
				return {
					render: () => `<div>${createdAt && fmt.format(new Date(createdAt))}</div>`
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
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
			const updatedAtSnippet = createRawSnippet<[{ updatedAt: string | undefined }]>(getUpdatedAt => {
				const { updatedAt } = getUpdatedAt();
				return {
					render: () => `<div>${updatedAt && fmt.format(new Date(updatedAt))}</div>`
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Identity']>, { row })
	}
];
