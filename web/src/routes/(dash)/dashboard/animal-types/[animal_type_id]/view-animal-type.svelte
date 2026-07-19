<script lang="ts">
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';
	import type { PageProps } from './$types';
	import IconTrash from '@lucide/svelte/icons/trash';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	async function deleteAnimalType(id: number) {
		try {
			const deleteAnimalTypeResult = await fluffly.DELETE('/animal_types/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteAnimalTypeResult.error) {
				toast.error(
					[deleteAnimalTypeResult.error.message, deleteAnimalTypeResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Animal type deleted');
			goto('/dashboard/animal-types');
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

{#if data.animalTypeResult?.data}
	<div class="flex gap-x-4">
		<a href="/dashboard/animal-types/{data.animalTypeResult.data.id}/update">
			<Button>Update animal type</Button>
		</a>
		<Button
			variant="destructive"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete animal type: ${data.animalTypeResult.data!.name}?`,
					destructive: true,
					async onConfirm() {
						await deleteAnimalType(data.animalTypeResult.data!.id);
					}
				});
			}}
		>
			Delete animal type
			<IconTrash />
		</Button>
	</div>

	<h1 class="mb-6 text-2xl font-bold">Animal type</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.animalTypeResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.animalTypeResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.animalTypeResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.animalTypeResult?.data.updated_at))}</time>
		</div>
	</div>
{/if}
