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
	import { page } from '$app/state';
	import { useDebounce } from 'runed';
	import { goto, invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { toast } from 'svelte-sonner';

	let { table }: { table: Table<TData> } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	let nameCol = $derived(table.getColumn('name'));

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
		if (nameCol && name) {
			nameCol.setFilterValue(name);
		}
	});

	async function deleteAnimalTypes(ids: number[]) {
		try {
			const deleteAnimalTypesResult = await fluffly.DELETE('/animal_types', {
				body: { ids }
			});
			if (deleteAnimalTypesResult.error) {
				toast.error(
					[deleteAnimalTypesResult.error.message, deleteAnimalTypesResult.error.reason].filter(Boolean).join(', ')
				);
				return;
			}
			toast.success('Animal types deleted');
			invalidate('data:dashboard-animal-types');
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
							await deleteAnimalTypes(ids).finally(() => {
								table.resetRowSelection();
							});
						}
					});
				}}
			>
				Delete {table.getFilteredSelectedRowModel().rows.length} animal types
				<IconTrash />
			</Button>
		{/if}
	</div>

	<DataTableViewOptions {table} />
</div>
