<script lang="ts">
	import type { PageProps } from './$types';

	let { data, params }: PageProps = $props();
</script>

{#if data.schema}
	<h1 class="mb-6 text-2xl font-bold">Identity schema - {data.id}</h1>
	<div>
		<pre>{JSON.stringify(data.schema, null, 2)}</pre>
	</div>
{/if}
