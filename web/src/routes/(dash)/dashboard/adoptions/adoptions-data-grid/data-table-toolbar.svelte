<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconX from '@lucide/svelte/icons/x';
	import type { Table } from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import DataTableViewOptions from '$lib/components/data-grid-shared/data-table-view-options.svelte';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { adoptionStatusValues } from '$lib/enum-values';
	import type { AdoptionStatus } from '$lib/gen/fluffly_openapi';
	import DataTableFacetedFilter from '$lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import { statuses } from './data';

	let { table }: { table: Table<TData> } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	let statusCol = $derived(table.getColumn('status'));
	let animalIdCol = $derived(table.getColumn('animal_id'));
	let organizationIdCol = $derived(table.getColumn('organization_id'));
	let userIdCol = $derived(table.getColumn('user_id'));

	let adoptionStatuses = $derived.by(() => {
		const params = page.url.searchParams.getAll('status');
		if (params.length > 0) {
			const statuses = params.filter(x => adoptionStatusValues.includes(x as AdoptionStatus)) as AdoptionStatus[];
			if (statuses.length > 0) {
				return statuses;
			}
		}
	});

	let animalIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('animal_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});
	let organizationIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('organization_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});
	let userIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('user_id');
		if (params.length > 0) {
			return params;
		}
	});

	function gotoWithFilters(params: URLSearchParams) {
		goto(page.url.pathname + params.size ? `?${params}` : '', { keepFocus: true });
	}

	onMount(() => {
		if (statusCol && adoptionStatuses) {
			statusCol.setFilterValue(adoptionStatuses);
		}
		if (animalIdCol && animalIds) {
			animalIdCol.setFilterValue(animalIds);
		}
		if (organizationIdCol && organizationIds) {
			organizationIdCol.setFilterValue(organizationIds);
		}
		if (userIdCol && userIds) {
			userIdCol.setFilterValue(userIds);
		}
	});
</script>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
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
	</div>

	<DataTableViewOptions {table} />
</div>
