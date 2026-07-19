<script lang="ts">
	import { type NavItem } from './nav-items';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';

	let { items }: { items: NavItem[] } = $props();
</script>

{#each items as item (item.title)}
	<Collapsible.Root title={item.title} open class="group/collapsible">
		<Sidebar.Group>
			<Sidebar.GroupLabel
				class="group/label text-sm text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
			>
				{#snippet child({ props })}
					<Collapsible.Trigger {...props}>
						{item.title}
						<ChevronRightIcon class="ms-auto transition-transform group-data-[state=open]/collapsible:rotate-90" />
					</Collapsible.Trigger>
				{/snippet}
			</Sidebar.GroupLabel>
			<Collapsible.Content>
				<Sidebar.GroupContent>
					<Sidebar.Menu>
						{#each item.items as subItem (subItem.title)}
							<Sidebar.MenuItem>
								<Sidebar.MenuButton isActive={subItem.isActive}>
									{#snippet child({ props })}
										<a href={subItem.url} {...props}>{subItem.title}</a>
									{/snippet}
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/each}
					</Sidebar.Menu>
				</Sidebar.GroupContent>
			</Collapsible.Content>
		</Sidebar.Group>
	</Collapsible.Root>
{/each}
