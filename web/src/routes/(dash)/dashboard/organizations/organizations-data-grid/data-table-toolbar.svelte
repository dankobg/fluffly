<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconX from '@lucide/svelte/icons/x';
	import type { Table } from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Input } from '$lib/components/ui/input/index';
	import DataTableFacetedFilter from '$lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '$lib/components/data-grid-shared/data-table-view-options.svelte';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { page } from '$app/state';
	import { useDebounce } from 'runed';
	import { goto, invalidate } from '$app/navigation';
	import { statuses } from './data';
	import { organizationStatusValues } from '$lib/enum-values';
	import type { OrganizationStatus } from '$lib/gen/fluffly_openapi';
	import { onMount } from 'svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';

	let { table }: { table: Table<TData> } = $props();

	let isFiltered = $derived(table.getState().columnFilters.length > 0);
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
			const sp = new URLSearchParams(page.url.searchParams);
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
		goto(page.url.pathname + params.size ? `?${params}` : '', { keepFocus: true });
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
		<Input
			placeholder="Filter by name..."
			value={name}
			oninput={e => {
				name = e.currentTarget.value;
				debounceNameFilter();
			}}
			onchange={e => {
				name = e.currentTarget.value;
				debounceNameFilter();
			}}
			class="h-8 w-[150px] lg:w-[250px]"
		/>

		{#if statusCol}
			<DataTableFacetedFilter
				column={statusCol}
				title="Status"
				options={statuses}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(page.url.searchParams);
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

	<DataTableViewOptions {table} />
</div>
