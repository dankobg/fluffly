<script lang="ts" generics="TData extends RowData">
	import Input from '#lib/components/ui/input/input.svelte';
	import type { RowData, Column } from '@tanstack/svelte-table';
	import type { AppFeatures } from './app-table';

	let {
		title,
		column
	}: {
		title?: string;
		column: Column<AppFeatures, TData>;
	} = $props();

	const placeholder = $derived(column.columnDef.meta?.label ?? title ?? column.id);
</script>

<Input
	placeholder="filter by {placeholder}..."
	value={(column?.getFilterValue() as string) ?? ''}
	oninput={e => {
		column?.setFilterValue(e.currentTarget.value);
	}}
	onchange={e => {
		column?.setFilterValue(e.currentTarget.value);
	}}
	class="h-8 w-[150px] lg:w-[250px]"
/>
