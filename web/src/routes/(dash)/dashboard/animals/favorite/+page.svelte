<script lang="ts">
	import type { PageProps } from './$types';
	import FavoriteAnimals from './favorite-animals.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Favorite animals</title>
</svelte:head>

<FavoriteAnimals {data} {params} />
