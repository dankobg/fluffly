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

	const chartConfigUsers = {
		active: { label: 'Active', color: 'var(--color-green-600)' },
		inactive: { label: 'Inactive', color: 'var(--color-red-600)' }
	} satisfies Chart.ChartConfig;

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

	let chartDataUsers = $derived.by(() => {
		if (!data?.analyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigUsers.active.label,
				value: data.analyticsStatsResult.data.stats.users.by_state.active,
				color: chartConfigUsers.active.color
			},
			{
				label: chartConfigUsers.inactive.label,
				value: data.analyticsStatsResult.data.stats.users.by_state.inactive,
				color: chartConfigUsers.inactive.color
			}
		];
	});

	let chartDataAdoptions = $derived.by(() => {
		if (!data?.analyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigAdoptions.approved.label,
				value: data.analyticsStatsResult.data.stats.adoptions.by_status.approved,
				color: chartConfigAdoptions.approved.color
			},
			{
				label: chartConfigAdoptions.pending.label,
				value: data.analyticsStatsResult.data.stats.adoptions.by_status.pending,
				color: chartConfigAdoptions.pending.color
			},
			{
				label: chartConfigAdoptions.rejected.label,
				value: data.analyticsStatsResult.data.stats.adoptions.by_status.rejected,
				color: chartConfigAdoptions.rejected.color
			}
		];
	});

	let chartDataOrganizations = $derived.by(() => {
		if (!data?.analyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigOrganizations.approved.label,
				value: data.analyticsStatsResult.data.stats.organizations.by_status.approved,
				color: chartConfigOrganizations.approved.color
			},
			{
				label: chartConfigOrganizations.pending.label,
				value: data.analyticsStatsResult.data.stats.organizations.by_status.pending,
				color: chartConfigOrganizations.pending.color
			},
			{
				label: chartConfigOrganizations.rejected.label,
				value: data.analyticsStatsResult.data.stats.organizations.by_status.rejected,
				color: chartConfigOrganizations.rejected.color
			}
		];
	});

	let chartDataAnimals = $derived.by(() => {
		if (!data?.analyticsStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigAnimals.adoptable.label,
				value: data.analyticsStatsResult.data.stats.animals.by_status.adoptable,
				color: chartConfigAnimals.adoptable.color
			},
			{
				label: chartConfigAnimals.pending.label,
				value: data.analyticsStatsResult.data.stats.animals.by_status.pending,
				color: chartConfigAnimals.pending.color
			},
			{
				label: chartConfigAnimals.adopted.label,
				value: data.analyticsStatsResult.data.stats.animals.by_status.adopted,
				color: chartConfigAnimals.adopted.color
			},
			{
				label: chartConfigAnimals.reserved.label,
				value: data.analyticsStatsResult.data.stats.animals.by_status.reserved,
				color: chartConfigAnimals.reserved.color
			},
			{
				label: chartConfigAnimals.rejected.label,
				value: data.analyticsStatsResult.data.stats.animals.by_status.rejected,
				color: chartConfigAnimals.rejected.color
			}
		];
	});

	let usersActivePercent = $derived.by(() => {
		if (!data?.analyticsStatsResult?.data) {
			return 0;
		}
		const percent =
			(data?.analyticsStatsResult?.data.stats.users.by_state.active /
				data?.analyticsStatsResult?.data.stats.users.total) *
			100;
		return Number(percent.toFixed(2));
	});

	let usersActivityText = $derived.by(() => {
		if (usersActivePercent === 0) return 'No active users';
		if (usersActivePercent < 40) return 'Low user activity';
		if (usersActivePercent < 70) return 'Moderate user activity';
		if (usersActivePercent < 90) return 'Healthy user activity';
		return 'Strong user activity';
	});
</script>

{#if data?.analyticsStatsResult?.data}
	<div class="grid grid-cols-1 gap-4 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Adoptions</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.analyticsStatsResult?.data?.stats?.adoptions?.by_status?.approved}
				</Card.Title>
				<Card.Action>
					<IconHandHeart />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending adoptions {data?.analyticsStatsResult?.data?.stats?.adoptions?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.analyticsStatsResult?.data?.stats?.organizations?.by_status.pending > 0
						? 'Required admin attention'
						: 'No pending adoptions'}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Users</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.analyticsStatsResult?.data?.stats?.users.total}
				</Card.Title>
				<Card.Action>
					<IconUsers />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					{usersActivePercent}% activity rate
				</div>
				<div class="text-muted-foreground">
					{usersActivityText}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Animals</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.analyticsStatsResult?.data?.stats?.animals?.total}
				</Card.Title>
				<Card.Action>
					<IconCat />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.analyticsStatsResult?.data?.stats?.animals?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.analyticsStatsResult?.data?.stats?.animals?.by_status.pending > 0
						? 'Required admin attention'
						: 'No pending animals'}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Organizations</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.analyticsStatsResult?.data?.stats?.organizations?.total}
				</Card.Title>
				<Card.Action>
					<IconBuilding2 />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.analyticsStatsResult?.data?.stats?.organizations?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.analyticsStatsResult?.data?.stats?.organizations?.by_status.pending > 0
						? 'Required admin attention'
						: 'No pending organizations'}
				</div>
			</Card.Footer>
		</Card.Root>
	</div>

	<div class="grid grid-cols-[repeat(auto-fill,minmax(min(35rem,100%),1fr))] justify-center gap-4">
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
								value={data.analyticsStatsResult.data?.stats.adoptions.total}
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
				<Card.Title>Users</Card.Title>
				<Card.Description>Users by state</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigUsers} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataUsers}
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
								value={data.analyticsStatsResult.data?.stats.users.total}
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
				<div class="leading-none text-muted-foreground">Users stats</div>
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
								value={data.analyticsStatsResult.data?.stats.animals.total}
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
								value={data.analyticsStatsResult.data?.stats.organizations.total}
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
