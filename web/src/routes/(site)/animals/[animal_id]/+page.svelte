<script lang="ts">
	import type { PageProps } from './$types';
	import AnimalPage from './animal.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Animal {params.animal_id}</title>
</svelte:head>

<AnimalPage {data} {params} />
