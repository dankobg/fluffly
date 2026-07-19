<script lang="ts">
	import type { PageProps } from './$types';
	import UpdateAnimal from './update-animal.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Update animal</title>
</svelte:head>

<UpdateAnimal {data} {params} />
