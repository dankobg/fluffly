<script lang="ts">
	import DataTable from './breeds-data-grid/data-table.svelte';
	import { columns } from './breeds-data-grid/columns';
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconPlus from '@lucide/svelte/icons/plus';

	let { data }: PageProps = $props();
</script>

{#if data.breedsResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Breeds</h2>
				<p class="text-muted-foreground">List of breeds</p>
			</div>
		</div>
		<div>
			<a href="/dashboard/breeds/create">
				<Button>
					Create new breed
					<IconPlus />
				</Button>
			</a>
		</div>
		<DataTable
			data={data.breedsResult?.data.data ?? []}
			meta={data.breedsResult?.data.meta}
			{columns}
			animalSpecies={data.animalSpeciesResults?.data?.data ?? []}
		/>
	</div>
{/if}
