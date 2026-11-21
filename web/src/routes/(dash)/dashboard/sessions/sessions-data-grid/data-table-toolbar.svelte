<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableFacetedFilter from '#lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import { aals } from './data';
	import Checkbox from '#lib/components/ui/checkbox/checkbox.svelte';
	import Label from '#lib/components/ui/label/label.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import { expandSessions } from '../expandSessions.svelte';
	import { PathsSessionsGetParametersQueryExpand } from '#lib/gen/fluffly_openapi.ts';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import FilterInput from '#lib/components/data-grid-shared/filter-input.svelte';

	const table = useTableContext();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	const emailCol = $derived(table.getColumn('email'));
	const fullnameCol = $derived(table.getColumn('full_name'));
	const aalCol = $derived(table.getColumn('authenticator_assurance_level'));
	const expandOptions = Object.values(PathsSessionsGetParametersQueryExpand);

	const searchParams = new SvelteURLSearchParams(new URL(page.url.href).searchParams);

	function onExpandCheckedChange(val: boolean, opt: PathsSessionsGetParametersQueryExpand) {
		const all = searchParams.getAll('expand') as PathsSessionsGetParametersQueryExpand[];
		let newUrl = page.url.pathname;
		if (val) {
			if (!all.includes(opt)) {
				searchParams.append('expand', opt);
				expandSessions.expanded = [...all, opt];
			}
		} else {
			searchParams.delete('expand');
			for (const p of all.filter(x => x !== opt)) {
				searchParams.append('expand', p);
			}
			expandSessions.expanded = expandSessions.expanded.filter(x => x !== opt);
		}
		if (searchParams.size > 0) {
			newUrl += `?${searchParams}`;
		}
		goto(newUrl);
	}

	$effect(() => {
		expandSessions.expanded = searchParams.getAll('expand') as PathsSessionsGetParametersQueryExpand[];
	});
</script>

<div class="flex gap-4">
	<p>Expand</p>
	<div class="flex gap-2">
		{#each expandOptions as opt (opt)}
			{@const checked = expandSessions.expanded.includes(opt)}
			<div class="flex items-center gap-1">
				<Checkbox
					id="expand-{opt}"
					aria-labelledby="expand-{opt}-label"
					{checked}
					value={opt}
					onCheckedChange={val => onExpandCheckedChange(val, opt)}
				/>
				<Label
					id="expand-{opt}-label"
					for="expand-{opt}"
					class="text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
				>
					{opt}
				</Label>
			</div>
		{/each}
	</div>
</div>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if expandSessions.expanded.includes(PathsSessionsGetParametersQueryExpand.identity)}
			{#if emailCol}
				<FilterInput column={emailCol} />
			{/if}
			{#if fullnameCol}
				<FilterInput column={fullnameCol} title="full name" />
			{/if}
		{/if}

		{#if aalCol}
			<DataTableFacetedFilter column={aalCol} title="Aal" options={aals} />
		{/if}

		{#if isFiltered}
			<Button variant="ghost" onclick={() => table.resetColumnFilters()} class="h-8 px-2 lg:px-3">
				Reset
				<IconX />
			</Button>
		{/if}
	</div>
	<DataTableViewOptions />
</div>
