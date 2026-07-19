import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, renderSnippet } from '$lib/components/ui/data-table/index';
import DataTableCheckbox from '$lib/components/data-grid-shared/data-table-checkbox.svelte';
import DataTableCellId from '$lib/components/data-grid-shared/data-table-cell-id.svelte';
import DataTableColumnHeader from '$lib/components/data-grid-shared/data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';
import { createRawSnippet } from 'svelte';
import DataTableCellStatus from './data-table-cell-status.svelte';
import DataTableCellType from './data-table-cell-type.svelte';
import DataTableCellTemplateType from './data-table-cell-template-type.svelte';
import type { components } from '$lib/gen/fluffly_openapi';

export const columns: ColumnDef<components['schemas']['Message']>[] = [
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				column,
				title: 'Id'
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellId, {
				value: row.original.id,
				href: `/dashboard/messages/${row.original.id}`
			});
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'subject',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Subject',
				column
			});
		},
		cell: ({ row }) => {
			const subjectSnippet = createRawSnippet<[{ subject: string }]>(getSubject => {
				const { subject } = getSubject();
				return {
					render: () => `<div>${subject}</div>`
				};
			});
			return renderSnippet(subjectSnippet, {
				subject: row.original.subject
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'body',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Body',
				column
			});
		},
		cell: ({ row }) => {
			const bodySnippet = createRawSnippet<[{ body: string }]>(getBody => {
				const { body } = getBody();
				return {
					render: () => `<div>${body}</div>`
				};
			});
			return renderSnippet(bodySnippet, {
				body: row.original.body
			});
		},
		enableSorting: false
	},
	{
		accessorKey: 'channel',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Channel',
				column
			});
		},
		cell: ({ row }) => {
			const channelSnippet = createRawSnippet<[{ channel: string | undefined }]>(getChannel => {
				const { channel } = getChannel();
				return {
					render: () => `<div>${channel}</div>`
				};
			});
			return renderSnippet(channelSnippet, {
				channel: row.original.channel
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'recipient',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Recipient',
				column
			});
		},
		cell: ({ row }) => {
			const recipientSnippet = createRawSnippet<[{ recipient: string }]>(getRecipient => {
				const { recipient } = getRecipient();
				return {
					render: () => `<div>${recipient}</div>`
				};
			});
			return renderSnippet(recipientSnippet, {
				recipient: row.original.recipient
			});
		},
		filterFn: (row, id, value) => {
			return (row.getValue(id) as string).includes(value);
		}
	},
	{
		accessorKey: 'status',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
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
		accessorKey: 'type',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Type',
				column
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellType, {
				value: row.original.type
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'template_type',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Template type',
				column
			});
		},
		cell: ({ row }) => {
			return renderComponent(DataTableCellTemplateType, {
				value: row.original.template_type
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'send_count',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
				title: 'Send count',
				column
			});
		},
		cell: ({ row }) => {
			const sendCountSnippet = createRawSnippet<[{ sendCount: number }]>(getSendCount => {
				const { sendCount } = getSendCount();
				return {
					render: () => `<div>${sendCount}</div>`
				};
			});
			return renderSnippet(sendCountSnippet, {
				sendCount: row.original.send_count
			});
		},
		filterFn: (row, id, value) => {
			return value.includes(row.getValue(id));
		}
	},
	{
		accessorKey: 'created_at',
		header: ({ column }) => {
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
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
			return renderComponent(DataTableColumnHeader<components['schemas']['Message'], unknown>, {
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
		cell: ({ row }) => renderComponent(DataTableRowActions<components['schemas']['Message']>, { row })
	}
];
