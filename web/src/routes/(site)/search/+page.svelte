<script lang="ts">
	import type { PageProps } from './$types';
	import SearchPage from './search.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Fluffly search</title>
</svelte:head>

<SearchPage {data} {params} />
