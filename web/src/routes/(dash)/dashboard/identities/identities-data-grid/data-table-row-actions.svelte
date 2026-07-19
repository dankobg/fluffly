<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconCopy from '@lucide/svelte/icons/copy';
	import IconPen from '@lucide/svelte/icons/pen';
	import IconFingerprint from '@lucide/svelte/icons/fingerprint';
	import IconTrash2 from '@lucide/svelte/icons/trash-2';
	import type { Row } from '@tanstack/table-core';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { invalidate } from '$app/navigation';
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

	async function onConfirmDeleteIdentity() {
		const identity = row.original as components['schemas']['Identity'];
		try {
			await fluffly.DELETE('/identities/{id}', {
				params: {
					path: { id: identity.id }
				}
			});
			toast.success('identity deleted');
			invalidate('data:dashboard-identities');
		} catch (error) {
			console.log('err', error);
			toast.error('identity delete failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeleteIdentityClick() {
		confirmation.openDialog({
			description: deleteIdentityDescriptionSnippet,
			onConfirm: onConfirmDeleteIdentity,
			destructive: true
		});
	}
</script>

{#snippet deleteIdentityDescriptionSnippet()}
	{@const identity = row?.original as components['schemas']['Identity']}
	{@const email = (identity?.traits as CustomTraits)?.['email']}
	This action cannot be undone. This will delete the identity <strong>{email}</strong> completely.
{/snippet}

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
		<a href="/dashboard/identities/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
		<a href="/dashboard/identities/{row.getValue('id')}/update">
			<DropdownMenu.Item class="cursor-pointer">
				<IconPen />
				Update
			</DropdownMenu.Item>
		</a>
		<a href="/dashboard/identities/{row.getValue('id')}/sessions">
			<DropdownMenu.Item class="cursor-pointer">
				<IconFingerprint />
				View sessions
			</DropdownMenu.Item>
		</a>
		<DropdownMenu.Item onclick={onDeleteIdentityClick}>
			<IconTrash2 />
			Delete
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
