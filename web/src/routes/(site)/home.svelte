<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import type { components } from '$lib/gen/fluffly_openapi';
	import type { PageProps } from './$types';
	import IconSearch from '@lucide/svelte/icons/search';
	import { SvelteMap } from 'svelte/reactivity';
	import { radiusMetersOptions } from '$lib/enum-values';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import { toast } from 'svelte-sonner';
	import { getCurrentPosition } from '$lib/geocoding/geocoding';
	import { cn } from '$lib/utils';

	let { data }: PageProps = $props();

	const speciesById = new Map<number, components['schemas']['AnimalSpecie']>();
	const speciesByTypeId = new Map<number, components['schemas']['AnimalSpecie'][]>();

	for (const s of data.animalSpeciesResult?.data?.data ?? []) {
		speciesById.set(s.id, s);

		const list = speciesByTypeId.get(s.animal_type_id);
		if (list) {
			list.push(s);
		} else {
			speciesByTypeId.set(s.animal_type_id, [s]);
		}
	}

	const animalTypesOrder = [
		'Cat',
		'Dog',
		'Rabbit',
		'Bird',
		'Small & Furry',
		'Scales, Fins & Other',
		'Horse',
		'Barnyard'
	];

	const animalTypeIcons = new Map([
		['Cat', '/images/animals/cat.svg'],
		['Dog', '/images/animals/dog.svg'],
		['Rabbit', '/images/animals/rabbit.svg'],
		['Bird', '/images/animals/bird-1.svg'],
		['Small & Furry', '/images/animals/mouse.svg'],
		['Scales, Fins & Other', '/images/animals/fish-2.svg'],
		['Horse', '/images/animals/horse-2.svg'],
		['Barnyard', '/images/animals/cow.svg']
	]);

	let animalTypeByName = $derived.by(() => {
		return new SvelteMap(data?.animalTypesResult?.data?.data?.map(t => [t.name, t]) ?? []);
	});

	let animalTypes = $derived.by(() => {
		return animalTypesOrder
			.map(name => {
				const type = animalTypeByName.get(name);
				if (!type) {
					return;
				}
				return {
					id: type.id,
					name: type.name,
					src: animalTypeIcons.get(name)
				};
			})
			.filter(Boolean);
	});

	let useMyLocation = $state<boolean>(false);
	let location = $state<string | undefined>();
	let radiusM = $state<number | undefined>();
	let position = $state<[number, number] | null>();
	let positionErr = $state<string | null>();

	let selectedRadiusM = $derived(radiusMetersOptions.find(x => x.value === String(radiusM))?.label ?? 'Radius');

	function resolveSearchUrl(animalTypeId: number): string {
		const sp = new URLSearchParams();
		sp.set('animal_type_id', String(animalTypeId));

		const species = speciesByTypeId.get(animalTypeId) ?? [];
		if (species?.length === 1) {
			sp.set('animal_specie_id', String(species[0]!.id));
		} else {
			sp.delete('animal_specie_id');
		}

		if (location) {
			sp.set('location', String(location));
		}
		if (position) {
			sp.set('lat', String(position[0]));
			sp.set('lon', String(position[1]));
		}
		if (radiusM !== undefined) {
			sp.set('radius_m', String(radiusM));
		}

		return `/search?${sp}`;
	}
</script>

<section
	class="flex max-h-[25rem] flex-col items-center justify-center gap-8 bg-[linear-gradient(rgb(144,66,35,0.3),rgb(144,66,35,0.3)),url('/images/pexels/pexels-12.webp')] bg-cover bg-center bg-no-repeat p-4 md:p-8 lg:p-12"
>
	<div>
		<h1 class="text-center text-4xl font-bold text-white">Find and adopt your new best pal with Fluffly</h1>
		<h2 class="mt-4 text-center text-xl font-semibold text-white">Browse your favorite pets</h2>
	</div>

	<div class="grid max-w-lg gap-4 rounded-xl bg-black/30 p-4">
		<div class="grid gap-2">
			<Label for="search-location" class="text-white">Search location</Label>
			<InputGroup.Root class="bg-input dark:bg-input">
				<InputGroup.Input
					id="search-location"
					placeholder="Enter street address..."
					bind:value={location}
					disabled={useMyLocation}
				/>
				<InputGroup.Addon>
					<IconSearch />
				</InputGroup.Addon>
			</InputGroup.Root>
		</div>

		<div class="flex flex-col justify-between gap-4 md:flex-row">
			<div class="flex items-center gap-2">
				<Checkbox
					class="bg-input dark:bg-input"
					id="use-my-location"
					checked={useMyLocation}
					onCheckedChange={async v => {
						useMyLocation = v;
						if (v) {
							try {
								position = await getCurrentPosition();
								positionErr = null;
							} catch (error: unknown) {
								if (error instanceof Error) {
									position = null;
									positionErr = [error.message, error.cause].filter(x => x).join(', ');
									toast.error(positionErr);
								}
							}
						}
					}}
				/>
				<Label for="use-my-location" class="text-white">Use my location</Label>
			</div>

			<div class="flex items-center gap-2">
				<Label for="search-distance" class="text-white">Within distance</Label>
				<Select.Root
					type="single"
					bind:value={
						() => (radiusM !== undefined ? String(radiusM) : undefined),
						v => {
							radiusM = v ? Number(v) : undefined;
						}
					}
				>
					<Select.Trigger id="search-distance" class="bg-input dark:bg-input hover:dark:bg-input">
						{selectedRadiusM}
					</Select.Trigger>
					<Select.Content>
						{#each radiusMetersOptions as radius (radius.value)}
							<Select.Item {...radius} />
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>
</section>

<div class="mx-auto mt-6 grid max-w-[50rem] gap-4">
	<section
		class="grid grid-cols-[repeat(auto-fill,minmax(min(10rem,100%),1fr))] justify-center gap-4 p-4
		lg:p-0 [&:has(.animal-card:hover)_.animal-card:not(:hover)]:scale-98 [&:has(.animal-card:hover)_.animal-card:not(:hover)]:grayscale-20"
	>
		{#each animalTypes ?? [] as item (item?.name)}
			{#if item}
				<a href={resolveSearchUrl(item.id)} class="group animal-card transition-all duration-200">
					<div
						class="flex flex-col items-center justify-center gap-2 rounded-xl border-4 border-pink-600 p-3
							transition-all duration-200 ease-out
							hover:-translate-y-1 hover:scale-[1.05] hover:border-pink-400 hover:bg-pink-300/20 hover:shadow-lg dark:hover:bg-pink-200/10"
					>
						<div class="overflow-hidden rounded-md">
							<img
								class={cn(
									'max-w-full object-cover transition-transform duration-200',
									Math.random() > 0.5 ? 'group-hover:rotate-7' : 'group-hover:-rotate-7'
								)}
								src={item.src}
								alt={item.name}
							/>
						</div>

						<p class="text-sm font-medium transition-colors">
							{item.name}
						</p>
					</div>
				</a>
			{/if}
		{/each}
	</section>
</div>

<h3 class="mt-6 text-center text-xl">Or search by organizations</h3>
<div class="mt-4 mb-10 flex justify-center">
	<a href="/search?organizations=true" class="group">
		<div
			class="flex w-full max-w-[12rem] flex-col items-center justify-center gap-2 rounded-md border-4 border-pink-600 p-3 transition-all duration-200 ease-out hover:-translate-y-1 hover:scale-[1.05] hover:border-pink-400 hover:bg-pink-300/20 hover:shadow-lg dark:hover:bg-pink-200/10"
		>
			<img
				class={cn(
					'max-w-full object-cover transition-transform duration-200',
					Math.random() > 0.5 ? 'group-hover:rotate-7' : 'group-hover:-rotate-7'
				)}
				src="/images/animals/shelter.svg"
				alt="Rescues & shelters"
			/>
			<p class="text-sm">Rescues & shelters</p>
		</div>
	</a>
</div>
