<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import * as Avatar from '$lib/components/ui/avatar/index';
	import { config } from '$lib/kratos/config';
	import ModeSwitcher from '$lib/components/mode-switcher/mode-switcher.svelte';
	import IconLogout from '@lucide/svelte/icons/log-out';
	import IconUserIcon from '@lucide/svelte/icons/user';
	import IconGauge from '@lucide/svelte/icons/gauge';
	import IconSidebar from '@lucide/svelte/icons/sidebar';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { User } from '$lib/kratos/service';
	import { getInitials } from '$lib/utils';

	type Props = {
		logoutUrl?: string;
		user: User | null;
	};

	let { logoutUrl, user }: Props = $props();

	const sidebar = Sidebar.useSidebar();
</script>

<header class="sticky top-0 z-50 mx-auto flex w-full max-w-[120rem] items-center border-b bg-background">
	<div class="flex h-(--header-height) w-full items-center gap-2 px-4">
		<Button class="size-8" variant="ghost" size="icon" onclick={sidebar.toggle}>
			<IconSidebar />
		</Button>
		<Separator orientation="vertical" class="me-2 h-4" />

		<a href="/" class="flex items-center gap-2 font-semibold md:text-base">
			<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
			<span class="text-lg">Fluffly</span>
		</a>

		<div class="ml-auto hidden items-center gap-4 sm:flex md:gap-2 lg:gap-4">
			<ModeSwitcher />
		</div>

		{#if user}
			<DropdownMenu.Root>
				<DropdownMenu.Trigger class="hidden sm:block">
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
			<a
				href={config.routes.login.path}
				class="hidden text-lg text-muted-foreground transition-all hover:text-primary sm:inline"
			>
				Login
			</a>
			<a
				href={config.routes.registration.path}
				class="hidden text-lg text-muted-foreground transition-all hover:text-primary sm:inline"
			>
				Register
			</a>
		{/if}
	</div>
</header>
