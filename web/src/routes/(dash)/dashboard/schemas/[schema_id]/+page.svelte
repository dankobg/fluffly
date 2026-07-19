<script lang="ts">
	import type { PageProps } from './$types';
	import ViewSchema from './view-schema.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View identity schema</title>
</svelte:head>

<ViewSchema {data} {params} />
