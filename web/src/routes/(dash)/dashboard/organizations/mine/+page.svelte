<script lang="ts">
	import type { PageProps } from './$types';
	import MyOrganizations from './my-organizations.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>My organizations</title>
</svelte:head>

<MyOrganizations {data} {params} />
