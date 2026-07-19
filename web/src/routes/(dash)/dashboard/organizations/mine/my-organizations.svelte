<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index';
	import * as Item from '$lib/components/ui/item/index.js';
	import IconBadgeCheck from '@lucide/svelte/icons/badge-check';
	import IconClock from '@lucide/svelte/icons/clock';
	import IconCircleX from '@lucide/svelte/icons/circle-x';
	import IconChevronRight from '@lucide/svelte/icons/chevron-right';
	import IconSquareDashed from '@lucide/svelte/icons/square-dashed';
	import { OrganizationStatus } from '$lib/gen/fluffly_openapi';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { cn } from '$lib/utils';

	let { data }: PageProps = $props();
</script>

<div class="grid grid-cols-[repeat(auto-fill,minmax(min(22rem,100%),1fr))] justify-center gap-4">
	<Card.Root>
		<Card.Header>
			<Card.Title>My organizations</Card.Title>
		</Card.Header>
		<Card.Content class="grid gap-4">
			{#if data.myOrganizationsResult?.data?.data && data.myOrganizationsResult?.data?.data.length > 0}
				{#each data.myOrganizationsResult?.data?.data as org (org.id)}
					<Item.Root variant="outline" size="sm">
						{#snippet child({ props })}
							<a href="/dashboard/organizations/{org.id}" {...props}>
								<Item.Media>
									{#if org.status === OrganizationStatus.approved}
										<IconBadgeCheck class="size-5 text-green-600" />
									{/if}
									{#if org.status === OrganizationStatus.pending}
										<IconClock class="size-5 text-yellow-600" />
									{/if}
									{#if org.status === OrganizationStatus.rejected}
										<IconCircleX class="size-5 text-red-600" />
									{/if}
								</Item.Media>
								<Item.Content>
									<Item.Title>
										{org.name}
										<Badge
											class={cn(
												org.status === OrganizationStatus.approved
													? 'bg-green-600'
													: org.status === OrganizationStatus.pending
														? 'bg-yellow-600'
														: org.status === OrganizationStatus.rejected
															? 'bg-red-600'
															: 'bg-gray-500'
											)}
										>
											{org.status}
										</Badge>
									</Item.Title>
								</Item.Content>
								<Item.Actions>
									<IconChevronRight class="size-4" />
								</Item.Actions>
							</a>
						{/snippet}
					</Item.Root>
				{/each}
			{:else}
				<Empty.Root class="flex-none border">
					<Empty.Header>
						<Empty.Media variant="icon">
							<IconSquareDashed />
						</Empty.Media>
						<Empty.Title>Not part of any organizations</Empty.Title>
						<Empty.Description>Apply for an organization so you can adopt animals faster</Empty.Description>
					</Empty.Header>
					<Empty.Content>
						<a href="/dashboard/organizations/apply" class="underline">Apply for an organization now</a>
					</Empty.Content>
				</Empty.Root>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
