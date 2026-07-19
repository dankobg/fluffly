<script lang="ts">
	import IconUser from '@lucide/svelte/icons/user';
	import IconChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import IconLogout from '@lucide/svelte/icons/log-out';
	import IconDashboard from '@lucide/svelte/icons/layout-dashboard';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { User } from '$lib/kratos/service';
	import { getInitials } from '$lib/utils';

	let { logoutUrl, user }: { logoutUrl?: string; user: User } = $props();

	const sidebar = Sidebar.useSidebar();
</script>

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
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={user?.avatarUrl} alt={user?.fullName} />
							<Avatar.Fallback class="rounded-lg">{getInitials(user?.fullName || user?.email)}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{user?.fullName}</span>
							<span class="truncate text-xs">{user?.email}</span>
						</div>
						<IconChevronsUpDown class="ms-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Group>
					<a href="/dashboard">
						<DropdownMenu.Item>
							<IconDashboard />
							Dashboard
						</DropdownMenu.Item>
					</a>
				</DropdownMenu.Group>
				<DropdownMenu.Group>
					<a href="/dashboard/account">
						<DropdownMenu.Item>
							<IconUser />
							Account
						</DropdownMenu.Item>
					</a>
				</DropdownMenu.Group>
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
