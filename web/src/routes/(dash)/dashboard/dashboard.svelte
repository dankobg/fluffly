<script lang="ts">
	import type { PageProps } from './$types';
	import AdminDashboard from './admin-dashboard.svelte';
	import UserDashboard from './user-dashboard.svelte';

	let { data, params }: PageProps = $props();
</script>

{#if data.auth?.user?.isDeveloper}
	<AdminDashboard {data} {params} />
{:else}
	<UserDashboard {data} {params} />
{/if}
