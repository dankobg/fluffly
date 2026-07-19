<script lang="ts">
	import DataTable from './sessions-data-grid/data-table.svelte';
	import { columns, columnsWithIdentity } from './sessions-data-grid/columns';
	import type { PageProps } from './$types';
	import { expandSessions } from './expandSessions.svelte';
	import { PathsSessionsGetParametersQueryExpand } from '$lib/gen/fluffly_openapi';

	let { data }: PageProps = $props();

	let cols = $derived(
		expandSessions.expanded.includes(PathsSessionsGetParametersQueryExpand.identity) ? columnsWithIdentity : columns
	);
</script>

{#if data.sessionsResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Sessions</h2>
				<p class="text-muted-foreground">List of sessions</p>
			</div>
		</div>
		<DataTable data={data.sessionsResult?.data ?? []} columns={cols} />
	</div>
{/if}
