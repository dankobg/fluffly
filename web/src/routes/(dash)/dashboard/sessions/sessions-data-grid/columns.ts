import { filterFn_inDateRange, renderComponent, renderSnippet } from '@tanstack/svelte-table';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import DataTableCellActive from './data-table-cell-active.svelte';
import type { components } from '#lib/gen/fluffly_openapi.ts';
import type { CustomTraits } from '#lib/kratos/service.ts';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';

const columnHelper = createAppColumnHelper<components['schemas']['Session']>();

export const columns = columnHelper.columns([
	columnHelper.display({
		id: 'select',
		header: ({ header }) => renderComponent(header.SelectHeader),
		cell: ({ cell }) => renderComponent(cell.SelectCell),
		enableSorting: false,
		enableHiding: false,
		enableResizing: false,
		maxSize: 40
	}),
	columnHelper.accessor('id', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'ID'
			});
		},
		cell: ({ cell, row }) => {
			return renderComponent(cell.DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/sessions/${row.original.id}`
			});
		}
	}),
	columnHelper.accessor('active', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Active'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellActive, {
				value: Boolean(row.original.active)
			});
		}
	}),
	columnHelper.accessor('authenticated_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Authenticated time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		},
		filterFn: filterFn_inDateRange
	}),
	columnHelper.accessor('authentication_methods', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Authentication methods'
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
		filterFn: 'includesString'
	}),
	columnHelper.accessor('authenticator_assurance_level', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Authenticator assurance level'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	columnHelper.accessor('expires_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Expires time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		},
		filterFn: filterFn_inDateRange
	}),
	columnHelper.accessor('issued_at', {
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'Issued time'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.DateCell);
		},
		filterFn: filterFn_inDateRange
	}),
	columnHelper.display({
		id: 'actions',
		cell: () => renderComponent(DataTableRowActions),
		enableSorting: false,
		enableHiding: false,
		enableResizing: false,
		maxSize: 40
	})
]);

const identityColumnHelper = createAppColumnHelper<components['schemas']['Session']>();

const identityColumns = identityColumnHelper.columns([
	identityColumnHelper.accessor(row => (row.identity?.traits as CustomTraits)?.['email'], {
		id: 'email',
		header: ({ header }) => {
			return renderComponent(header.DataTableColumnHeader, {
				title: 'E-Mail'
			});
		},
		cell: ({ cell }) => {
			return renderComponent(cell.TextCell);
		},
		filterFn: 'includesString'
	}),
	identityColumnHelper.accessor(
		row => {
			const traits = row.identity?.traits as CustomTraits;
			return `${traits?.first_name ?? ''}${traits?.last_name ? ' ' : ''}${traits?.last_name ?? ''}`;
		},
		{
			id: 'full_name',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Full name'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			},
			filterFn: 'includesString'
		}
	)
]);

export const columnsWithIdentity = columns.toSpliced(2, 0, ...identityColumns);
