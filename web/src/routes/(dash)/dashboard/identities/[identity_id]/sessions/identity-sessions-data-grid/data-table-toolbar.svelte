<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableFacetedFilter from '#lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import { aals } from './data';
	import { confirmation } from '#lib/components/confirmation-dialog/confirmation-dialog-state.svelte.ts';
	import { fluffly } from '#lib/fluffly/client.ts';
	import { toast } from 'svelte-sonner';
	import { invalidate } from '$app/navigation';
	import type { components } from '#lib/gen/fluffly_openapi.ts';
	import type { CustomTraits } from '#lib/kratos/service.ts';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';

	const table = useTableContext();
	let { identity }: { identity: components['schemas']['Identity'] | undefined } = $props();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	const aalCol = $derived(table.getColumn('authenticator_assurance_level'));

	async function onConfirmDeleteIdentitySessions() {
		if (!identity) {
			return;
		}
		try {
			await fluffly.DELETE('/identities/{id}/sessions', {
				params: {
					path: { id: identity.id }
				}
			});
			toast.success('sessions deleted');
			invalidate(`data:dashboard-identities-${identity.id}-sessions`);
		} catch (error) {
			console.log('err', error);
			toast.error('sessions delete failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeleteIdentitySessionsClick() {
		confirmation.openDialog({
			description: deleteIdentitySessionsDescriptionSnippet,
			onConfirm: onConfirmDeleteIdentitySessions
		});
	}
</script>

{#snippet deleteIdentitySessionsDescriptionSnippet()}
	This action cannot be undone. This will delete all sessions for user: <strong>
		{(identity?.traits as CustomTraits)?.['email']}
	</strong>
{/snippet}

<div class="flex gap-4">
	<Button onclick={onDeleteIdentitySessionsClick} disabled={table.getRowModel().rows.length === 0}>
		Delete all sessions
	</Button>
</div>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if aalCol}
			<DataTableFacetedFilter column={aalCol} title="Aal" options={aals} />
		{/if}

		{#if isFiltered}
			<Button variant="ghost" onclick={() => table.resetColumnFilters()} class="h-8 px-2 lg:px-3">
				Reset
				<IconX />
			</Button>
		{/if}
	</div>
	<DataTableViewOptions />
</div>
