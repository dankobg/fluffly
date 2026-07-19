<script lang="ts">
	import type { PageProps } from './$types';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import AnimalCard from '$lib/components/animal-card/animal-card.svelte';
	import { page } from '$app/state';
	import { goto, invalidate } from '$app/navigation';

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
		goto('/dashboard/animals/favorite' + params.size ? `?${params}` : '', { keepFocus: true });
	}
</script>

<div class="grid w-full grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] gap-4 p-4">
	{#each data?.animalsResult?.data?.data ?? [] as animal (animal.id)}
		<AnimalCard
			{animal}
			showPostedBy={true}
			showLikeUnlikeButton={data?.auth?.session?.active}
			className="dark:bg-background/70"
			onLiked={() => {
				invalidate('data:dashboard-favorite-animals');
			}}
			onUnliked={() => {
				invalidate('data:dashboard-favorite-animals');
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
