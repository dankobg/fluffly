<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import IconEye from '@lucide/svelte/icons/eye';
	import IconCircleCheck from '@lucide/svelte/icons/circle-check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import IconCopy from '@lucide/svelte/icons/copy';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.ts';
	import Button from '#lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { confirmation } from '#lib/components/confirmation-dialog/confirmation-dialog-state.svelte.ts';
	import { invalidate } from '$app/navigation';
	import { fluffly } from '#lib/fluffly/client.ts';
	import { AdoptionStatus } from '#lib/gen/fluffly_openapi.ts';
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
	async function approveAdoption(animalId: number, adoptionId: number) {
		try {
			const approveAdoptionResult = await fluffly.POST('/animals/{id}/adoptions/{adoption_id}/approve', {
				params: {
					path: { id: animalId, adoption_id: adoptionId }
				}
			});
			if (approveAdoptionResult.error) {
				toast.error(
					[approveAdoptionResult.error.message, approveAdoptionResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Adoption approved');
			invalidate('data:dashboard-adoptions');
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

	async function rejectAdoption(animalId: number, adoptionId: number) {
		try {
			const rejectAdoptionResult = await fluffly.POST('/animals/{id}/adoptions/{adoption_id}/reject', {
				params: {
					path: { id: animalId, adoption_id: adoptionId }
				}
			});
			if (rejectAdoptionResult.error) {
				toast.error([rejectAdoptionResult.error.message, rejectAdoptionResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Adoption rejected');
			invalidate('data:dashboard-adoptions');
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
		<a href="/dashboard/countries/{row.getValue('id')}">
			<DropdownMenu.Item class="cursor-pointer">
				<IconEye />
				View
			</DropdownMenu.Item>
		</a>
		{#if row.getValue('status') === AdoptionStatus.pending}
			<DropdownMenu.Item
				class="cursor-pointer"
				onclick={() => {
					confirmation.openDialog({
						title: `Approve adoption: ${row.getValue('name')}?`,
						destructive: true,
						async onConfirm() {
							await approveAdoption(row.getValue('animal_id'), row.getValue('id'));
						}
					});
				}}
			>
				<IconCircleCheck />
				Approve adoption
			</DropdownMenu.Item>
			<DropdownMenu.Item
				class="cursor-pointer"
				onclick={() => {
					confirmation.openDialog({
						title: `Reject adoption: ${row.getValue('name')}?`,
						destructive: true,
						async onConfirm() {
							await rejectAdoption(row.getValue('animal_id'), row.getValue('id'));
						}
					});
				}}
			>
				<IconCircleX />
				Reject adoption
			</DropdownMenu.Item>
		{/if}
	</DropdownMenu.Content>
</DropdownMenu.Root>
