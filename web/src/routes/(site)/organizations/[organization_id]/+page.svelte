<script lang="ts">
	import type { PageProps } from './$types';
	import OrganizationPage from './organization.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Organization {params.organization_id}</title>
</svelte:head>

<OrganizationPage {data} {params} />
