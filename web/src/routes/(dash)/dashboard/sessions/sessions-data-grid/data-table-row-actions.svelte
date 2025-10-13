<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconCopy from '@lucide/svelte/icons/copy';
	import IconTrash2 from '@lucide/svelte/icons/trash-2';
	import IconCirclePlus from '@lucide/svelte/icons/circle-plus';
	import type { Row } from '@tanstack/table-core';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { fluffly } from '$lib/fluffly/client';
	import { invalidate } from '$app/navigation';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import type { components } from '$lib/gen/fluffly_openapi';
	import type { CustomTraits } from '$lib/kratos/service';

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

	async function onConfirmDeactivateSession() {
		try {
			await fluffly.DELETE('/sessions/{id}', {
				params: {
					path: { id: row.getValue('id') }
				}
			});
			toast.success('session deactivated');
			invalidate('data:sessions');
		} catch (error) {
			console.log('err', error);
			toast.error('session deactivation failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function onConfirmExtendSession() {
		try {
			await fluffly.PATCH('/sessions/{id}/extend', {
				params: {
					path: { id: row.getValue('id') }
				}
			});
			toast.success('session extended');
			invalidate('data:sessions');
		} catch (error) {
			console.log('err', error);
			toast.error('session extend failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeactivateSessionClick(row: Row<TData>) {
		confirmation.openDialog({
			description: deactivateSessionDescriptionSnippet,
			onConfirm: onConfirmDeactivateSession
		});
	}

	function onExtendSessionClick(row: Row<TData>) {
		confirmation.openDialog({
			description: extendSessionDescriptionSnippet,
			onConfirm: onConfirmExtendSession
		});
	}
</script>

{#snippet deactivateSessionDescriptionSnippet()}
	{@const sess = row.original as components['schemas']['Session']}
	This action cannot be undone. This will deactive (invalidate) the session
	{#if sess?.identity}
		for user: <strong>{(sess?.identity?.traits as CustomTraits)?.['email']}</strong>
	{:else}
		<strong>{sess.id}</strong>
	{/if}
	so they will have to login again.
{/snippet}

{#snippet extendSessionDescriptionSnippet()}
	{@const sess = row.original as components['schemas']['Session']}
	This will extend the session
	{#if sess?.identity}
		for user: <strong>{(sess?.identity?.traits as CustomTraits)?.['email']}</strong>
	{:else}
		<strong>{sess.id}</strong>
	{/if}
	so they will have to login again.
{/snippet}

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="ghost" class="data-[state=open]:bg-muted flex h-8 w-8 p-0">
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
		<a href="/dashboard/sessions/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
		<DropdownMenu.Item onclick={() => onDeactivateSessionClick(row)} disabled={!row.getValue('active')}>
			<IconTrash2 />
			Deactivate (Invalidate)
		</DropdownMenu.Item>
		<DropdownMenu.Item onclick={() => onExtendSessionClick(row)} disabled={!row.getValue('active')}>
			<IconCirclePlus />
			Extend
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
