<script lang="ts">
	import DataTable from './identity-sessions-data-grid/data-table.svelte';
	import { columns } from './identity-sessions-data-grid/columns';
	import type { PageProps } from './$types';
	import type { CustomTraits } from '$lib/kratos/service';

	let { data, params }: PageProps = $props();
</script>

<div class="flex h-full flex-1 flex-col space-y-8 p-8">
	<div class="flex items-center justify-between space-y-2">
		<div>
			<h2 class="text-2xl font-bold tracking-tight">
				Identity sessions - <strong>{(data.identityResult?.data?.traits as CustomTraits)?.['email']}</strong>
			</h2>
			<p class="text-muted-foreground">List of sessions</p>
		</div>
	</div>
	<DataTable data={data.sessionsResult?.data ?? []} {columns} identity={data.identityResult?.data} />
</div>
