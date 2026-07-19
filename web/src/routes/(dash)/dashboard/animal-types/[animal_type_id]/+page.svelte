<script lang="ts">
	import type { PageProps } from './$types';
	import ViewAnimalType from './view-animal-type.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View animal type</title>
</svelte:head>

<ViewAnimalType {data} {params} />
