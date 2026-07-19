<script lang="ts">
	import type { PageProps } from './$types';
	import PostedAnimals from './posted-animals.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>My posted animals</title>
</svelte:head>

<PostedAnimals {data} {params} />
