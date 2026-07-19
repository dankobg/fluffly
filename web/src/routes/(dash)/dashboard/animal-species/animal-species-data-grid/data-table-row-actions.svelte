<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconPen from '@lucide/svelte/icons/pen';
	import IconTrash from '@lucide/svelte/icons/trash';
	import IconCopy from '@lucide/svelte/icons/copy';
	import type { Row } from '@tanstack/table-core';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { invalidate } from '$app/navigation';
	import { fluffly } from '$lib/fluffly/client';

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

	async function deleteAnimalSpecie(id: number) {
		try {
			const deleteAnimalSpecieResult = await fluffly.DELETE('/animal_species/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteAnimalSpecieResult.error) {
				toast.error(
					[deleteAnimalSpecieResult.error.message, deleteAnimalSpecieResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Animal specie deleted');
			invalidate('data:dashboard-animal-species');
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
		<a href="/dashboard/animal-species/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
		<a href="/dashboard/animal-species/{row.getValue('id')}/update">
			<DropdownMenu.Item class="cursor-pointer">
				<IconPen />
				Update
			</DropdownMenu.Item>
		</a>
		<DropdownMenu.Item
			class="cursor-pointer"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete animal specie: ${row.getValue('name')}?`,
					destructive: true,
					async onConfirm() {
						await deleteAnimalSpecie(row.getValue('id'));
					}
				});
			}}
		>
			<IconTrash />
			Delete
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
