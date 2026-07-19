<script lang="ts">
	import DataTable from './countries-data-grid/data-table.svelte';
	import { columns } from './countries-data-grid/columns';
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconPlus from '@lucide/svelte/icons/plus';

	let { data }: PageProps = $props();
</script>

{#if data.countriesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Countries</h2>
				<p class="text-muted-foreground">List of Countries</p>
			</div>
		</div>
		<div>
			<a href="/dashboard/countries/create">
				<Button>
					Create new country
					<IconPlus />
				</Button>
			</a>
		</div>
		<DataTable data={data.countriesResult?.data.data ?? []} meta={data.countriesResult?.data.meta} {columns} />
	</div>
{/if}
