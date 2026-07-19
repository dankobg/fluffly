<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconTrash from '@lucide/svelte/icons/trash';
	import IconPen from '@lucide/svelte/icons/pen';
	import IconCopy from '@lucide/svelte/icons/copy';
	import IconCircleCheck from '@lucide/svelte/icons/circle-check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import type { Row } from '@tanstack/table-core';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { invalidate } from '$app/navigation';
	import { fluffly } from '$lib/fluffly/client';
	import { OrganizationStatus } from '$lib/gen/fluffly_openapi';

	let { row }: { row: Row<TData> } = $props();

	const hasId = $derived(
		typeof row.original === 'object' && !!row.original && 'id' in row.original && typeof row.original.id === 'string'
	);

	function copyIdToClipboard() {
		try {
			navigator.clipboard.writeText((row.original as TData & { id: string }).id).then(() => {
				toast.success('id coppied');
			});
		} catch (error) {
			if (error instanceof Error) toast.error('failed to copy id: ' + error.message);
		}
	}

	async function approveOrganization(id: number) {
		try {
			const approveOrganizationResult = await fluffly.POST('/organizations/{id}/application/approve', {
				params: {
					path: { id }
				}
			});
			if (approveOrganizationResult.error) {
				toast.error(
					[approveOrganizationResult.error.message, approveOrganizationResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organization approved');
			invalidate('data:dashboard-organizations');
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function rejectOrganization(id: number) {
		try {
			const rejectOrganizationResult = await fluffly.POST('/organizations/{id}/application/reject', {
				params: {
					path: { id }
				}
			});
			if (rejectOrganizationResult.error) {
				toast.error(
					[rejectOrganizationResult.error.message, rejectOrganizationResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organization rejected');
			invalidate('data:dashboard-organizations');
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function deleteOrganization(id: number) {
		try {
			const deleteOrganizationResult = await fluffly.DELETE('/organizations/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteOrganizationResult.error) {
				toast.error(
					[deleteOrganizationResult.error.message, deleteOrganizationResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organization deleted');
			invalidate('data:dashboard-organizations');
		} catch (error: unknown) {
			if (error instanceof TypeError) {
				toast.error(`${error}`);
				return;
			}
			console.error('unexpected error');
		} finally {
			confirmation.closeDialog();
		}
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="ghost" class="flex h-8 w-8 p-0 data-[state=open]:bg-muted">
				<IconEllipsis />
				<span class="sr-only">Open Menu</span>
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		{#if hasId}
			<DropdownMenu.Item onclick={copyIdToClipboard}>
				<IconCopy />
				Copy ID to clipboard
			</DropdownMenu.Item>
		{/if}
		<a href="/dashboard/organizations/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
		<a href="/dashboard/organizations/{row.getValue('id')}/update">
			<DropdownMenu.Item class="cursor-pointer">
				<IconPen />
				Edit
			</DropdownMenu.Item>
		</a>
		<DropdownMenu.Item
			class="cursor-pointer"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete organization: ${row.getValue('name')}?`,
					destructive: true,
					async onConfirm() {
						await deleteOrganization(row.getValue('id'));
					}
				});
			}}
		>
			<IconTrash />
			Delete
		</DropdownMenu.Item>
		{#if row.getValue('status') === OrganizationStatus.pending}
			<DropdownMenu.Item
				class="cursor-pointer"
				onclick={() => {
					confirmation.openDialog({
						title: `Approve organization: ${row.getValue('name')}?`,
						destructive: true,
						async onConfirm() {
							await approveOrganization(row.getValue('id'));
						}
					});
				}}
			>
				<IconCircleCheck />
				Approve
			</DropdownMenu.Item>
			<DropdownMenu.Item
				class="cursor-pointer"
				onclick={() => {
					confirmation.openDialog({
						title: `Reject organization: ${row.getValue('name')}?`,
						destructive: true,
						async onConfirm() {
							await rejectOrganization(row.getValue('id'));
						}
					});
				}}
			>
				<IconCircleX />
				Reject
			</DropdownMenu.Item>
		{/if}
	</DropdownMenu.Content>
</DropdownMenu.Root>
