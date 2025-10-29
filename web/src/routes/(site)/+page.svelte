<script lang="ts">
	import type { PageProps } from './$types';
	import Home from './home.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Fluffly home</title>
</svelte:head>

<Home {data} {params} />
