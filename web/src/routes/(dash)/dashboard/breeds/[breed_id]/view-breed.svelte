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
			goto('/dashboard/breeds');
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

{#if data.breedResult?.data}
	<div class="flex gap-x-4">
		<a href="/dashboard/breeds/{data.breedResult.data.id}/update">
			<Button>Update breed</Button>
		</a>
		<Button
			variant="destructive"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete breed: ${data.breedResult.data!.name}?`,
					destructive: true,
					async onConfirm() {
						await deleteBreed(data.breedResult.data!.id);
					}
				});
			}}
		>
			Delete breed
			<IconTrash />
		</Button>
	</div>

	<h1 class="mb-6 text-2xl font-bold">Breed</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.breedResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Animal species ID</span>
			<span class="font-medium">{data.breedResult?.data.animal_specie_id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.breedResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.breedResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.breedResult?.data.updated_at))}</time>
		</div>
	</div>
{/if}
