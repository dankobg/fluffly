<script lang="ts">
	import DataTable from './animal-species-data-grid/data-table.svelte';
	import { columns } from './animal-species-data-grid/columns';
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconPlus from '@lucide/svelte/icons/plus';

	let { data }: PageProps = $props();
</script>

{#if data.animalSpeciesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Animal species</h2>
				<p class="text-muted-foreground">List of animal species</p>
			</div>
		</div>
		<div>
			<a href="/dashboard/animal-species/create">
				<Button>
					Create new animal species
					<IconPlus />
				</Button>
			</a>
		</div>
		<DataTable
			data={data.animalSpeciesResult?.data.data ?? []}
			meta={data.animalSpeciesResult?.data.meta}
			{columns}
			animalTypes={data?.animalTypesResult?.data?.data ?? []}
		/>
	</div>
{/if}
