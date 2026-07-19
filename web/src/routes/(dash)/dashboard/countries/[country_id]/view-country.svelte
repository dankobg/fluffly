<script lang="ts">
	import { goto } from '$app/navigation';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';
	import type { PageProps } from './$types';
	import IconTrash from '@lucide/svelte/icons/trash';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	async function deleteCountry(id: number) {
		try {
			const deleteCountryResult = await fluffly.DELETE('/countries/{id}', {
				params: {
					path: { id }
				}
			});
			if (deleteCountryResult.error) {
				toast.error([deleteCountryResult.error.message, deleteCountryResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Country deleted');
			goto('/dashboard/countries');
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

{#if data.countryResult?.data}
	<div class="flex gap-x-4">
		<a href="/dashboard/countries/{data.countryResult.data.id}/update">
			<Button>Update country</Button>
		</a>
		<Button
			variant="destructive"
			onclick={() => {
				confirmation.openDialog({
					title: `Delete country: ${data.countryResult.data!.name}?`,
					destructive: true,
					async onConfirm() {
						await deleteCountry(data.countryResult.data!.id);
					}
				});
			}}
		>
			Delete country
			<IconTrash />
		</Button>
	</div>

	<h1 class="mb-6 text-2xl font-bold">Country</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.countryResult?.data.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Name</span>
			<span class="font-medium">{data.countryResult?.data.name}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Iso alpha 2</span>
			<span class="font-medium">{data.countryResult?.data.iso_alpha2}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Iso alpha 3</span>
			<span class="font-medium">{data.countryResult?.data.iso_alpha3}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Iso numeric</span>
			<span class="font-medium">{data.countryResult?.data.iso_numeric}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(new Date(data.countryResult?.data.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(new Date(data.countryResult?.data.updated_at))}</time>
		</div>
	</div>
{/if}
