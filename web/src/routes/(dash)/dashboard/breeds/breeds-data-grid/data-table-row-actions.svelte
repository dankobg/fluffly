<script lang="ts">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconPen from '@lucide/svelte/icons/pen';
	import IconTrash from '@lucide/svelte/icons/trash';
	import IconCopy from '@lucide/svelte/icons/copy';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.ts';
	import Button from '#lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { confirmation } from '#lib/components/confirmation-dialog/confirmation-dialog-state.svelte.ts';
	import { invalidate } from '$app/navigation';
	import { fluffly } from '#lib/fluffly/client.ts';
	import { useCellContext } from '#lib/components/data-grid-shared/app-table.ts';

	const cell = useCellContext();
	const row = $derived(cell.row);

	const hasId = $derived(
		typeof row.original === 'object' && !!row.original && 'id' in row.original && typeof row.original.id === 'string'
	);

	function copyIdToClipboard() {
		try {
			navigator.clipboard.writeText(row.original.id).then(() => {
				toast.success('id coppied');
			});
		} catch (error) {
			if (Error.isError(error)) {
				toast.error('failed to copy id: ' + error.message);
			}
		}
	}

	async function deleteBreed(id: number) {
		try {
			const deleteBreedResult = await fluffly.DELETE('/animal_breeds/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteBreedResult.error) {
				toast.error([deleteBreedResult.error.message, deleteBreedResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Breed deleted');
			invalidate('data:dashboard-breeds');
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
		<a href="/dashboard/breeds/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
		<a href="/dashboard/breeds/{row.getValue('id')}/update">
			<DropdownMenu.Item class="cursor-pointer">
				<IconPen />
				Update
			</DropdownMenu.Item>
		</a>
		<DropdownMenu.Item
			class="cursor-pointer"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete breed: ${row.getValue('name')}?`,
					destructive: true,
					async onConfirm() {
						await deleteBreed(row.getValue('id'));
					}
				});
			}}
		>
			<IconTrash />
			Delete
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
