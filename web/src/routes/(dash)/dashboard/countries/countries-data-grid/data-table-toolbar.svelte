<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import { confirmation } from '#lib/components/confirmation-dialog/confirmation-dialog-state.svelte.ts';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';
	import { useDebounce } from 'runed';
	import { onMount } from 'svelte';
	import { fluffly } from '#lib/fluffly/client.ts';
	import { toast } from 'svelte-sonner';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import FilterInputServer from '#lib/components/data-grid-shared/filter-input-server.svelte';

	const table = useTableContext();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	let nameCol = $derived(table.getColumn('name'));
	let isoAlpha2Col = $derived(table.getColumn('iso_alpha2'));
	let isoAlpha3Col = $derived(table.getColumn('iso_alpha3'));

	let name = $derived(page.url.searchParams.get('name') ?? '');
	let isoAlpha2 = $derived(page.url.searchParams.get('iso_alpha2') ?? '');
	let isoAlpha3 = $derived(page.url.searchParams.get('iso_alpha3') ?? '');

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

	const debounceIsoAlpha2Filter = useDebounce(
		() => {
			const sp = new URLSearchParams(new URL(page.url.href).searchParams);
			if (isoAlpha2) {
				sp.set('iso_alpha2', isoAlpha2);
			} else {
				sp.delete('iso_alpha2');
			}
			table.getColumn('iso_alpha2')?.setFilterValue(isoAlpha2 ?? '');
			gotoWithFilters(sp);
		},
		() => 700
	);

	const debounceIsoAlpha3Filter = useDebounce(
		() => {
			const sp = new URLSearchParams(new URL(page.url.href).searchParams);
			if (isoAlpha3) {
				sp.set('iso_alpha3', isoAlpha3);
			} else {
				sp.delete('iso_alpha3');
			}
			table.getColumn('iso_alpha3')?.setFilterValue(isoAlpha3 ?? '');
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
		if (isoAlpha2Col && isoAlpha2) {
			isoAlpha2Col.setFilterValue(isoAlpha2);
		}
		if (isoAlpha3Col && isoAlpha3) {
			isoAlpha3Col.setFilterValue(isoAlpha3);
		}
	});

	async function deleteCountries(ids: number[]) {
		try {
			const deleteCountriesResult = await fluffly.DELETE('/countries', {
				body: { ids }
			});
			if (deleteCountriesResult.error) {
				toast.error(
					[deleteCountriesResult.error.message, deleteCountriesResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Countries deleted');
			invalidate('data:dashboard-countries');
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

		{#if isoAlpha2Col}
			<FilterInputServer column={isoAlpha2Col} bind:value={isoAlpha2} onFilterChange={debounceIsoAlpha2Filter} />
		{/if}

		{#if isoAlpha3Col}
			<FilterInputServer column={isoAlpha3Col} bind:value={isoAlpha3} onFilterChange={debounceIsoAlpha3Filter} />
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
						title: `Delete ${table.getFilteredSelectedRowModel().rows.length} countries?`,
						destructive: true,
						async onConfirm() {
							await deleteCountries(ids).finally(() => {
								table.resetRowSelection();
							});
						}
					});
				}}
			>
				Delete {table.getFilteredSelectedRowModel().rows.length} countries
				<IconTrash />
			</Button>
		{/if}
	</div>

	<DataTableViewOptions />
</div>
