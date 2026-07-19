<script lang="ts">
	import type { PageProps } from './$types';
	import AdoptedAnimals from './adopted-animals.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>My adoptions</title>
</svelte:head>

<AdoptedAnimals {data} {params} />
