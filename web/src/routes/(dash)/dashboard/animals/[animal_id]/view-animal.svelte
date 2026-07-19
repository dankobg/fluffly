<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import IconTrash from '@lucide/svelte/icons/trash';
	import IconCircleCheck from '@lucide/svelte/icons/circle-check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import type { PageProps } from './$types';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { getErrorFallbackAnimalImage } from '$lib/utils';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';
	import { goto, invalidate } from '$app/navigation';
	import { AnimalStatus } from '$lib/gen/fluffly_openapi';

	let { data }: PageProps = $props();
	const timeFmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});
	const listFmt = new Intl.ListFormat(undefined, {
		style: 'long',
		type: 'conjunction'
	});

	type GroupedFields = Record<string, FieldInfo[]>;
	type FieldInfo = {
		key: string;
		title: string;
		value: unknown;
		type: string;
		enum?: string[];
	};

	let groupedFields: GroupedFields = {};

	for (const [key, val] of Object.entries(data?.animalSpeciesResult?.data?.properties_schema?.['properties'] ?? {})) {
		const group = val['x-group'] || 'Other';
		if (!groupedFields[group]) {
			groupedFields[group] = [];
		}

		groupedFields[group].push({
			key,
			title: val.title || key,
			value: data?.animalResult?.data?.properties?.[key],
			type: val.type,
			enum: val.enum
		});
	}

	const groupOrder = ['Traits', 'Behaviour', 'Health'];

	async function approveAnimal(id: number) {
		try {
			const approveAnimalResult = await fluffly.POST('/animals/{id}/submissions/approve', {
				params: {
					path: { id }
				}
			});
			if (approveAnimalResult.error) {
				toast.error([approveAnimalResult.error.message, approveAnimalResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Animal approved');
			invalidate(`data:dashboard-animals-${id}`);
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

	async function rejectAnimal(id: number) {
		try {
			const rejectAnimalResult = await fluffly.POST('/animals/{id}/submissions/reject', {
				params: {
					path: { id }
				}
			});
			if (rejectAnimalResult.error) {
				toast.error([rejectAnimalResult.error.message, rejectAnimalResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Animal rejected');
			invalidate(`data:dashboard-animals-${id}`);
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

	async function deleteAnimal(id: number) {
		try {
			const deleteAnimalResult = await fluffly.DELETE('/animals/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteAnimalResult.error) {
				toast.error([deleteAnimalResult.error.message, deleteAnimalResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Animal deleted');
			goto('/dashboard/animals');
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
			invalidate(`data:dashboard-animals-${animalId}`);
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
			invalidate(`data:dashboard-animals-${animalId}`);
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

{#if data.animalResult?.data}
	<div class="flex gap-x-4">
		<a href="/dashboard/animals/{data.animalResult.data.id}/update">
			<Button>Update animal</Button>
		</a>
		<Button
			variant="destructive"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete animal: ${data.animalResult.data!.name}?`,
					destructive: true,
					async onConfirm() {
						await deleteAnimal(data.animalResult.data!.id);
					}
				});
			}}
		>
			Delete animal
			<IconTrash />
		</Button>

		{#if data?.animalResult?.data?.status === AnimalStatus.pending}
			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Approve animal: ${data.animalResult.data!.name}?`,
						destructive: true,
						async onConfirm() {
							await approveAnimal(data.animalResult.data!.id);
						}
					});
				}}
			>
				Approve animal
				<IconCircleCheck />
			</Button>

			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Reject animal: ${data.animalResult.data!.name}?`,
						destructive: true,
						async onConfirm() {
							await rejectAnimal(data.animalResult.data!.id);
						}
					});
				}}
			>
				Reject animal
				<IconCircleX />
			</Button>
		{/if}

		{#if data?.animalResult?.data?.status === AnimalStatus.reserved}
			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Approve adoption: ${data.animalResult.data!.name}?`,
						destructive: true,
						async onConfirm() {
							await approveAdoption(data.animalResult.data!.id, data.animalResult.data!.adoption_id!);
						}
					});
				}}
			>
				Approve adoption
				<IconCircleCheck />
			</Button>

			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Reject adoption: ${data.animalResult.data!.name}?`,
						destructive: true,
						async onConfirm() {
							await rejectAdoption(data.animalResult.data!.id, data.animalResult.data!.adoption_id!);
						}
					});
				}}
			>
				Reject adoption
				<IconCircleX />
			</Button>
		{/if}
	</div>

	<h1 class="mb-6 text-2xl font-bold">Animal</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.animalResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.animalResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Type</span>
			<span class="font-medium">{data.animalResult?.data.type?.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Species</span>
			<span class="font-medium">{data.animalResult?.data.specie?.name}</span>
		</div>
		{#if data.animalResult?.data?.breeds}
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Breeds</span>
				<span class="font-medium">{listFmt.format(data.animalResult?.data.breeds?.map(x => x.name))}</span>
			</div>
		{/if}
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Gender</span>
			<span class="font-medium">{data.animalResult?.data.gender}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Hermaphrodite</span>
			<span class="font-medium">{data.animalResult?.data.hermaphrodite}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Age</span>
			<span class="font-medium">{data.animalResult?.data.age}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Size</span>
			<span class="font-medium">{data.animalResult?.data.size}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Status</span>
			<span class="font-medium">{data.animalResult?.data.status}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Description</span>
			<span class="font-medium">{data.animalResult?.data.description}</span>
		</div>
		{#if data.animalResult?.data?.tags}
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Tags</span>
				<span class="font-medium">{listFmt.format(data.animalResult?.data?.tags?.map(x => x.name))}</span>
			</div>
		{/if}
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{timeFmt.format(new Date(data.animalResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{timeFmt.format(new Date(data.animalResult?.data.updated_at))}</time>
		</div>
	</div>

	{#if data.animalResult.data?.microchip}
		<h2 class="my-6 text-xl font-bold">Microchip information</h2>
		<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Number</span>
				<span class="font-medium">{data.animalResult?.data.microchip.number}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Brand</span>
				<span class="font-medium">{data.animalResult?.data.microchip.brand}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Location</span>
				<span class="font-medium">{data.animalResult?.data.microchip.location}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Description</span>
				<span class="font-medium">{data.animalResult?.data.microchip.description}</span>
			</div>
		</div>
	{/if}

	<h2 class="my-6 text-xl font-bold">Properties</h2>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		{#each groupOrder as group (group)}
			{#if groupedFields[group]}
				<div>
					<h3 class="text-l my-6 font-bold">{group}</h3>
					<ul>
						{#each groupedFields[group] as field (field.title)}
							<li>
								<strong>{field.title}:</strong>
								{#if field.type === 'boolean'}
									{field.value ? 'Yes' : 'No'}
								{:else if field.enum}
									{field.value}
								{:else}
									{field.value}
								{/if}
							</li>
						{/each}
					</ul>
				</div>
			{/if}
		{/each}

		{#each Object.keys(groupedFields).filter(g => !groupOrder.includes(g)) as group (group)}
			<div>
				<h3 class="text-l my-6 font-bold">{group}</h3>
				<ul>
					{#each groupedFields[group] as field (field.title)}
						{#if field.value}
							<li>
								<strong>{field.title}:</strong>
								{#if field.type === 'boolean'}
									{field.value ? 'Yes' : 'No'}
								{:else if field.enum}
									{field.value}
								{:else}
									{field.value}
								{/if}
							</li>
						{/if}
					{/each}
				</ul>
			</div>
		{/each}
	</div>

	{#if data.animalResult.data?.image_full_url}
		<div class="mt-8 grid grid-cols-1 text-sm">
			<span class="mb-2 text-muted-foreground">Main photo</span>
			<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
				<div class="group grid h-64 overflow-hidden">
					<img
						src={data.animalResult.data?.image_full_url}
						alt="animal main"
						class="h-full w-full object-cover object-center"
						loading="lazy"
						onerror={e => {
							const url = getErrorFallbackAnimalImage(data?.animalResult?.data?.type.name ?? '');
							(e.currentTarget as HTMLImageElement).src = url;
						}}
					/>
				</div>
			</div>
		</div>
	{/if}

	{#if data.animalResult.data?.photos && data.animalResult.data?.photos?.length > 0}
		<div class="mt-8 grid grid-cols-1 text-sm">
			<span class="mb-2 text-muted-foreground">Photos</span>
			<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
				{#each data.animalResult.data.photos as photo, i (photo.id)}
					<div class="group grid h-64 overflow-hidden">
						<img
							src={photo.full_url}
							alt="animal photo {i + 1}"
							class="h-full w-full object-cover object-center"
							loading="lazy"
						/>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	{#if data.animalResult.data?.videos && data.animalResult.data?.videos?.length > 0}
		<div class="mt-8 grid grid-cols-1 text-sm">
			<span class="mb-2 text-muted-foreground">Videos</span>
			<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
				{#each data.animalResult.data.videos as video (video.id)}
					<div class="group grid h-64 overflow-hidden">
						<video controls class="h-full w-full object-cover" muted>
							<track kind="captions" />
							<source src={video.url} />
						</video>
					</div>
				{/each}
			</div>
		</div>
	{/if}
{/if}
