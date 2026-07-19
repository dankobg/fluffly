<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index';
	import { PieChart, Text } from 'layerchart';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import IconUsers from '@lucide/svelte/icons/users';
	import IconCat from '@lucide/svelte/icons/cat';
	import IconBuilding2 from '@lucide/svelte/icons/building-2';
	import IconHandHeart from '@lucide/svelte/icons/hand-heart';

	let { data }: PageProps = $props();

	const chartConfigAdoptions = {
		approved: { label: 'Approved', color: 'var(--color-green-600)' },
		pending: { label: 'Pending', color: 'var(--color-yellow-600)' },
		rejected: { label: 'Rejected', color: 'var(--color-red-600)' }
	} satisfies Chart.ChartConfig;

	const chartConfigOrganizations = {
		approved: { label: 'Approved', color: 'var(--color-green-600)' },
		pending: { label: 'Pending', color: 'var(--color-yellow-600)' },
		rejected: { label: 'Rejected', color: 'var(--color-red-600)' }
	} satisfies Chart.ChartConfig;

	const chartConfigAnimals = {
		adoptable: { label: 'Adoptable', color: 'var(--color-sky-600)' },
		pending: { label: 'Pending', color: 'var(--color-yellow-600)' },
		adopted: { label: 'Adopted', color: 'var(--color-green-600)' },
		reserved: { label: 'Reserved', color: 'var(--color-sky-600)' },
		rejected: { label: 'Rejected', color: 'var(--color-red-600)' }
	} satisfies Chart.ChartConfig;

	let chartDataAdoptions = $derived.by(() => {
		if (!data?.myAnalyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigOrganizations.approved.label,
				value: data.myAnalyticsStatsResult.data.stats.adoptions.by_status.approved,
				color: chartConfigOrganizations.approved.color
			},
			{
				label: chartConfigOrganizations.pending.label,
				value: data.myAnalyticsStatsResult.data.stats.adoptions.by_status.pending,
				color: chartConfigOrganizations.pending.color
			},
			{
				label: chartConfigOrganizations.rejected.label,
				value: data.myAnalyticsStatsResult.data.stats.adoptions.by_status.rejected,
				color: chartConfigOrganizations.rejected.color
			}
		];
	});

	let chartDataOrganizations = $derived.by(() => {
		if (!data?.myAnalyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigOrganizations.approved.label,
				value: data.myAnalyticsStatsResult.data.stats.organizations.by_status.approved,
				color: chartConfigOrganizations.approved.color
			},
			{
				label: chartConfigOrganizations.pending.label,
				value: data.myAnalyticsStatsResult.data.stats.organizations.by_status.pending,
				color: chartConfigOrganizations.pending.color
			},
			{
				label: chartConfigOrganizations.rejected.label,
				value: data.myAnalyticsStatsResult.data.stats.organizations.by_status.rejected,
				color: chartConfigOrganizations.rejected.color
			}
		];
	});

	let chartDataAnimals = $derived.by(() => {
		if (!data?.myAnalyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigAnimals.adoptable.label,
				value: data.myAnalyticsStatsResult.data.stats.animals.posted.by_status.adoptable,
				color: chartConfigAnimals.adoptable.color
			},
			{
				label: chartConfigAnimals.pending.label,
				value: data.myAnalyticsStatsResult.data.stats.animals.posted.by_status.pending,
				color: chartConfigAnimals.pending.color
			},
			{
				label: chartConfigAnimals.adopted.label,
				value: data.myAnalyticsStatsResult.data.stats.animals.posted.by_status.adopted,
				color: chartConfigAnimals.adopted.color
			},
			{
				label: chartConfigAnimals.reserved.label,
				value: data.myAnalyticsStatsResult.data.stats.animals.posted.by_status.reserved,
				color: chartConfigAnimals.reserved.color
			},
			{
				label: chartConfigAnimals.rejected.label,
				value: data.myAnalyticsStatsResult.data.stats.animals.posted.by_status.rejected,
				color: chartConfigAnimals.rejected.color
			}
		];
	});

	let favoriteAnimalsDesc = $derived.by(() => {
		const favorites = data?.myAnalyticsStatsResult?.data?.stats?.animals?.favorites || 0;
		if (favorites > 0) {
			return 'Possible future best friends';
		} else {
			return 'You have not found your perfect friend yet';
		}
	});

	let favoriteAnimalsText = $derived.by(() => {
		const favorites = data?.myAnalyticsStatsResult?.data?.stats?.animals?.favorites || 0;
		if (favorites === 0) return 'No liked animals';
		if (favorites < 5) return 'Few liked animals';
		if (favorites < 10) return 'Several liked animals';
		return 'Many liked animals';
	});
</script>

{#if data?.myAnalyticsStatsResult?.data}
	<div class="grid grid-cols-1 gap-4 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Adoptions</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.myAnalyticsStatsResult?.data?.stats?.adoptions.total}
				</Card.Title>
				<Card.Action>
					<IconHandHeart />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.myAnalyticsStatsResult?.data?.stats?.adoptions?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.myAnalyticsStatsResult?.data?.stats?.organizations?.by_status.pending > 0
						? 'Waiting for admin approval'
						: 'No pending adoptions'}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Favorites</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.myAnalyticsStatsResult?.data?.stats?.animals.favorites}
				</Card.Title>
				<Card.Action>
					<IconUsers />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">{favoriteAnimalsDesc}</div>
				<div class="text-muted-foreground">
					{favoriteAnimalsText}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Animals</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.myAnalyticsStatsResult?.data?.stats?.animals.posted.total}
				</Card.Title>
				<Card.Action>
					<IconCat />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.myAnalyticsStatsResult?.data?.stats?.animals?.posted.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.myAnalyticsStatsResult?.data?.stats?.animals?.posted.by_status.pending > 0
						? 'Waiting for admin approval'
						: 'No pending animals'}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Organizations</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.myAnalyticsStatsResult?.data?.stats?.organizations?.total}
				</Card.Title>
				<Card.Action>
					<IconBuilding2 />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.myAnalyticsStatsResult?.data?.stats?.organizations?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.myAnalyticsStatsResult?.data?.stats?.organizations?.by_status.pending > 0
						? 'Waiting for admin approval'
						: 'No pending organizations'}
				</div>
			</Card.Footer>
		</Card.Root>
	</div>

	<div class="grid grid-cols-[repeat(auto-fill,minmax(min(28rem,100%),1fr))] justify-center gap-4">
		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Adoptions</Card.Title>
				<Card.Description>Adoptions by status</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigAdoptions} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataAdoptions}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.myAnalyticsStatsResult.data?.stats.adoptions.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Adoptions stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Animals</Card.Title>
				<Card.Description>Animals by status</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigAnimals} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataAnimals}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.myAnalyticsStatsResult.data?.stats.animals.posted.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Animals stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Organizations</Card.Title>
				<Card.Description>Organizations by status</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigOrganizations} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataOrganizations}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.myAnalyticsStatsResult.data?.stats.organizations.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Organizations stats</div>
			</Card.Footer>
		</Card.Root>
	</div>
{/if}
