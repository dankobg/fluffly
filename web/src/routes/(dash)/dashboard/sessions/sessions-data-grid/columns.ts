import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import DataTableCellActive from './data-table-cell-active.svelte';
import type { components } from '$lib/gen/fluffly_openapi';
import type { CustomTraits } from '$lib/kratos/service';

export const columns: ColumnDef<components['schemas']['Session']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/sessions/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'active',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Active',
				column
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellActive, {
				value: Boolean(row.original.active)
			});
		}
	},
	{
		accessorKey: 'authenticatedAt',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Authenticated time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const authenticatedAtSnippet = createRawSnippet<[string]>(getAuthenticatedAt => {
				const authenticatedAt = getAuthenticatedAt();
				return {
					render: () => `<div>${authenticatedAt}</div>`
				};
			});
			return renderSnippet(authenticatedAtSnippet, fmt.format(row.getValue('authenticatedAt')));
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'authenticationMethods',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Authentication methods',
				column
			});
		},
		cell: ({ row }) => {
			const authMethodsSnippet = createRawSnippet<[string]>(getAuthMethods => {
				const authMethods = getAuthMethods();
				return {
					render: () => `<div>${authMethods}</div>`
				};
			});
			const value =
				row.original.authentication_methods
					?.map(x => `${x.method} - ${x.aal}${x.provider ? ` - ${x.provider}` : ''}`)
					.join(', ') ?? '';
			return renderSnippet(authMethodsSnippet, value);
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'authenticatorAssuranceLevel',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Authenticator assurance level',
				column
			});
		},
		cell: ({ row }) => {
			const aalSnippet = createRawSnippet<[string]>(getAal => {
				const aal = getAal();
				return {
					render: () => `<div>${aal}</div>`
				};
			});
			return renderSnippet(aalSnippet, row.getValue('authenticatorAssuranceLevel'));
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'expiresAt',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Expires time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const expiresAtSnippet = createRawSnippet<[string]>(getExpiresAt => {
				const expiresAt = getExpiresAt();
				return {
					render: () => `<div>${expiresAt}</div>`
				};
			});
			return renderSnippet(expiresAtSnippet, fmt.format(row.getValue('expiresAt')));
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'issuedAt',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Issued time',
				column
			});
		},
		cell: ({ row }) => {
			const fmt = new Intl.DateTimeFormat(undefined, {
				dateStyle: 'short',
				timeStyle: 'short',
				hour12: false
			});
			const issuedAtSnippet = createRawSnippet<[string]>(getIssuedAt => {
				const issuedAt = getIssuedAt();
				return {
					render: () => `<div>${issuedAt}</div>`
				};
			});
			return renderSnippet(issuedAtSnippet, fmt.format(row.getValue('issuedAt')));
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'createdAt',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Session']>, { row })
	}
];

const identityColumns: ColumnDef<components['schemas']['Session']>[] = [
	{
		accessorFn: row => (row.identity?.traits as CustomTraits)?.['email'],
		id: 'email',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'E-Mail',
				column
			});
		},
		cell: ({ row }) => {
			const identityEmailSnippet = createRawSnippet<[string]>(getIdentityEmail => {
				const email = getIdentityEmail();
				return {
					render: () => `<div>${email ?? ''}</div>`
				};
			});
			return renderSnippet(identityEmailSnippet, row.getValue('email'));
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorFn: row => (row.identity?.traits as CustomTraits)?.['username'],
		id: 'username',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Username',
				column
			});
		},
		cell: ({ row }) => {
			const identityUsernameSnippet = createRawSnippet<[string]>(getUsername => {
				const username = getUsername();
				return {
					render: () => `<div>${username ?? ''}</div>`
				};
			});
			return renderSnippet(identityUsernameSnippet, row.getValue('username'));
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorFn: row => row.identity?.traits,
		id: 'fullName',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Full name',
				column
			});
		},
		cell: ({ row }) => {
			const identityFullNameSnippet = createRawSnippet<[string]>(getIdentityFullName => {
				const fullName = getIdentityFullName();
				return {
					render: () => `<div>${fullName}</div>`
				};
			});
			const traits = row.getValue('fullName') as { first_name?: string; last_name?: string };
			const fullName = `${traits?.first_name ?? ''}${traits?.last_name ? ' ' : ''}${traits?.last_name ?? ''}`;
			return renderSnippet(identityFullNameSnippet, fullName);
		},
		filterFn: (row, id, value) => {
			const traits = row.getValue(id) as { first_name?: string; last_name?: string };
			return `${traits.first_name} ${traits.last_name}`.toLowerCase().includes((value as string).toLowerCase());
		}
	}
];

export const columnsWithIdentity: ColumnDef<components['schemas']['Session']>[] = columns.toSpliced(
	2,
	0,
	...identityColumns
);
