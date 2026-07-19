<script lang="ts">
	import type { PageProps } from './$types';
	import UpdateOrganization from './update-organization.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Update organization</title>
</svelte:head>

<UpdateOrganization {data} {params} />
