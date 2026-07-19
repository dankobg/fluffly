<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { mainNavItems } from './nav-items';
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>Site navigation</Sidebar.GroupLabel>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			{#each mainNavItems as mainNavItem (mainNavItem.title)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton isActive={mainNavItem.isActive}>
						{#snippet child({ props })}
							<a href={mainNavItem.url} {...props}>{mainNavItem.title}</a>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Menu>
	</Sidebar.GroupContent>
</Sidebar.Group>
