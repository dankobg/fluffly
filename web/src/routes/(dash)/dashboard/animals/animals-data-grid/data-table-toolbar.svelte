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
	import { page } from '$app/state';
	import { useDebounce } from 'runed';
	import { goto, invalidate } from '$app/navigation';
	import { statuses } from './data';
	import type { AnimalAge, AnimalSize, AnimalGender, AnimalStatus, components } from '$lib/gen/fluffly_openapi';
	import { capitalize } from '$lib/utils';
	import { animalAgeValues, animalGenderValues, animalSizeValues, animalStatusValues } from '$lib/enum-values';
	import { onMount } from 'svelte';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';

	let {
		table,
		animalTypes,
		animalSpecies
	}: {
		table: Table<TData>;
		animalTypes: components['schemas']['AnimalType'][];
		animalSpecies: components['schemas']['AnimalSpecie'][];
	} = $props();

	let isFiltered = $derived(table.getState().columnFilters.length > 0);
	let nameCol = $derived(table.getColumn('name'));
	let typeCol = $derived(table.getColumn('type'));
	let speciesCol = $derived(table.getColumn('specie'));
	let genderCol = $derived(table.getColumn('gender'));
	let ageCol = $derived(table.getColumn('age'));
	let sizeCol = $derived(table.getColumn('size'));
	let statusCol = $derived(table.getColumn('status'));

	const ageOptions = animalAgeValues.map(x => ({
		value: x,
		label: capitalize(x)
	}));

	const sizeOptions = animalSizeValues.map(x => ({
		value: x,
		label: capitalize(x)
	}));

	const genderOptions = animalGenderValues.map(x => ({
		value: x,
		label: x === 'm' ? 'Male' : 'Female'
	}));

	let animalTypeOptions = animalTypes.map(x => ({ label: capitalize(x.name), value: String(x.id) }));
	let animalSpeciesOptions = animalSpecies.map(x => ({ label: capitalize(x.name), value: String(x.id) }));

	let animalTypeIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('animal_type_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});
	let animalSpecieIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('animal_specie_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});
	let animalAges = $derived.by(() => {
		const params = page.url.searchParams.getAll('age');
		if (params.length > 0) {
			const ages = params.filter(x => animalAgeValues.includes(x as AnimalAge)) as AnimalAge[];
			if (ages.length > 0) {
				return ages;
			}
		}
	});
	let animalSizes = $derived.by(() => {
		const params = page.url.searchParams.getAll('size');
		if (params.length > 0) {
			const sizes = params.filter(x => animalSizeValues.includes(x as AnimalSize)) as AnimalSize[];
			if (sizes.length > 0) {
				return sizes;
			}
		}
	});
	let animalGenders = $derived.by(() => {
		const params = page.url.searchParams.getAll('gender');
		if (params.length > 0) {
			const genders = params.filter(x => animalGenderValues.includes(x as AnimalGender)) as AnimalGender[];
			if (genders.length > 0) {
				return genders;
			}
		}
	});
	let animalStatuses = $derived.by(() => {
		const params = page.url.searchParams.getAll('status');
		if (params.length > 0) {
			const statuses = params.filter(x => animalStatusValues.includes(x as AnimalStatus)) as AnimalStatus[];
			if (statuses.length > 0) {
				return statuses;
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
		if (typeCol && animalTypeIds) {
			typeCol.setFilterValue(animalTypeIds);
		}
		if (speciesCol && animalSpecieIds) {
			speciesCol.setFilterValue(animalSpecieIds);
		}
		if (nameCol && name) {
			nameCol.setFilterValue(name);
		}
		if (ageCol && animalAges) {
			ageCol.setFilterValue(animalAges);
		}
		if (sizeCol && animalSizes) {
			sizeCol.setFilterValue(animalSizes);
		}
		if (genderCol && animalGenders) {
			genderCol.setFilterValue(animalGenders);
		}
		if (statusCol && animalStatuses) {
			statusCol.setFilterValue(animalStatuses);
		}
	});

	async function deleteAnimals(ids: number[]) {
		try {
			const deleteAnimalsResult = await fluffly.DELETE('/animals', {
				body: { ids }
			});
			if (deleteAnimalsResult.error) {
				toast.error([deleteAnimalsResult.error.message, deleteAnimalsResult.error.reason].filter(Boolean).join(', '));
				return;
			}
			toast.success('Animals deleted');
			invalidate('data:dashboard-animals');
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

		{#if typeCol}
			<DataTableFacetedFilter
				column={typeCol}
				title="Type"
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

		{#if speciesCol}
			<DataTableFacetedFilter
				column={speciesCol}
				title="Species"
				options={animalSpeciesOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(page.url.searchParams);
					if (selected) {
						sp.append('animal_specie_id', opt.value);
					} else {
						sp.delete('animal_specie_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if genderCol}
			<DataTableFacetedFilter
				column={genderCol}
				title="Gender"
				options={genderOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(page.url.searchParams);
					if (selected) {
						sp.append('gender', opt.value);
					} else {
						sp.delete('gender', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if ageCol}
			<DataTableFacetedFilter
				column={ageCol}
				title="Age"
				options={ageOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(page.url.searchParams);
					if (selected) {
						sp.append('age', opt.value);
					} else {
						sp.delete('age', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if sizeCol}
			<DataTableFacetedFilter
				column={sizeCol}
				title="Size"
				options={sizeOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(page.url.searchParams);
					if (selected) {
						sp.append('size', opt.value);
					} else {
						sp.delete('size', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

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
							await deleteAnimals(ids).finally(() => {
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
