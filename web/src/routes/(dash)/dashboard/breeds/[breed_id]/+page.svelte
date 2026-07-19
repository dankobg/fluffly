<script lang="ts">
	import type { PageProps } from './$types';
	import ViewBreed from './view-breed.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View breed</title>
</svelte:head>

<ViewBreed {data} {params} />
