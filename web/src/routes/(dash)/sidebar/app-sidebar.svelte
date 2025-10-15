<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import IconChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import IconLogout from '@lucide/svelte/icons/log-out';
	import { dashboardNavItems } from './dashboard-nav-items';
	import type { User } from '$lib/kratos/service';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { page } from '$app/state';

	let {
		ref = $bindable(null),
		user,
		logoutUrl,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & { user?: User; logoutUrl?: string } = $props();

	const sidebar = Sidebar.useSidebar();
</script>

{#if user}
	<Sidebar.Root {...restProps} bind:ref>
		<Sidebar.Header>
			<div class="flex h-14 items-center justify-center gap-4">
				<img src="/images/logo.svg" alt="Fluffly logo" class="m-0 h-full max-w-full p-0" />
				<span>Fluffly</span>
			</div>
		</Sidebar.Header>

		<Sidebar.Content>
			{#each dashboardNavItems as group (group.title)}
				<Sidebar.Group>
					<Sidebar.GroupLabel>{group.title}</Sidebar.GroupLabel>
					<Sidebar.GroupContent>
						<Sidebar.Menu>
							{#each group.items as item (item.title)}
								<Sidebar.MenuItem>
									<Sidebar.MenuButton isActive={page.url.pathname === item.url}>
										{#snippet child({ props })}
											<a href={item.url} {...props}>
												{#if item.icon}
													{@const Icon = item.icon}
													<Icon />
												{/if}
												{item.title}
											</a>
										{/snippet}
									</Sidebar.MenuButton>
								</Sidebar.MenuItem>
							{/each}
						</Sidebar.Menu>
					</Sidebar.GroupContent>
				</Sidebar.Group>
			{/each}
		</Sidebar.Content>
		<Sidebar.Footer>
			{@render sidebarFooterContent()}
		</Sidebar.Footer>
		<Sidebar.Rail />
	</Sidebar.Root>
{:else}
	<Sidebar.Root {...restProps} bind:ref>
		<Sidebar.Header>
			<div class="flex h-14 items-center justify-center gap-4">
				<img src="/images/logo.svg" alt="Fluffly logo" class="m-0 h-full max-w-full p-0" />
				<span>Fluffly</span>
			</div>
		</Sidebar.Header>

		<Sidebar.Content>
			{#each dashboardNavItems as group (group.title)}
				<Sidebar.Group>
					<Sidebar.GroupLabel>
						<Skeleton class="h-4 w-24" />
					</Sidebar.GroupLabel>
					<Sidebar.GroupContent>
						<Sidebar.Menu>
							{#each group.items as item (item.title)}
								<Sidebar.MenuItem>
									<Sidebar.MenuButton>
										{#snippet child()}
											<Skeleton class="h-4 w-full" />
										{/snippet}
									</Sidebar.MenuButton>
								</Sidebar.MenuItem>
							{/each}
						</Sidebar.Menu>
					</Sidebar.GroupContent>
				</Sidebar.Group>
			{/each}
		</Sidebar.Content>
		<Sidebar.Footer>
			<Skeleton class="h-12 w-full" />
		</Sidebar.Footer>
		<Sidebar.Rail />
	</Sidebar.Root>
{/if}

{#snippet sidebarFooterContent()}
	<Sidebar.Menu>
		<Sidebar.MenuItem>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuButton
							size="lg"
							class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							{...props}
						>
							<Avatar.Root class="h-8 w-8 rounded-lg">
								<Avatar.Image src={user?.avatarUrl} alt={user?.fullName} />
								<Avatar.Fallback class="rounded-lg">CN</Avatar.Fallback>
							</Avatar.Root>
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-semibold">{user?.fullName}</span>
								<span class="truncate text-xs">{user?.email}</span>
							</div>
							<IconChevronsUpDown class="ml-auto size-4" />
						</Sidebar.MenuButton>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content
					class="w-[--bits-dropdown-menu-anchor-width] min-w-56 rounded-lg"
					side={sidebar.isMobile ? 'bottom' : 'right'}
					align="end"
					sideOffset={4}
				>
					<DropdownMenu.Label class="p-0 font-normal">
						<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
							<Avatar.Root class="h-8 w-8 rounded-lg">
								<Avatar.Image src={user?.avatarUrl} alt={user?.fullName} />
								<Avatar.Fallback class="rounded-lg">CN</Avatar.Fallback>
							</Avatar.Root>
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-semibold">{user?.fullName}</span>
								<span class="truncate text-xs">{user?.email}</span>
							</div>
						</div>
					</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<a href={logoutUrl}>
						<DropdownMenu.Item class="cursor-pointer">
							<IconLogout />
							Log out
						</DropdownMenu.Item>
					</a>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</Sidebar.MenuItem>
	</Sidebar.Menu>
{/snippet}
