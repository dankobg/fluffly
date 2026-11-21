<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableFacetedFilter from '#lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import { confirmation } from '#lib/components/confirmation-dialog/confirmation-dialog-state.svelte.ts';
	import { page } from '$app/state';
	import { useDebounce } from 'runed';
	import { goto, invalidate } from '$app/navigation';
	import { statuses } from './data';
	import { organizationStatusValues } from '#lib/enum-values.ts';
	import type { OrganizationStatus } from '#lib/gen/fluffly_openapi.ts';
	import { onMount } from 'svelte';
	import { fluffly } from '#lib/fluffly/client.ts';
	import { toast } from 'svelte-sonner';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import FilterInputServer from '#lib/components/data-grid-shared/filter-input-server.svelte';

	const table = useTableContext();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	let nameCol = $derived(table.getColumn('name'));
	let statusCol = $derived(table.getColumn('status'));

	let name = $derived(page.url.searchParams.get('name') ?? '');

	let organizationStatuses = $derived.by(() => {
		const params = page.url.searchParams.getAll('status');
		if (params.length > 0) {
			const statuses = params.filter(x =>
				organizationStatusValues.includes(x as OrganizationStatus)
			) as OrganizationStatus[];
			if (statuses.length > 0) {
				return statuses;
			}
		}
	});

	const debounceNameFilter = useDebounce(
		() => {
			const sp = new URLSearchParams(new URL(page.url.href).searchParams);
			if (name) {
				sp.set('name', name);
			} else {
				sp.delete('name');
			}
			table.getColumn('name')?.setFilterValue(name ?? '');
			gotoWithFilters(sp);
		},
		() => 700
	);

	function gotoWithFilters(params: URLSearchParams) {
		goto(page.url.pathname + params.size ? `?${params}` : '', { reset: false });
	}

	onMount(() => {
		if (nameCol && name) {
			nameCol.setFilterValue(name);
		}
		if (statusCol && organizationStatuses) {
			statusCol.setFilterValue(organizationStatuses);
		}
	});

	async function deleteOrganizations(ids: number[]) {
		try {
			const deleteOrganizationsResult = await fluffly.DELETE('/organizations', {
				body: { ids }
			});
			if (deleteOrganizationsResult.error) {
				toast.error(
					[deleteOrganizationsResult.error.message, deleteOrganizationsResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Organizations deleted');
			invalidate('data:dashboard-organizations');
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

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if nameCol}
			<FilterInputServer column={nameCol} bind:value={name} onFilterChange={debounceNameFilter} />
		{/if}

		{#if statusCol}
			<DataTableFacetedFilter
				column={statusCol}
				title="Status"
				options={statuses}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('status', opt.value);
					} else {
						sp.delete('status', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if isFiltered}
			<Button
				variant="ghost"
				onclick={() => {
					table.resetColumnFilters();
					gotoWithFilters(new URLSearchParams());
				}}
				class="h-8 px-2 lg:px-3"
			>
				Reset
				<IconX />
			</Button>
		{/if}

		{#if table.getFilteredSelectedRowModel().rows.length > 0}
			<Button
				size="sm"
				variant="destructive"
				onclick={() => {
					const ids: number[] = table.getFilteredSelectedRowModel().rows.map(row => row.getValue('id'));

					confirmation.openDialog({
						title: `Delete ${table.getFilteredSelectedRowModel().rows.length} organizations?`,
						destructive: true,
						async onConfirm() {
							await deleteOrganizations(ids).finally(() => {
								table.resetRowSelection();
							});
						}
					});
				}}
			>
				Delete {table.getFilteredSelectedRowModel().rows.length} organizations
				<IconTrash />
			</Button>
		{/if}
	</div>

	<DataTableViewOptions />
</div>
