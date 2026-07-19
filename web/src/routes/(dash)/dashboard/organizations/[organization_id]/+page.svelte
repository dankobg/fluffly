<script lang="ts">
	import type { PageProps } from './$types';
	import ViewOrganization from './view-organization.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>View organization</title>
</svelte:head>

<ViewOrganization {data} {params} />
