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
			goto('/dashboard/animal-species');
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

{#if data.animalSpeciesResult?.data}
	<div class="flex gap-x-4">
		<a href="/dashboard/animal-species/{data.animalSpeciesResult.data.id}/update">
			<Button>Update animal species</Button>
		</a>
		<Button
			variant="destructive"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete animal specie: ${data.animalSpeciesResult.data!.name}?`,
					destructive: true,
					async onConfirm() {
						await deleteAnimalSpecie(data.animalSpeciesResult.data!.id);
					}
				});
			}}
		>
			Delete animal specie
			<IconTrash />
		</Button>
	</div>

	<h1 class="mb-6 text-2xl font-bold">Animal species</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.animalSpeciesResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Animal type ID</span>
			<span class="font-medium">{data.animalSpeciesResult?.data.animal_type_id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.animalSpeciesResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.animalSpeciesResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.animalSpeciesResult?.data.updated_at))}</time>
		</div>
	</div>

	<div class="flex flex-col justify-center">
		<span class="text-muted-foreground">Properties schema</span>
		<pre class="font-medium">{JSON.stringify(data.animalSpeciesResult?.data.properties_schema ?? {}, null, 2)}</pre>
	</div>
{/if}
