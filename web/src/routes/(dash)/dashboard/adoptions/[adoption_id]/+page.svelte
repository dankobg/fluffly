<script lang="ts">
	import type { PageProps } from './$types';
	import ViewAdoption from './view-adoption.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View adoption</title>
</svelte:head>

<ViewAdoption {data} {params} />
