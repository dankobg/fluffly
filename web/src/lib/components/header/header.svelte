<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import * as Avatar from '$lib/components/ui/avatar/index';
	import { config } from '$lib/kratos/config';
	import { page } from '$app/state';
	import { navItems } from './nav-items';
	import ModeSwitcher from '$lib/components/mode-switcher/mode-switcher.svelte';
	import IconLogout from '@lucide/svelte/icons/log-out';
	import IconUserIcon from '@lucide/svelte/icons/user';
	import IconGauge from '@lucide/svelte/icons/gauge';
	import IconMenu from '@lucide/svelte/icons/menu';
	import * as Sheet from '$lib/components/ui/sheet';
	import type { User } from '$lib/kratos/service';
	import { getInitials } from '$lib/utils';

	type Props = {
		logoutUrl?: string;
		user?: User;
	};

	let { logoutUrl, user }: Props = $props();
</script>

<header class="bg-background sticky top-0 flex h-16 max-w-[1920px] items-center gap-4 px-4 md:mx-auto md:px-6">
	<nav class="flex w-full items-center justify-start text-lg font-medium md:gap-5 lg:gap-6">
		<Sheet.Root>
			<Sheet.Trigger class="flex md:hidden">
				<IconMenu class="mr-4 size-6" />
				<span class="sr-only">Open sidebar navigation</span>
			</Sheet.Trigger>
			<Sheet.Content side="left">
				<nav class="grid gap-6 text-lg font-medium">
					<a href="/" class="flex items-center gap-2 text-lg font-semibold">
						<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
						<span class="text-lg">Fluffly</span>
					</a>
					{#each navItems as navItem (navItem.label)}
						<a
							href={navItem.href}
							class="text-muted-foreground hover:text-primary transition-all"
							class:text-primary={page.url.pathname === navItem.href}>{navItem.label}</a
						>
					{/each}
				</nav>
			</Sheet.Content>
		</Sheet.Root>

		<a href="/" class="flex items-center gap-2 font-semibold md:text-base">
			<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
			<span class="text-lg">Fluffly</span>
		</a>

		<div class="hidden items-center gap-2 md:flex">
			{#each navItems as item (item.label)}
				<a
					href={item.href}
					class="text-muted-foreground hover:text-primary text-lg transition-all"
					class:text-primary={page.url.pathname === item.href}>{item.label}</a
				>
			{/each}
		</div>
	</nav>

	<div class="ml-auto flex items-center gap-4 md:gap-2 lg:gap-4">
		<ModeSwitcher />
	</div>

	{#if user}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				<Avatar.Root>
					<Avatar.Image src={user.avatarUrl} alt="user avatar" />
					<Avatar.Fallback class="text-md bg-secondary font-bold">
						{getInitials(user.fullName ?? user.email)}
					</Avatar.Fallback>
				</Avatar.Root>
				<span class="sr-only">Toggle user menu</span>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content class="w-56">
				<DropdownMenu.Group>
					<DropdownMenu.GroupHeading>{user.fullName ?? user.email}</DropdownMenu.GroupHeading>
					<DropdownMenu.Separator />
					<DropdownMenu.Group>
						<a href="/dashboard">
							<DropdownMenu.Item class="cursor-pointer">
								<IconGauge class="mr-2 size-4" />
								<span>Dashboard</span>
							</DropdownMenu.Item>
						</a>
						<a href="/dashboard/account">
							<DropdownMenu.Item class="cursor-pointer">
								<IconUserIcon class="mr-2 size-4" />
								<span>Account</span>
							</DropdownMenu.Item>
						</a>
					</DropdownMenu.Group>
					<DropdownMenu.Separator />
					<a href={logoutUrl}>
						<DropdownMenu.Item class="cursor-pointer">
							<IconLogout class="mr-2 size-4" />
							<span>Log out</span>
						</DropdownMenu.Item>
					</a>
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{:else}
		<a href={config.routes.login.path} class="text-muted-foreground hover:text-primary text-lg transition-all">Login</a>
		<a href={config.routes.registration.path} class="text-muted-foreground hover:text-primary text-lg transition-all">
			Register
		</a>
	{/if}
</header>
