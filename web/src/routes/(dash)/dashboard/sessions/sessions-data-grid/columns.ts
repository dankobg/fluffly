import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import DataTableCellActive from './data-table-cell-active.svelte';
import type { AuthenticatorAssuranceLevel, components } from '$lib/gen/fluffly_openapi';
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
		accessorKey: 'authenticated_at',
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
			const authenticatedAtSnippet = createRawSnippet<[{ authenticatedAt: string | undefined }]>(getAuthenticatedAt => {
				const { authenticatedAt } = getAuthenticatedAt();
				return {
					render: () => `<div>${authenticatedAt && fmt.format(new Date(authenticatedAt))}</div>`
				};
			});
			return renderSnippet(authenticatedAtSnippet, {
				authenticatedAt: row.original.authenticated_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'authentication_methods',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Authentication methods',
				column
			});
		},
		cell: ({ row }) => {
			const authMethodsSnippet = createRawSnippet<[{ authenticationMethods: string }]>(getAuthMethods => {
				const { authenticationMethods } = getAuthMethods();
				return {
					render: () => `<div>${authenticationMethods}</div>`
				};
			});
			const value =
				row.original.authentication_methods
					?.map(x => `${x.method} - ${x.aal}${x.provider ? ` - ${x.provider}` : ''}`)
					.join(', ') ?? '';
			return renderSnippet(authMethodsSnippet, {
				authenticationMethods: value
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'authenticator_assurance_level',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Authenticator assurance level',
				column
			});
		},
		cell: ({ row }) => {
			const aalSnippet = createRawSnippet<[{ aal: AuthenticatorAssuranceLevel | undefined }]>(getAal => {
				const { aal } = getAal();
				return {
					render: () => `<div>${aal && aal}</div>`
				};
			});
			return renderSnippet(aalSnippet, {
				aal: row.original.authenticator_assurance_level
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'expires_at',
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
			const expiresAtSnippet = createRawSnippet<[{ expiresAt: string | undefined }]>(getExpiresAt => {
				const { expiresAt } = getExpiresAt();
				return {
					render: () => `<div>${expiresAt && fmt.format(new Date(expiresAt))}</div>`
				};
			});
			return renderSnippet(expiresAtSnippet, {
				expiresAt: row.original.expires_at
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'issued_at',
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
			const issuedAtSnippet = createRawSnippet<[{ issuedAt: string | undefined }]>(getIssuedAt => {
				const { issuedAt } = getIssuedAt();
				return {
					render: () => `<div>${issuedAt && fmt.format(new Date(issuedAt))}</div>`
				};
			});
			return renderSnippet(issuedAtSnippet, {
				issuedAt: row.original.issued_at
			});
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
			const identityEmailSnippet = createRawSnippet<[{ email: string }]>(getIdentityEmail => {
				const { email } = getIdentityEmail();
				return {
					render: () => `<div>${email ?? ''}</div>`
				};
			});
			return renderSnippet(identityEmailSnippet, {
				email: row.getValue('email') as string
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorFn: row => row.identity?.traits,
		id: 'full_name',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Session'], unknown>, {
				title: 'Full name',
				column
			});
		},
		cell: ({ row }) => {
			const identityFullNameSnippet = createRawSnippet<[{ fullName: string }]>(getIdentityFullName => {
				const { fullName } = getIdentityFullName();
				return {
					render: () => `<div>${fullName}</div>`
				};
			});
			const traits = row.getValue('full_name') as CustomTraits;
			const fullName = `${traits?.first_name ?? ''}${traits?.last_name ? ' ' : ''}${traits?.last_name ?? ''}`;
			return renderSnippet(identityFullNameSnippet, {
				fullName
			});
		},
		filterFn: (row, id, value) => {
			const traits = row.getValue(id) as CustomTraits;
			return `${traits.first_name} ${traits.last_name}`.toLowerCase().includes((value as string).toLowerCase());
		}
	}
];

export const columnsWithIdentity: ColumnDef<components['schemas']['Session']>[] = columns.toSpliced(
	2,
	0,
	...identityColumns
);
