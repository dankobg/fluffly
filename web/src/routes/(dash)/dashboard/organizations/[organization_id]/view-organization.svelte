<script lang="ts">
	const backendBaseUrl = import.meta.env['VITE_PUBLIC_BACKEND_BASE'] as string;
	import Button from '$lib/components/ui/button/button.svelte';
	import IconTrash from '@lucide/svelte/icons/trash';
	import IconCircleCheck from '@lucide/svelte/icons/circle-check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import type { PageProps } from './$types';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import * as Item from '$lib/components/ui/item/index.js';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';
	import { goto, invalidate } from '$app/navigation';
	import { OrganizationStatus } from '$lib/gen/fluffly_openapi';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	async function approveOrganization(id: number) {
		try {
			const approveOrganizationResult = await fluffly.POST('/organizations/{id}/application/approve', {
				params: {
					path: { id }
				}
			});
			if (approveOrganizationResult.error) {
				toast.error(
					[approveOrganizationResult.error.message, approveOrganizationResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organization approved');
			invalidate(`data:dashboard-organizations-${id}`);
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

	async function rejectOrganization(id: number) {
		try {
			const rejectOrganizationResult = await fluffly.POST('/organizations/{id}/application/reject', {
				params: {
					path: { id }
				}
			});
			if (rejectOrganizationResult.error) {
				toast.error(
					[rejectOrganizationResult.error.message, rejectOrganizationResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organization rejected');
			invalidate(`data:dashboard-organizations-${id}`);
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

	async function deleteOrganization(id: number) {
		try {
			const deleteOrganizationResult = await fluffly.DELETE('/organizations/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteOrganizationResult.error) {
				toast.error(
					[deleteOrganizationResult.error.message, deleteOrganizationResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organization deleted');
			goto('/dashboard/organizations');
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

{#if data.organizationResult?.data}
	<div class="flex gap-x-4">
		<a href="/dashboard/organizations/{data.organizationResult.data.id}/update">
			<Button>Update organization</Button>
		</a>
		<Button
			variant="destructive"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete organization: ${data.organizationResult.data!.name}?`,
					destructive: true,
					async onConfirm() {
						await deleteOrganization(data.organizationResult.data!.id);
					}
				});
			}}
		>
			Delete organization
			<IconTrash />
		</Button>

		{#if data?.organizationResult?.data?.status === OrganizationStatus.pending}
			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Approve organization: ${data.organizationResult.data!.name}?`,
						destructive: true,
						async onConfirm() {
							await approveOrganization(data.organizationResult.data!.id);
						}
					});
				}}
			>
				Approve organization
				<IconCircleCheck />
			</Button>

			<Button
				onclick={() => {
					confirmation.openDialog({
						title: `Reject organization: ${data.organizationResult.data!.name}?`,
						destructive: true,
						async onConfirm() {
							await rejectOrganization(data.organizationResult.data!.id);
						}
					});
				}}
			>
				Reject organization
				<IconCircleX />
			</Button>
		{/if}
	</div>

	<h1 class="mb-6 text-2xl font-bold">Organization</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.organizationResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.organizationResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Status</span>
			<span class="font-medium">{data.organizationResult?.data.status}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Website</span>
			<span class="font-medium">{data.organizationResult?.data.website}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Mission statement</span>
			<span class="font-medium">{data.organizationResult?.data.mission_statement}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Adoption policy</span>
			<span class="font-medium">{data.organizationResult?.data.adoption_policy}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Adoption URL</span>
			<span class="font-medium">{data.organizationResult?.data.adoption_url}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.organizationResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.organizationResult?.data.updated_at))}</time>
		</div>
	</div>

	{#if data.organizationResult?.data?.contact}
		<h2 class="my-6 text-xl font-bold">Contact information</h2>
		<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Phone</span>
				<span class="font-medium">{data.organizationResult?.data.contact.phone}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">E-Mail</span>
				<span class="font-medium">{data.organizationResult?.data.contact.email}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Country</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.country.name}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">City</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.city}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Region</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.region}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Street address</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.street_address}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Street number</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.street_number}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Unit number</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.unit_number}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Postal code</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.postal_code}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Latitude/Longitude</span>
				<span class="font-medium"
					>{data.organizationResult?.data.contact.address.lat}, {data.organizationResult?.data.contact.address
						.lon}</span
				>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Note</span>
				<span class="font-medium">{data.organizationResult?.data.contact.address.note}</span>
			</div>
		</div>
	{/if}

	{#if data.organizationResult.data?.socials?.length && data.organizationResult.data?.socials?.length > 0}
		<div class="grid grid-cols-1 text-sm">
			<span class="mb-2 text-muted-foreground">Social platforms</span>
			<Item.Group>
				{#each data.organizationResult.data?.socials as social (social.id + social.platform)}
					<Item.Root class="p-2">
						<Item.Content class="gap-0">
							<Item.Title>{social.platform}</Item.Title>
							<Item.Description>
								<a href={social.url}>{social.url}</a>
							</Item.Description>
						</Item.Content>
					</Item.Root>
				{/each}
			</Item.Group>
		</div>
	{/if}

	{#if data.organizationResult.data?.work_hour}
		{#if Object.values(data.organizationResult.data.work_hour).length > 0}
			<div class="grid grid-cols-1 text-sm">
				<span class="mb-2 text-muted-foreground">Work hours</span>
				<Item.Group>
					{#each ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'] as const as day (day)}
						{#if data.organizationResult.data?.work_hour?.[day]}
							<Item.Root class="p-1">
								<Item.Content class="flex flex-row gap-2">
									<Item.Title class="text-muted-foreground">{day[0]?.toUpperCase()}{day.slice(1)}:</Item.Title>
									<Item.Description class="text-white">
										{data.organizationResult.data?.work_hour?.[day]}
									</Item.Description>
								</Item.Content>
							</Item.Root>
						{/if}
					{/each}
				</Item.Group>
			</div>
		{/if}
	{/if}

	{#if data.organizationResult.data?.photos && data.organizationResult.data?.photos?.length > 0}
		<div class="mt-8 grid grid-cols-1 text-sm">
			<span class="mb-2 text-muted-foreground">Photos</span>
			<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
				{#each data.organizationResult.data.photos as photo, i (photo.id)}
					<div class="group grid h-64 overflow-hidden">
						<img
							src={photo.full_url}
							alt="organization photo {i + 1}"
							class="h-full w-full object-cover object-center"
							loading="lazy"
							onerror={e => {
								const num = Math.floor(Math.random() * 5) + 1;
								(e.currentTarget as HTMLImageElement).src =
									`${backendBaseUrl}/public/images/placeholder/shelter-${num}.svg`;
							}}
						/>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	{#if data.organizationResult.data?.videos && data.organizationResult.data?.videos?.length > 0}
		<div class="mt-8 grid grid-cols-1 text-sm">
			<span class="mb-2 text-muted-foreground">Videos</span>
			<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
				{#each data.organizationResult.data.videos as video (video.id)}
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
