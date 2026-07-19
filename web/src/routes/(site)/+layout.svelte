<script lang="ts">
	import type { LayoutData } from '../$types';
	import { type Snippet } from 'svelte';
	import Header from '$lib/components/shell/header.svelte';
	import SiteSidebar from '$lib/components/shell/site-sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col" open={false}>
		<Header user={data.auth.user} logoutUrl={data.logoutUrl} />
		<div class="flex flex-1">
			<SiteSidebar user={data?.auth?.user} logoutUrl={data?.logoutUrl} />
			<Sidebar.Inset class="mx-auto max-w-[120rem]">
				{@render children()}
			</Sidebar.Inset>
		</div>
	</Sidebar.Provider>
</div>
