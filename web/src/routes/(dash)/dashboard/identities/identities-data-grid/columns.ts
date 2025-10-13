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
		accessorKey: 'schemaId',
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
			const emailSnippet = createRawSnippet<[string]>(getEmail => {
				const email = getEmail();
				return {
					render: () => `<div>${email}</div>`
				};
			});
			return renderSnippet(emailSnippet, row.getValue('email'));
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
			const usernameSnippet = createRawSnippet<[string]>(getUsername => {
				const username = getUsername();
				return {
					render: () => `<div>${username}</div>`
				};
			});
			return renderSnippet(usernameSnippet, row.getValue('username'));
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorFn: row => (row.traits as CustomTraits)?.['first_name'],
		id: 'firstName',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'First name',
				column
			}),
		cell: ({ row }) => {
			const firstNameSnippet = createRawSnippet<[string]>(getFirstName => {
				const firstName = getFirstName();
				return {
					render: () => `<div>${firstName}</div>`
				};
			});
			return renderSnippet(firstNameSnippet, row.getValue('firstName'));
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorFn: row => (row.traits as CustomTraits)?.['last_name'],
		id: 'lastName',
		header: ({ column }) =>
			renderComponent(DataTableColumnHeader<components['schemas']['Identity'], unknown>, {
				title: 'Last name',
				column
			}),
		cell: ({ row }) => {
			const lastNameSnippet = createRawSnippet<[string]>(getFirstName => {
				const lastName = getFirstName();
				return {
					render: () => `<div>${lastName}</div>`
				};
			});
			return renderSnippet(lastNameSnippet, row.getValue('lastName'));
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
		accessorKey: 'stateChangedAt',
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
			const stateChangedAtSnippet = createRawSnippet<[string]>(getStateChangedAt => {
				const stateChangedAt = getStateChangedAt();
				return {
					render: () => `<div>${stateChangedAt}</div>`
				};
			});
			return renderSnippet(stateChangedAtSnippet, fmt.format(row.getValue('stateChangedAt')));
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
			const mainRecoveryAddressSnippet = createRawSnippet<[string]>(getMainRecoveryAddress => {
				const mainRecoveryAddress = getMainRecoveryAddress();
				return {
					render: () => `<div>${mainRecoveryAddress}</div>`
				};
			});
			return renderSnippet(mainRecoveryAddressSnippet, row.getValue('mainRecoveryAddress'));
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
		accessorKey: 'createdAt',
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
			const createdAtSnippet = createRawSnippet<[string]>(getCreatedAt => {
				const createdAt = getCreatedAt();
				return {
					render: () => `<div>${createdAt}</div>`
				};
			});
			return renderSnippet(createdAtSnippet, fmt.format(row.getValue('createdAt')));
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'updatedAt',
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
			const updatedAtSnippet = createRawSnippet<[string]>(getUpdatedAt => {
				const updatedAt = getUpdatedAt();
				return {
					render: () => `<div>${updatedAt}</div>`
				};
			});
			return renderSnippet(updatedAtSnippet, fmt.format(row.getValue('updatedAt')));
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
