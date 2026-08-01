<script lang="ts">
	import type { PageProps } from './$types';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import IconTrash2 from '@lucide/svelte/icons/trash-2';
	import IconSearch from '@lucide/svelte/icons/search';
	import RemovableTag from '$lib/components/removable-tag/removable-tag.svelte';
	import OrganizationCard from '$lib/components/organization-card/organization-card.svelte';
	import { useDebounce } from 'runed';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';

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
		goto('/organizations' + params.size ? `?${params}` : '', { keepFocus: true });
	}

	let name = $derived(page.url.searchParams.get('name') ?? undefined);

	const debounceNameFilter = useDebounce(
		() => {
			const sp = new URLSearchParams(page.url.searchParams);
			if (name) {
				sp.set('name', name);
			} else {
				sp.delete('name');
			}
			gotoWithFilters(sp);
		},
		() => 700
	);

	let filters = $derived.by(() => {
		return {
			name
		};
	});

	let hasFilters = $derived(Object.values(filters).filter(x => !!x).length > 0);
</script>

<div class="mt-4 flex justify-center">
	<div class="w-full max-w-lg">
		<InputGroup.Root>
			<InputGroup.Input
				placeholder="Search organization by name"
				bind:value={
					() => name,
					v => {
						name = v;
						debounceNameFilter();
					}
				}
			/>
			<InputGroup.Addon>
				<IconSearch />
			</InputGroup.Addon>
		</InputGroup.Root>
	</div>
</div>

<div class="w-full">
	{#if hasFilters}
		<div class="p-4">
			<div class="flex flex-wrap items-center gap-2">
				<Button
					variant="destructive"
					onclick={() => {
						gotoWithFilters(new URLSearchParams());
					}}
				>
					Clear all filters
					<IconTrash2 />
				</Button>

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
				{/each}
			</div>
		</div>
	{/if}

	<div class="grid w-full grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] gap-4 p-4">
		{#each data?.organizationsResult?.data?.data ?? [] as organization (organization.id)}
			<OrganizationCard {organization} />
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
			count={data?.organizationsResult?.data?.meta?.total ?? 0}
			page={data?.organizationsResult?.data?.meta?.page ?? pageNum}
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
