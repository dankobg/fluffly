<script lang="ts">
	import type { LayoutData } from '../$types';
	import { type Snippet } from 'svelte';
	import Header from '$lib/components/header/header.svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<!-- <Header user={data.auth.user ?? undefined} logoutUrl={data.logoutUrl} /> -->
<Header user={undefined} logoutUrl={undefined} />
{@render children()}
