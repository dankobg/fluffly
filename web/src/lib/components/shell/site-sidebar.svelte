<script lang="ts" module>
</script>

<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import IconEllipsis from '@lucide/svelte/icons/ellipsis';
	import type { User } from '$lib/kratos/service';
	import { config } from '$lib/kratos/config';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { resetMode, setMode } from 'mode-watcher';

	let {
		logoutUrl,
		user,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & { logoutUrl?: string; user: User | null } = $props();

	const sidebar = Sidebar.useSidebar();
</script>

<Sidebar.Root class="top-(--header-height) h-[calc(100svh-var(--header-height))]!" {...restProps}>
	<Sidebar.Content>
		<NavMain />
	</Sidebar.Content>
	<Sidebar.Footer>
		{@render modeSwitcher()}
		{@render footerContent()}
	</Sidebar.Footer>
</Sidebar.Root>

{#snippet modeSwitcher()}
	<Sidebar.Menu>
		<DropdownMenu.Root>
			<Sidebar.MenuItem>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuButton
							class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							{...props}
						>
							Mode switcher
							<IconEllipsis class="ms-auto" />
						</Sidebar.MenuButton>
					{/snippet}
				</DropdownMenu.Trigger>

				<DropdownMenu.Content
					class="min-w-56 rounded-lg"
					side={sidebar.isMobile ? 'bottom' : 'right'}
					align={sidebar.isMobile ? 'end' : 'start'}
				>
					<DropdownMenu.Item onclick={() => setMode('light')}>Light</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => setMode('dark')}>Dark</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => resetMode()}>System</DropdownMenu.Item>
				</DropdownMenu.Content>
			</Sidebar.MenuItem>
		</DropdownMenu.Root>
	</Sidebar.Menu>
{/snippet}

{#snippet footerContent()}
	{#if user}
		<NavUser {user} {logoutUrl} />
	{:else}
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<a href={config.routes.login.path} {...props}>Login</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<a href={config.routes.registration.path} {...props}>Register</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	{/if}
{/snippet}
