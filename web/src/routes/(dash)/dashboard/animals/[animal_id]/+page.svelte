<script lang="ts">
	import type { PageProps } from './$types';
	import ViewAnimal from './view-animal.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View animal</title>
</svelte:head>

<ViewAnimal {data} {params} />
