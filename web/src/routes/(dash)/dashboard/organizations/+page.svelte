<script lang="ts">
	import DataTable from './organizations-data-grid/data-table.svelte';
	import { columns } from './organizations-data-grid/columns';
	import type { PageProps } from './$types';
	import Button from '$lib/components/ui/button/button.svelte';
	import IconPlus from '@lucide/svelte/icons/plus';

	let { data }: PageProps = $props();
</script>

{#if data.organizationsResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Organizations</h2>
				<p class="text-muted-foreground">List of Organizations</p>
			</div>
		</div>
		<div>
			<a href="/dashboard/organizations/create">
				<Button>
					Create new organization
					<IconPlus />
				</Button>
			</a>
		</div>
		<DataTable data={data.organizationsResult?.data.data ?? []} meta={data.organizationsResult?.data.meta} {columns} />
	</div>
{/if}
