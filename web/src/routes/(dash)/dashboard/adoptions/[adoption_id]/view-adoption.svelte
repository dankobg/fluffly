<script lang="ts">
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconCircleCheck from '@lucide/svelte/icons/circle-check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import IconClock from '@lucide/svelte/icons/clock';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';
	import type { PageProps } from './$types';
	import { invalidate } from '$app/navigation';
	import { AdoptionStatus, OrganizationStatus } from '$lib/gen/fluffly_openapi';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

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
			invalidate(`data:dashboard-adoptions-${adoptionId}`);
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
			invalidate(`data:dashboard-adoptions-${adoptionId}`);
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

{#if data.adoptionResult?.data}
	<div class="flex gap-x-4">
		{#if data?.adoptionResult?.data?.status === AdoptionStatus.pending}
			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Approve adoption: ${data.adoptionResult.data!.id}?`,
						destructive: true,
						async onConfirm() {
							await approveAdoption(data.adoptionResult.data!.animal_id, data.adoptionResult.data!.id!);
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
						title: `Reject adoption: ${data.adoptionResult.data!.id}?`,
						destructive: true,
						async onConfirm() {
							await rejectAdoption(data.adoptionResult.data!.animal_id, data.adoptionResult.data!.id!);
						}
					});
				}}
			>
				Reject adoption
				<IconCircleX />
			</Button>
		{/if}
	</div>

	<h1 class="mb-6 text-2xl font-bold">Adoption</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.adoptionResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">User ID</span>
			<span class="font-medium">{data.adoptionResult?.data.user_id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Organization ID</span>
			<span class="font-medium">{data.adoptionResult?.data.organization_id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Animal ID</span>
			<span class="font-medium">{data.adoptionResult?.data.animal_id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Status</span>
			<div class="flex gap-2">
				<span class="font-medium">{data.adoptionResult?.data.status}</span>
				{#if data.adoptionResult?.data.status === AdoptionStatus.pending}
					<IconClock class="text-yellow-600" />
				{:else if data.adoptionResult?.data.status === AdoptionStatus.approved}
					<IconCheck class="text-green-600" />
				{:else}
					<IconX class="text-red-600" />
				{/if}
			</div>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Permanent</span>
			<span class="font-medium">{data.adoptionResult?.data.is_permanent}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Adopted at</span>
			<time class="font-medium">
				{data.adoptionResult?.data.adopted_at && fmt.format(new Date(data.adoptionResult?.data.adopted_at))}
			</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Note</span>
			<span class="font-medium">{data.adoptionResult?.data.note}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Returned at</span>
			<time class="font-medium">
				{data.adoptionResult?.data.returned_at && fmt.format(new Date(data.adoptionResult?.data.returned_at))}
			</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.adoptionResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.adoptionResult?.data.updated_at))}</time>
		</div>
	</div>

	{#if data?.adoptionResult?.data?.organization}
		<h1 class="my-6 text-xl font-bold">Adopter organization</h1>
		<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">ID</span>
				<span class="font-medium">{data.adoptionResult?.data?.organization?.id}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Name</span>
				<span class="font-medium">{data.adoptionResult?.data?.organization?.name}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Status</span>
				<div class="flex gap-2">
					<span class="font-medium">{data.adoptionResult?.data?.organization?.status}</span>
					{#if data.adoptionResult?.data?.organization?.status === OrganizationStatus.pending}
						<IconClock class="text-yellow-600" />
					{:else if data.adoptionResult?.data?.organization?.status === OrganizationStatus.approved}
						<IconCheck class="text-green-600" />
					{:else}
						<IconX class="text-red-600" />
					{/if}
				</div>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Website</span>
				<span class="font-medium">{data.adoptionResult?.data?.organization?.website}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Mission statement</span>
				<span class="font-medium">{data.adoptionResult?.data?.organization?.mission_statement}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Adoption policy</span>
				<span class="font-medium">{data.adoptionResult?.data?.organization?.adoption_policy}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Adoption URL</span>
				<span class="font-medium">{data.adoptionResult?.data?.organization?.adoption_url}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Created time</span>
				<time class="font-medium">
					{data.adoptionResult?.data?.organization &&
						fmt.format(new Date(data.adoptionResult?.data?.organization?.created_at))}
				</time>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Updated time</span>
				<time class="font-medium">
					{data.adoptionResult?.data?.organization &&
						fmt.format(new Date(data.adoptionResult?.data?.organization?.updated_at))}
				</time>
			</div>
		</div>

		{#if data.adoptionResult?.data?.organization?.contact}
			<h2 class="my-6 text-xl font-bold">Organization contact information</h2>
			<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Phone</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.phone}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">E-Mail</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.email}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Country</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.country.name}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">City</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.city}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Region</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.region}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Street address</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.street_address}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Street number</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.street_number}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Unit number</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.unit_number}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Postal code</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.postal_code}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Latitude/Longitude</span>
					<span class="font-medium">
						{data.adoptionResult?.data?.organization?.contact.address.lat}, {data.adoptionResult?.data?.organization
							?.contact.address.lon}
					</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Note</span>
					<span class="font-medium">{data.adoptionResult?.data?.organization?.contact.address.note}</span>
				</div>
			</div>
		{/if}
	{/if}
{/if}
