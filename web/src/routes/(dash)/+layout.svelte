<script lang="ts">
	import type { LayoutData } from '../$types';
	import { type Snippet } from 'svelte';
	import Header from './dashboard-shell/header.svelte';
	import SiteSidebar from './dashboard-shell/site-sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { MediaQuery } from 'svelte/reactivity';
	import { page } from '$app/state';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index';
	import * as Drawer from '$lib/components/ui/drawer/index';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import { buttonVariants } from '$lib/components/ui/button/index';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	const ITEMS_TO_DISPLAY = 3;
	let open = $state(false);
	const isDesktop = new MediaQuery('(min-width: 768px)');
	let crumbs: Array<{ label: string; href: string }> = $state([]);
	const homeCrumb = { label: 'Home', href: '/' };

	$effect(() => {
		const tokens = page.url.pathname.split('/').filter(t => t !== '');
		let href = '';
		const arr = tokens.map(t => {
			href += '/' + t;
			const label = t.charAt(0).toUpperCase() + t.slice(1);
			return {
				label: page.data['label'] || label,
				href
			};
		});
		crumbs = [homeCrumb, ...arr];
	});
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col" open>
		<Header user={data.auth.user} logoutUrl={data.logoutUrl} />
		<div class="flex flex-1">
			<SiteSidebar
				user={data?.auth?.user}
				logoutUrl={data?.logoutUrl}
				isDeveloper={Boolean(data.auth.user?.isDeveloper)}
			/>
			<Sidebar.Inset class="mx-auto max-w-[120rem]">
				{@render dashboardBreadcrumbs()}

				<div class="flex flex-1 flex-col p-4">
					<div class="@container/main flex flex-1 flex-col gap-4 rounded-xl bg-muted/50 p-4">
						{#if data.auth.session?.active}
							{@render children()}
						{:else}
							<span>Loading...</span>
						{/if}
					</div>
				</div>
			</Sidebar.Inset>
		</div>
	</Sidebar.Provider>
</div>

{#snippet dashboardBreadcrumbs()}
	<div class="flex h-10 shrink-0 items-center gap-2 border-b px-4">
		{#if data.auth.session?.active}
			<Breadcrumb.Root>
				<Breadcrumb.List>
					{#if crumbs.length > 2}
						<Breadcrumb.Item>
							<Breadcrumb.Link href={crumbs?.[0]?.href}>
								{crumbs?.[0]?.label}
							</Breadcrumb.Link>
						</Breadcrumb.Item>
						<Breadcrumb.Separator />
					{/if}

					{#if crumbs.length > ITEMS_TO_DISPLAY}
						<Breadcrumb.Item>
							{#if isDesktop.current}
								<DropdownMenu.Root bind:open>
									<DropdownMenu.Trigger class="flex items-center gap-1" aria-label="Toggle menu">
										<Breadcrumb.Ellipsis class="size-4" />
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="start">
										{#each crumbs.slice(1, -2) as item (item.label)}
											<DropdownMenu.Item>
												<a href={item.href ? item.href : '#'}>
													{item.label}
												</a>
											</DropdownMenu.Item>
										{/each}
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							{:else}
								<Drawer.Root bind:open>
									<Drawer.Trigger aria-label="Toggle Menu">
										<Breadcrumb.Ellipsis class="size-4" />
									</Drawer.Trigger>
									<Drawer.Content>
										<Drawer.Header class="text-left">
											<Drawer.Title>Navigate to</Drawer.Title>
											<Drawer.Description>Select a page to navigate to.</Drawer.Description>
										</Drawer.Header>
										<div class="grid gap-1 px-4">
											{#each crumbs.slice(1, -2) as item (item.label)}
												<a href={item.href ? item.href : '#'} class="py-1 text-sm">
													{item.label}
												</a>
											{/each}
										</div>
										<Drawer.Footer class="pt-4">
											<Drawer.Close class={buttonVariants({ variant: 'outline' })}>Close</Drawer.Close>
										</Drawer.Footer>
									</Drawer.Content>
								</Drawer.Root>
							{/if}
						</Breadcrumb.Item>
						<Breadcrumb.Separator />
					{/if}

					{#each crumbs.slice(-ITEMS_TO_DISPLAY + 1) as item (item.label)}
						<Breadcrumb.Item>
							{#if item.href}
								<Breadcrumb.Link href={item.href} class="max-w-20 truncate md:max-w-none">
									{item.label}
								</Breadcrumb.Link>
								<Breadcrumb.Separator />
							{:else}
								<Breadcrumb.Page class="max-w-20 truncate md:max-w-none">
									{item.label}
								</Breadcrumb.Page>
							{/if}
						</Breadcrumb.Item>
					{/each}
				</Breadcrumb.List>
			</Breadcrumb.Root>
		{:else}
			<Skeleton class="h-4 w-48" />
		{/if}
	</div>
{/snippet}
