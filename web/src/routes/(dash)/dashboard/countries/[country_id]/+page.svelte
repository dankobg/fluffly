<script lang="ts">
	import type { PageProps } from './$types';
	import ViewCountry from './view-country.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View country</title>
</svelte:head>

<ViewCountry {data} {params} />
