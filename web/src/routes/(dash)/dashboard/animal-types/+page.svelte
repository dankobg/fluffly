<script lang="ts">
	import DataTable from './animal-types-data-grid/data-table.svelte';
	import { columns } from './animal-types-data-grid/columns';
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconPlus from '@lucide/svelte/icons/plus';

	let { data }: PageProps = $props();
</script>

{#if data.animalTypesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Animal types</h2>
				<p class="text-muted-foreground">List of animal types</p>
			</div>
		</div>
		<div>
			<a href="/dashboard/animal-types/create">
				<Button>
					Create new animal type
					<IconPlus />
				</Button>
			</a>
		</div>
		<DataTable data={data.animalTypesResult?.data.data ?? []} meta={data.animalTypesResult?.data.meta} {columns} />
	</div>
{/if}
