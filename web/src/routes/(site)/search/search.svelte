<script lang="ts">
	import type { PageProps } from './$types';
	import SearchSidebar from './search-sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';
	import AnimalCard from '$lib/components/animal-card/animal-card.svelte';
	import type { AnimalAge, AnimalGender, AnimalSize } from '$lib/gen/fluffly_openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import IconTrash2 from '@lucide/svelte/icons/trash-2';
	import RemovableTag from '$lib/components/removable-tag/removable-tag.svelte';

	let { data }: PageProps = $props();

	let pageNum = $derived.by(() => {
		const param = page.url.searchParams.get('page');
		if (!param) {
			return 1;
		}
		const n = Number.parseInt(param);
		return !Number.isNaN(n) ? n : 1;
	});
	let pageSize = $derived.by(() => {
		const param = page.url.searchParams.get('page_size');
		if (!param) {
			return 50;
		}
		const n = Number.parseInt(param);
		return !Number.isNaN(n) ? n : 50;
	});

	const pageSizes = [10, 20, 30, 50, 100];
	const selectedPageSizeText = $derived(pageSizes.find(x => x === pageSize) ?? 'Page size');

	function gotoWithFilters(params: URLSearchParams) {
		goto('/search' + params.size ? `?${params}` : '', { keepFocus: true });
	}

	let use_my_location = $derived.by(() => {
		const latParam = page.url.searchParams.get('lat');
		const lonParam = page.url.searchParams.get('lon');
		if (!latParam || !lonParam) {
			return;
		}
		const lat = Number.parseFloat(latParam);
		const lon = Number.parseFloat(lonParam);
		return !Number.isNaN(lat) && !Number.isNaN(lon);
	});

	let location = $derived(page.url.searchParams.get('location') ?? undefined);

	let radius_m = $derived.by(() => {
		const param = page.url.searchParams.get('radius_m');
		if (param) {
			const n = Number.parseFloat(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animal_type_id = $derived.by(() => {
		const param = page.url.searchParams.get('animal_type_id');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animal_specie_id = $derived.by(() => {
		const param = page.url.searchParams.get('animal_specie_id');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let animal_breed_id = $derived.by(() => {
		const params = page.url.searchParams.getAll('animal_breed_id');
		const ids = params.map(Number).filter(Number.isFinite);
		if (ids.length > 0) {
			return ids;
		}
	});
	let age = $derived.by(() => {
		const params = page.url.searchParams.getAll('age');
		if (params.length > 0) {
			return params as AnimalAge[];
		}
	});
	let size = $derived.by(() => {
		const params = page.url.searchParams.getAll('size');
		if (params.length > 0) {
			return params as AnimalSize[];
		}
	});
	let gender = $derived.by(() => {
		const params = page.url.searchParams.getAll('gender');
		if (params.length > 0) {
			return params as AnimalGender[];
		}
	});
	let hermaphrodite = $derived.by(() => {
		const param = page.url.searchParams.get('hermaphrodite');
		if (param) {
			return param === 'true';
		}
	});
	let microchip = $derived.by(() => {
		const param = page.url.searchParams.get('microchip');
		if (param) {
			return param === 'true';
		}
	});
	let name = $derived(page.url.searchParams.get('name') ?? undefined);

	let tag = $derived.by(() => {
		const params = page.url.searchParams.getAll('tag');
		if (params.length > 0) {
			return params;
		}
	});
	let days_lt = $derived.by(() => {
		const param = page.url.searchParams.get('days_lt');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});
	let days_gt = $derived.by(() => {
		const param = page.url.searchParams.get('days_gt');
		if (param) {
			const n = Number.parseInt(param);
			if (!Number.isNaN(n)) {
				return n;
			}
		}
	});

	let properties = $derived.by(() => {
		const props = page.url.searchParams.entries().reduce(
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

		if (Object.keys(props).length > 0) {
			return props;
		}
	});

	let filters = $derived.by(() => {
		return {
			use_my_location,
			location,
			radius_m,
			animal_type_id,
			animal_specie_id,
			animal_breed_id,
			age,
			size,
			gender,
			hermaphrodite,
			microchip,
			name,
			tag,
			days_lt,
			days_gt,
			properties
		};
	});

	let hasFilters = $derived(Object.values(filters).filter(x => !!x).length > 0);
</script>

<Sidebar.Provider style="--sidebar-width: 19rem;">
	<SearchSidebar
		animalTypes={data?.animalTypesResult?.data?.data ?? []}
		animalSpecies={data?.animalSpeciesResult?.data?.data ?? []}
		animalBreeds={data?.animalBreedsResult?.data?.data ?? []}
	/>
	<Sidebar.Inset class="mx-auto max-w-[120rem]">
		{@render searchContent()}
	</Sidebar.Inset>
</Sidebar.Provider>

{#snippet searchContent()}
	<div class="w-full">
		<label
			class="mt-4 ml-4 flex w-fit items-center rounded-md border bg-background px-2 shadow-xs hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50"
		>
			<Sidebar.Trigger />
			Toggle filters sidebar
		</label>

		<div class="p-4">
			<div class="flex flex-wrap items-center gap-2">
				{#if hasFilters}
					<Button
						variant="destructive"
						onclick={() => {
							gotoWithFilters(new URLSearchParams());
						}}
					>
						Clear all filters
						<IconTrash2 />
					</Button>
				{/if}

				{#each Object.entries(filters) as [key, val] (key)}
					{#if ['string', 'number', 'boolean'].includes(typeof val)}
						<RemovableTag
							onClose={() => {
								const sp = new URLSearchParams(page.url.searchParams);
								sp.delete(key, `${val}`);
								gotoWithFilters(sp);
							}}
							class="rounded-full bg-fuchsia-800"
						>
							{key}
						</RemovableTag>
					{/if}

					{#if typeof val === 'object'}
						{#if Array.isArray(val)}
							{#each val as item (item)}
								<RemovableTag
									onClose={() => {
										const sp = new URLSearchParams(page.url.searchParams);
										sp.delete(key, `${item}`);
										gotoWithFilters(sp);
									}}
									class="rounded-full bg-fuchsia-800"
								>
									{item}
								</RemovableTag>
							{/each}
						{:else}
							{#each Object.entries(val) as [propName, propValues] (propName)}
								{#each propValues as prop (prop)}
									<RemovableTag
										onClose={() => {
											const sp = new URLSearchParams(page.url.searchParams);
											const newItems = propValues.filter(x => x !== prop);
											propValues.forEach((pv, idx) => {
												sp.delete(`${key}[${propName}][${idx}]`, `${pv}`);
											});
											newItems.forEach((pv, idx) => {
												sp.append(`${key}[${propName}][${idx}]`, `${pv}`);
											});
											gotoWithFilters(sp);
										}}
										class="rounded-full bg-fuchsia-800"
									>
										{#if prop === 'true'}
											<IconCheck />
											{propName}
										{:else if prop === 'false'}
											<IconCircleX />
											{propName}
										{:else}
											{prop}
										{/if}
									</RemovableTag>
								{/each}
							{/each}
						{/if}
					{/if}
				{/each}
			</div>
		</div>

		<div class="grid w-full grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] gap-4 p-4">
			{#each data?.animalsResult?.data?.data ?? [] as animal (animal.id)}
				<AnimalCard
					{animal}
					showPostedBy={true}
					showLikeUnlikeButton={data?.auth?.session?.active}
					onLiked={() => {
						invalidate('data:search');
					}}
					onUnliked={() => {
						invalidate('data:search');
					}}
				/>
			{/each}
		</div>

		<div class="flex items-center gap-4 p-4">
			<div class="flex items-center gap-2">
				<span>Page size</span>
				<Select.Root
					type="single"
					name="favoriteFruit"
					bind:value={
						() => String(pageSize),
						v => {
							const sp = new URLSearchParams(page.url.searchParams);
							sp.set('page_size', v);
							gotoWithFilters(sp);
						}
					}
				>
					<Select.Trigger>
						{selectedPageSizeText}
					</Select.Trigger>
					<Select.Content>
						<Select.Group>
							<Select.Label>Page size</Select.Label>
							{#each pageSizes as num (num)}
								<Select.Item value={String(num)} label={String(num)}>
									{num}
								</Select.Item>
							{/each}
						</Select.Group>
					</Select.Content>
				</Select.Root>
			</div>

			<Pagination.Root
				class="ml-auto flex-1 justify-end"
				count={data?.animalsResult?.data?.meta?.total ?? 0}
				page={data?.animalsResult?.data?.meta?.page ?? pageNum}
				perPage={pageSize}
			>
				{#snippet children({ pages, currentPage })}
					<Pagination.Content>
						<Pagination.Item>
							<Pagination.Previous
								onclick={() => {
									const sp = new URLSearchParams(page.url.searchParams);
									sp.set('page', String(currentPage - 1));
									gotoWithFilters(sp);
								}}
							/>
						</Pagination.Item>
						{#each pages as pageItem (pageItem.key)}
							{#if pageItem.type === 'ellipsis'}
								<Pagination.Item>
									<Pagination.Ellipsis />
								</Pagination.Item>
							{:else}
								<Pagination.Item>
									<Pagination.Link
										page={pageItem}
										isActive={currentPage === pageItem.value}
										onclick={() => {
											const sp = new URLSearchParams(page.url.searchParams);
											sp.set('page', String(pageItem.value));
											gotoWithFilters(sp);
										}}
									>
										{pageItem.value}
									</Pagination.Link>
								</Pagination.Item>
							{/if}
						{/each}
						<Pagination.Item>
							<Pagination.Ellipsis />
						</Pagination.Item>
						<Pagination.Item>
							<Pagination.Next
								onclick={() => {
									const sp = new URLSearchParams(page.url.searchParams);
									sp.set('page', String(currentPage + 1));
									gotoWithFilters(sp);
								}}
							/>
						</Pagination.Item>
					</Pagination.Content>
				{/snippet}
			</Pagination.Root>
		</div>
	</div>
{/snippet}
