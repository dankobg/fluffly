<script lang="ts">
	import { type ComponentProps } from 'svelte';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import { AnimalAge, AnimalGender, AnimalSize, type components } from '$lib/gen/fluffly_openapi';
	import IconSearch from '@lucide/svelte/icons/search';
	import IconChevronDown from '@lucide/svelte/icons/chevron-down';
	import {
		animalAgeValues,
		animalGenderValues,
		daysLessThanOptions,
		daysGreaterThanOptions,
		animalSizeValues,
		radiusMetersOptions
	} from '$lib/enum-values';
	import { capitalize } from '$lib/utils';
	import Label from '$lib/components/ui/label/label.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import TagsInput from '$lib/components/tags-input/tags-input.svelte';
	import { useDebounce } from 'runed';
	import { getCurrentPosition } from '$lib/geocoding/geocoding';
	import { toast } from 'svelte-sonner';

	let {
		animalTypes,
		animalSpecies,
		animalBreeds,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & {
		animalTypes: components['schemas']['AnimalType'][];
		animalSpecies: components['schemas']['AnimalSpecie'][];
		animalBreeds: components['schemas']['Breed'][];
	} = $props();

	let position = $state<[number, number] | null>();
	let positionErr = $state<string | null>(null);
	let useMyLocation = $derived.by(() => {
		const latParam = page.url.searchParams.get('lat');
		const lonParam = page.url.searchParams.get('lon');
		if (!latParam || !lonParam) {
			return false;
		}
		const lat = Number.parseFloat(latParam);
		const lon = Number.parseFloat(lonParam);
		return !Number.isNaN(lat) && !Number.isNaN(lon);
	});

	let location = $derived(page.url.searchParams.get('location') ?? '');
	let radiusM = $derived.by(() => {
		const param = page.url.searchParams.get('radius_m');
		if (param) {
			const n = Number.parseFloat(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animalTypeId = $derived.by(() => {
		const param = page.url.searchParams.get('animal_type_id');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animalSpecieId = $derived.by(() => {
		const param = page.url.searchParams.get('animal_specie_id');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animalBreedIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('animal_breed_id');
		const ids = params.map(Number).filter(Number.isFinite);
		if (ids.length > 0) {
			return ids;
		}
	});
	let animalAges = $derived.by(() => {
		const params = page.url.searchParams.getAll('age');
		if (params.length > 0) {
			return params as AnimalAge[];
		}
	});
	let animalSizes = $derived.by(() => {
		const params = page.url.searchParams.getAll('size');
		if (params.length > 0) {
			return params as AnimalSize[];
		}
	});
	let animalGenders = $derived.by(() => {
		const params = page.url.searchParams.getAll('gender');
		if (params.length > 0) {
			return params as AnimalGender[];
		}
	});
	let animalHermaphrodite = $derived.by(() => {
		const param = page.url.searchParams.get('hermaphrodite');
		if (param) {
			return param === 'true';
		}
	});
	let animalMicrochip = $derived.by(() => {
		const param = page.url.searchParams.get('microchip');
		if (param) {
			return param === 'true';
		}
	});
	let animalName = $derived(page.url.searchParams.get('name') ?? '');
	let animalTags = $derived.by(() => {
		const params = page.url.searchParams.getAll('tag');
		if (params.length > 0) {
			return params;
		}
	});
	let animalDaysLt = $derived.by(() => {
		const param = page.url.searchParams.get('days_lt');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animalDatsGt = $derived.by(() => {
		const param = page.url.searchParams.get('days_gt');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});

	let properties = $derived.by(() => {
		return page.url.searchParams.entries().reduce(
			(acc, [key, value]) => {
				const match = key.match(/^properties\[(\w+)\]\[(\d+)\]$/);
				if (match) {
					const [, name, indexStr] = match;
					const index = Number(indexStr);
					if (name) {
						acc[name] ||= [];
						acc[name][index] = value;
					}
				}
				return acc;
			},
			{} as Record<string, string[]>
		);
	});

	const typeById = new Map(animalTypes.map(t => [t.id, t]));
	const speciesById = new Map<number, components['schemas']['AnimalSpecie']>();
	const breedById = new Map<number, components['schemas']['Breed']>();
	const speciesByTypeId = new Map<number, components['schemas']['AnimalSpecie'][]>();
	const breedsBySpecieId = new Map<number, components['schemas']['Breed'][]>();

	for (const s of animalSpecies) {
		speciesById.set(s.id, s);

		const list = speciesByTypeId.get(s.animal_type_id);
		if (list) {
			list.push(s);
		} else {
			speciesByTypeId.set(s.animal_type_id, [s]);
		}
	}
	for (const b of animalBreeds) {
		breedById.set(b.id, b);

		const list = breedsBySpecieId.get(b.animal_specie_id);
		if (list) {
			list.push(b);
		} else {
			breedsBySpecieId.set(b.animal_specie_id, [b]);
		}
	}

	let validSpecies = $derived.by(() => {
		if (animalTypeId) {
			return speciesByTypeId.get(animalTypeId);
		}
		if (animalBreedIds && animalBreedIds.length > 0) {
			const out: components['schemas']['AnimalSpecie'][] = [];
			const breeds = animalBreedIds.map(id => breedById.get(Number(id))).filter(Boolean);
			for (const breed of breeds) {
				if (breed?.animal_specie_id) {
					const specie = speciesById.get(breed.animal_specie_id);
					if (specie) {
						out.push(specie);
					}
				}
			}
			return out;
		}
		return animalSpecies;
	});

	let validBreeds = $derived.by(() => {
		if (!animalSpecieId && animalTypeId) {
			const species = speciesByTypeId.get(animalTypeId);
			const out: components['schemas']['Breed'][] = [];
			for (const s of species ?? []) {
				const breeds = breedsBySpecieId.get(s.id);
				if (breeds) {
					out.push(...breeds);
				}
			}
			return out;
		}
		if (animalSpecieId) {
			return breedsBySpecieId.get(animalSpecieId);
		}
		return animalBreeds;
	});

	let selectedRadiusM = $derived(
		radiusMetersOptions.find(x => x.value === String(radiusM))?.label ?? 'Select distance'
	);

	let selectedAnimalTypeName = $derived(animalTypes.find(x => x.id === animalTypeId)?.name ?? 'Select animal type');
	let selectedAnimalSpecie = $derived(animalSpecies.find(x => x.id === animalSpecieId));
	let selectedBreedNames = $derived.by(() => {
		const arr = animalBreeds.filter(x => Boolean(animalBreedIds?.includes(x.id))).map(x => x.name);
		if (arr.length > 2) {
			return `Selected ${arr.length} breeds`;
		}
		return arr.join(', ');
	});
	let selectedDaysLt = $derived(
		daysLessThanOptions.find(x => x.value === String(animalDaysLt))?.label ?? 'Select recent'
	);
	let selectedDaysGt = $derived(
		daysGreaterThanOptions.find(x => x.value === String(animalDatsGt))?.label ?? 'Select older'
	);

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

	const debounceNameFilter = useDebounce(
		() => {
			const sp = new URLSearchParams(page.url.searchParams);
			if (animalName) {
				sp.set('name', animalName);
			} else {
				sp.delete('name');
			}
			gotoWithFilters(sp);
		},
		() => 700
	);

	const debounceLocationFilter = useDebounce(
		() => {
			const sp = new URLSearchParams(page.url.searchParams);
			if (location) {
				sp.set('location', location);
			} else {
				sp.delete('location');
			}
			gotoWithFilters(sp);
		},
		() => 700
	);

	function gotoWithFilters(params: URLSearchParams) {
		goto('/search' + params.size ? `?${params}` : '', { keepFocus: true });
	}

	type SchemaProperty = {
		description: string;
		enum?: string[];
		title: string;
		type: string;
		'x-group': string;
	};

	type DynamicFilter = {
		key: string;
		val: SchemaProperty;
	};

	let groupedFilters = $derived.by(() => {
		if (!selectedAnimalSpecie) {
			return {};
		}

		const properties = (selectedAnimalSpecie?.properties_schema?.['properties'] ?? {}) as Record<
			string,
			SchemaProperty
		>;
		return Object.entries(properties).reduce(
			(acc, [key, val]) => {
				if (typeof val !== 'object') {
					return acc;
				}
				const group = val?.['x-group'] ?? 'Other';

				acc[group] ||= [];
				acc[group].push({ key, val });

				return acc;
			},
			{} as Record<string, DynamicFilter[]>
		);
	});

	function removeOldProperties(params: URLSearchParams): URLSearchParams {
		const remove: string[] = [];
		params.entries().forEach(([key]) => {
			const match = key.match(/^properties\[(\w+)\]\[(\d+)\]$/);
			if (match) {
				const [, name, indexStr] = match;
				const index = Number(indexStr);
				if (name) {
					remove.push(`properties[${name}][${index}]`);
				}
			}
		});
		remove.forEach(x => {
			params.delete(x);
		});
		return new URLSearchParams(params);
	}
</script>

<Sidebar.Root variant="floating" {...restProps} class="top-(--header-height) h-[calc(100svh-var(--header-height))]!">
	<Sidebar.Header>
		<h1 class="text-xl">Filter animals</h1>
	</Sidebar.Header>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>Geolocation</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<!-- Use current location filter -->
				<div>
					<div class="flex items-center gap-2">
						<Checkbox
							id="use_current_location_search_filter"
							checked={useMyLocation}
							onCheckedChange={async v => {
								useMyLocation = v;
								const sp = new URLSearchParams(page.url.searchParams);
								if (v) {
									try {
										position = await getCurrentPosition();
										positionErr = null;
										sp.set('lat', String(position[0]));
										sp.set('lon', String(position[1]));
										sp.delete('location');
									} catch (error: unknown) {
										if (error instanceof Error) {
											position = null;
											positionErr = [error.message, error.cause].filter(x => x).join(', ');
											toast.error(positionErr);
										}
									}
								} else {
									sp.delete('lat');
									sp.delete('lon');
								}
								gotoWithFilters(sp);
							}}
						/>
						<div class="space-y-1 leading-none">
							<Label for="use_current_location_search_filter">Use my location</Label>
						</div>
					</div>
				</div>

				<!-- Location filter -->
				{#if !useMyLocation}
					<div class="mt-4">
						<InputGroup.Root>
							<InputGroup.Input
								placeholder="Search near location"
								bind:value={
									() => location,
									v => {
										location = v;
										debounceLocationFilter();
									}
								}
							/>
							<InputGroup.Addon>
								<IconSearch />
							</InputGroup.Addon>
						</InputGroup.Root>
					</div>
				{/if}

				<!-- Radius meters filter -->
				<div class="mt-4">
					<Select.Root
						type="single"
						onValueChange={v => {
							const sp = new URLSearchParams(page.url.searchParams);
							sp.set('radius_m', v);
							gotoWithFilters(sp);
						}}
						value={radiusM !== undefined ? String(radiusM) : undefined}
					>
						<Select.Trigger class="w-full">
							{selectedRadiusM}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Distance</Select.Label>
								{#each radiusMetersOptions as distance (distance.value)}
									<Select.Item value={String(distance.value)} label={distance.label}>
										{distance.label}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
			</Sidebar.GroupContent>
		</Sidebar.Group>

		<Sidebar.Group>
			<Sidebar.GroupLabel>Common</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<!-- Animal type filter -->
				<div class="mt-4">
					<Select.Root
						type="single"
						onValueChange={v => {
							const sp = removeOldProperties(page.url.searchParams);
							sp.set('animal_type_id', v);

							// set animal_specie_id if only 1 exists for given animal_type_id
							const species = speciesByTypeId.get(Number(v)) ?? [];
							if (species?.length === 1) {
								sp.set('animal_specie_id', String(species[0]!.id));
							} else {
								sp.delete('animal_specie_id');
							}

							const allowedSpecieIds = new Set(species.map(s => s.id));
							const breedIds = (animalBreedIds ?? []).filter(breedId => {
								const breed = breedById.get(breedId);
								return breed && allowedSpecieIds.has(breed.animal_specie_id);
							});

							sp.delete('animal_breed_id');
							for (const breedId of breedIds) {
								sp.append('animal_breed_id', `${breedId}`);
							}

							gotoWithFilters(sp);
						}}
						value={animalTypeId ? String(animalTypeId) : undefined}
					>
						<Select.Trigger class="w-full">
							{selectedAnimalTypeName}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Animal type</Select.Label>
								{#each animalTypes as animalType (animalType.id)}
									<Select.Item value={String(animalType.id)} label={animalType.name}>
										{animalType.name}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>

				<!-- Animal species filter -->
				<div class="mt-4">
					<Select.Root
						type="single"
						onValueChange={v => {
							const sp = removeOldProperties(page.url.searchParams);
							sp.set('animal_specie_id', v);

							const species = speciesById.get(Number(v));
							if (species && animalTypeId !== species.animal_type_id) {
								sp.set('animal_type_id', String(species.animal_type_id));
							}
							const allowedBreedIds = new Set(breedsBySpecieId.get(Number(v))?.map(b => b.id) ?? []);
							const breedIds = (animalBreedIds ?? []).filter(id => allowedBreedIds.has(id));

							sp.delete('animal_breed_id');
							for (const breedId of breedIds) {
								sp.append('animal_breed_id', `${breedId}`);
							}

							gotoWithFilters(sp);
						}}
						value={animalSpecieId ? String(animalSpecieId) : undefined}
					>
						<Select.Trigger class="w-full">
							{selectedAnimalSpecie?.name ?? 'Select animal species'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Animal species</Select.Label>
								{#each validSpecies as animalSpecies (animalSpecies.id)}
									<Select.Item value={String(animalSpecies.id)} label={animalSpecies.name}>
										{animalSpecies.name}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>

				<!-- Animal breed filter -->
				<div class="mt-4">
					<Select.Root
						type="multiple"
						bind:value={
							() => (animalBreedIds ?? []).map(String),
							v => {
								const sp = new URLSearchParams(page.url.searchParams);
								sp.delete('animal_breed_id');
								for (const id of v) {
									sp.append('animal_breed_id', id);
								}

								if (v.length > 0) {
									const breed = breedById.get(Number(v![0]));
									if (!breed) {
										return;
									}
									const species = speciesById.get(breed.animal_specie_id);
									if (!species) {
										return;
									}
									const type = typeById.get(species.animal_type_id);
									if (!type) {
										return;
									}

									// set animal_specie_id that matches the selected breed
									if (animalSpecieId !== species.id) {
										sp.set('animal_specie_id', String(species.id));
									}
									// set animal_type_id that matches the selected breed
									if (animalTypeId !== type.id) {
										sp.set('animal_type_id', String(type.id));
									}
								}

								gotoWithFilters(sp);
							}
						}
					>
						<Select.Trigger class="w-full">
							{selectedBreedNames || 'Select animal breeds'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Animal breed</Select.Label>
								{#each validBreeds as animalBreed (animalBreed.id)}
									<Select.Item value={String(animalBreed.id)} label={animalBreed.name}>
										{animalBreed.name}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>

				<!-- Animal age filter -->
				<div class="mt-4">
					<Collapsible.Root class="group/collapsible w-full space-y-2">
						<div class="flex items-center justify-between space-x-4">
							<h4 class="text-sm font-semibold">Age</h4>
							<Collapsible.Trigger class={buttonVariants({ variant: 'ghost', size: 'sm', class: 'w-9 p-0' })}>
								<IconChevronDown class="ms-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
								<span class="sr-only">Toggle age filters</span>
							</Collapsible.Trigger>
						</div>
						<Collapsible.Content class="space-y-2 pl-4">
							{#each ageOptions as age (age.value)}
								<div class="flex items-center gap-2">
									<Checkbox
										id={`age-${age.value}`}
										bind:checked={
											() => Boolean(animalAges?.includes(age.value)),
											v => {
												const sp = new URLSearchParams(page.url.searchParams);
												if (v) {
													sp.append('age', age.value);
												} else {
													sp.delete('age', age.value);
												}
												gotoWithFilters(sp);
											}
										}
									/>
									<div class="space-y-1 leading-none">
										<Label for={`age-${age.value}`}>{age.label}</Label>
									</div>
								</div>
							{/each}
						</Collapsible.Content>
					</Collapsible.Root>
				</div>

				<!-- Animal size filter -->
				<div class="mt-4">
					<Collapsible.Root class="group/collapsible w-full space-y-2">
						<div class="flex items-center justify-between space-x-4">
							<h4 class="text-sm font-semibold">Size</h4>
							<Collapsible.Trigger class={buttonVariants({ variant: 'ghost', size: 'sm', class: 'w-9 p-0' })}>
								<IconChevronDown class="ms-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
								<span class="sr-only">Toggle size filters</span>
							</Collapsible.Trigger>
						</div>
						<Collapsible.Content class="space-y-2 pl-4">
							{#each sizeOptions as size (size.value)}
								<div class="flex items-center gap-2">
									<Checkbox
										id={`size-${size.value}`}
										bind:checked={
											() => Boolean(animalSizes?.includes(size.value)),
											v => {
												const sp = new URLSearchParams(page.url.searchParams);
												if (v) {
													sp.append('size', size.value);
												} else {
													sp.delete('size', size.value);
												}
												gotoWithFilters(sp);
											}
										}
									/>
									<div class="space-y-1 leading-none">
										<Label for={`size-${size.value}`}>{size.label}</Label>
									</div>
								</div>
							{/each}
						</Collapsible.Content>
					</Collapsible.Root>
				</div>

				<div class="mt-4">
					<Collapsible.Root class="group/collapsible w-full space-y-2">
						<div class="flex items-center justify-between space-x-4">
							<h4 class="text-sm font-semibold">Gender</h4>
							<Collapsible.Trigger class={buttonVariants({ variant: 'ghost', size: 'sm', class: 'w-9 p-0' })}>
								<IconChevronDown class="ms-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
								<span class="sr-only">Toggle gender filters</span>
							</Collapsible.Trigger>
						</div>
						<Collapsible.Content class="space-y-2 pl-4">
							{#each genderOptions as gender (gender.value)}
								<div class="flex items-center gap-2">
									<Checkbox
										id={`gender-${gender.value}`}
										bind:checked={
											() => Boolean(animalGenders?.includes(gender.value)),
											v => {
												const sp = new URLSearchParams(page.url.searchParams);
												if (v) {
													sp.append('gender', gender.value);
												} else {
													sp.delete('gender', gender.value);
												}
												gotoWithFilters(sp);
											}
										}
									/>
									<div class="space-y-1 leading-none">
										<Label for={`gender-${gender.value}`}>{gender.label}</Label>
									</div>
								</div>
							{/each}
						</Collapsible.Content>
					</Collapsible.Root>
				</div>

				<!-- Animal hermaphrodite filter -->
				<div class="mt-4">
					<div class="flex items-center gap-2">
						<Checkbox
							id="gender"
							bind:checked={
								() => Boolean(animalHermaphrodite),
								v => {
									const sp = new URLSearchParams(page.url.searchParams);
									sp.set('hermaphrodite', String(v));
									gotoWithFilters(sp);
								}
							}
						/>
						<div class="space-y-1 leading-none">
							<Label for="gender">Hermaphrodite</Label>
						</div>
					</div>
				</div>

				<!-- Animal microchip filter -->
				<div class="mt-4">
					<div class="flex items-center gap-2">
						<Checkbox
							id="microchip"
							bind:checked={
								() => Boolean(animalMicrochip),
								v => {
									const sp = new URLSearchParams(page.url.searchParams);
									sp.set('microchip', String(v));
									gotoWithFilters(sp);
								}
							}
						/>
						<div class="space-y-1 leading-none">
							<Label for="microchip">Microchip</Label>
						</div>
					</div>
				</div>

				<div class="mt-4">
					<InputGroup.Root>
						<InputGroup.Input
							placeholder="Search pet by name"
							bind:value={
								() => animalName,
								v => {
									animalName = v;
									debounceNameFilter();
								}
							}
						/>
						<InputGroup.Addon>
							<IconSearch />
						</InputGroup.Addon>
					</InputGroup.Root>
				</div>

				<!-- Animal tag filter -->
				<div class="mt-4">
					<TagsInput
						placeholder="Search pet by tags"
						bind:value={
							() => animalTags ?? [],
							v => {
								const sp = new URLSearchParams(page.url.searchParams);
								sp.delete('tag');
								v.forEach(x => {
									sp.append('tag', x);
								});
								gotoWithFilters(sp);
							}
						}
					/>
				</div>
			</Sidebar.GroupContent>
		</Sidebar.Group>

		<!-- Animal properties filter -->
		{#each Object.entries(groupedFilters) as [group, dynamicFilters] (group)}
			<Sidebar.Group>
				<Sidebar.GroupLabel>{group}</Sidebar.GroupLabel>
				<Sidebar.GroupContent>
					{#each dynamicFilters as dynamicFilter (dynamicFilter.key)}
						{#if dynamicFilter.val.type === 'boolean'}
							<div class="mt-4">
								<div class="flex items-center gap-2">
									<Checkbox
										id={`${dynamicFilter.val.title}-${dynamicFilter.key}`}
										bind:checked={
											() => Boolean(properties?.[dynamicFilter.key]?.includes('true')),
											v => {
												const sp = new URLSearchParams(page.url.searchParams);
												sp.set(`properties[${dynamicFilter.key}][0]`, String(v));
												gotoWithFilters(sp);
											}
										}
									/>
									<div class="space-y-1 leading-none">
										<Label for={`${dynamicFilter.val.title}-${dynamicFilter.key}`}>{dynamicFilter.val.title}</Label>
									</div>
								</div>
							</div>
						{/if}

						{#if dynamicFilter.val.type === 'string'}
							{#if dynamicFilter.val.enum}
								<div class="mt-4">
									<Collapsible.Root class="group/collapsible w-full space-y-2">
										<div class="flex items-center justify-between space-x-4">
											<h4 class="text-sm font-semibold">{dynamicFilter.val.title}</h4>
											<Collapsible.Trigger class={buttonVariants({ variant: 'ghost', size: 'sm', class: 'w-9 p-0' })}>
												<IconChevronDown
													class="ms-auto transition-transform group-data-[state=open]/collapsible:rotate-180"
												/>
												<span class="sr-only">Toggle {dynamicFilter.key} filters</span>
											</Collapsible.Trigger>
										</div>
										<Collapsible.Content class="space-y-2 pl-4">
											{#each dynamicFilter.val.enum ?? [] as item (item)}
												<div class="flex items-center gap-2">
													<Checkbox
														id={`${dynamicFilter.key}-${item}`}
														bind:checked={
															() => Boolean(properties?.[dynamicFilter.key]?.includes(item)),
															v => {
																const sp = new URLSearchParams(page.url.searchParams);
																if (v) {
																	const n: number = (properties?.[dynamicFilter.key] ?? []).length;
																	sp.append(`properties[${dynamicFilter.key}][${n}]`, item);
																} else {
																	const newItems = properties[dynamicFilter.key]?.filter(x => x !== item);
																	properties[dynamicFilter.key]?.forEach((_, i) => {
																		sp.delete(`properties[${dynamicFilter.key}][${i}]`);
																	});
																	newItems?.forEach((x, i) => {
																		sp.append(`properties[${dynamicFilter.key}][${i}]`, x);
																	});
																}
																gotoWithFilters(sp);
															}
														}
													/>
													<div class="space-y-1 leading-none">
														<Label for={`${dynamicFilter.key}-${item}`}>{capitalize(item)}</Label>
													</div>
												</div>
											{/each}
										</Collapsible.Content>
									</Collapsible.Root>
								</div>
							{:else}
								<div class="mt-4">
									<div class="flex items-center gap-2">
										<TagsInput
											id={`${dynamicFilter.val.title}-${dynamicFilter.key}`}
											placeholder={`${dynamicFilter.val.title}`}
											bind:value={
												() => properties?.[dynamicFilter.key] ?? [],
												v => {
													const sp = new URLSearchParams(page.url.searchParams);
													for (const key of [...sp.keys()]) {
														if (key.startsWith(`properties[${dynamicFilter.key}]`)) {
															sp.delete(key);
														}
													}
													(v as string[]).forEach((x, i) => {
														sp.append(`properties[${dynamicFilter.key}][${i}]`, x);
													});
													gotoWithFilters(sp);
												}
											}
										/>
									</div>
								</div>
							{/if}
						{/if}
					{/each}
				</Sidebar.GroupContent>
			</Sidebar.Group>
		{/each}

		<Sidebar.Group>
			<Sidebar.GroupLabel>Time on fluffly</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<!-- Animal recent (days less than) on fluffly -->
				<div class="mt-4">
					<Select.Root
						type="single"
						onValueChange={v => {
							const sp = new URLSearchParams(page.url.searchParams);
							sp.set('days_lt', v);
							gotoWithFilters(sp);
						}}
						value={animalDaysLt ? String(animalDaysLt) : undefined}
					>
						<Select.Trigger class="w-full">
							{selectedDaysLt}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Animal type</Select.Label>
								{#each daysLessThanOptions as day (day.value)}
									<Select.Item value={String(day.value)} label={day.label}>
										{day.label}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>

				<!-- Animal overlooked (days greater than) on fluffly -->
				<div class="mt-4">
					<Select.Root
						type="single"
						onValueChange={v => {
							const sp = new URLSearchParams(page.url.searchParams);
							sp.set('days_gt', v);
							gotoWithFilters(sp);
						}}
						value={animalDatsGt ? String(animalDatsGt) : undefined}
					>
						<Select.Trigger class="w-full">
							{selectedDaysGt}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Overlooked animals</Select.Label>
								{#each daysGreaterThanOptions as day (day.value)}
									<Select.Item value={String(day.value)} label={day.label}>
										{day.label}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
</Sidebar.Root>
