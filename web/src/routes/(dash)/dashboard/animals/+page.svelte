<script lang="ts">
	import DataTable from './animals-data-grid/data-table.svelte';
	import { columns } from './animals-data-grid/columns';
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconPlus from '@lucide/svelte/icons/plus';

	let { data }: PageProps = $props();
</script>

{#if data.animalsResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Animals</h2>
				<p class="text-muted-foreground">List of Animals</p>
			</div>
		</div>
		<div>
			<a href="/dashboard/animals/create">
				<Button>
					Create new animal
					<IconPlus />
				</Button>
			</a>
		</div>
		<DataTable
			data={data.animalsResult?.data.data ?? []}
			meta={data.animalsResult?.data.meta}
			{columns}
			animalTypes={data.animalTypesResult?.data?.data ?? []}
			animalSpecies={data.animalSpeciesResult?.data?.data ?? []}
		/>
	</div>
{/if}
