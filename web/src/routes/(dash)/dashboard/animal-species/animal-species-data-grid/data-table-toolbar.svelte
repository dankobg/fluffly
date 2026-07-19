<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconX from '@lucide/svelte/icons/x';
	import type { Table } from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Input } from '$lib/components/ui/input/index';
	import DataTableViewOptions from '$lib/components/data-grid-shared/data-table-view-options.svelte';
	import IconTrash from '@lucide/svelte/icons/trash-2';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import type { components } from '$lib/gen/fluffly_openapi';
	import { useDebounce } from 'runed';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { capitalize } from '$lib/utils';
	import DataTableFacetedFilter from '$lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';

	let { table, animalTypes }: { table: Table<TData>; animalTypes: components['schemas']['AnimalType'][] } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	let nameCol = $derived(table.getColumn('name'));
	let animalTypeIdCol = $derived(table.getColumn('animal_type_id'));

	let animalTypeOptions = animalTypes.map(x => ({ label: capitalize(x.name), value: String(x.id) }));

	let animalTypeIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('animal_type_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});

	let name = $derived(page.url.searchParams.get('name') ?? '');

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
		if (animalTypeIdCol && animalTypeIds) {
			animalTypeIdCol.setFilterValue(animalTypeIds);
		}
		if (nameCol && name) {
			nameCol.setFilterValue(name);
		}
	});

	async function deleteAnimalSpecies(ids: number[]) {
		try {
			const deleteAnimalSpeciesResult = await fluffly.DELETE('/animal_species', {
				body: { ids }
			});
			if (deleteAnimalSpeciesResult.error) {
				toast.error(
					[deleteAnimalSpeciesResult.error.message, deleteAnimalSpeciesResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Animal species deleted');
			invalidate('data:dashboard-animal-species');
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

		{#if animalTypeIdCol}
			<DataTableFacetedFilter
				column={animalTypeIdCol}
				title="Animal type"
				options={animalTypeOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(page.url.searchParams);
					if (selected) {
						sp.append('animal_type_id', opt.value);
					} else {
						sp.delete('animal_type_id', opt.value);
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
						title: `Delete ${table.getFilteredSelectedRowModel().rows.length} animal types?`,
						destructive: true,
						async onConfirm() {
							await deleteAnimalSpecies(ids).finally(() => {
								table.resetRowSelection();
							});
						}
					});
				}}
			>
				Delete {table.getFilteredSelectedRowModel().rows.length} species
				<IconTrash />
			</Button>
		{/if}
	</div>

	<DataTableViewOptions {table} />
</div>
