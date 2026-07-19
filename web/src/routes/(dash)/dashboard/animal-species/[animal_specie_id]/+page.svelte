<script lang="ts">
	import type { PageProps } from './$types';
	import ViewAnimalSpecie from './view-animal-specie.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View animal species</title>
</svelte:head>

<ViewAnimalSpecie {data} {params} />
