<script lang="ts">
</script>

SOME DASHBOARD PLACEHOLDER

<!-- <script lang="ts">
	import * as Chart from '$lib/components/ui/chart/index';
	import * as Card from '$lib/components/ui/card/index';
	import Button from '$lib/components/ui/button/button.svelte';
	import { PieChart, Text } from 'layerchart';
	import type { PageProps } from './$types';
	
	let { data, params }: PageProps = $props();





	const chartConfig = {
		win: { label: 'Wins', color: 'var(--chart-2)' },
		loss: { label: 'Losses', color: 'var(--chart-5)' },
		draw: { label: 'Draws', color: 'var(--chart-3)' },
		interrupted: { label: 'Interrupted', color: 'var(--chart-4)' }
	} satisfies Chart.ChartConfig;

	const facts = $derived.by(() => {
		const timeControls = Object.entries(data.gameStats ?? {}).filter(([k]) => k !== 'all');
		let mostPlayed = null;
		let highestWinRate = null;
		let bestPerformance = null;
		let worstPerformance = null;
		for (const [timeControl, stats] of timeControls) {
			const { win, draw, loss, total } = stats;
			const winRate = total ? win / total : 0;
			const performanceScore = win * 1 + draw * 0.5 - loss * 1;
			if (!mostPlayed || stats.total > mostPlayed.total) {
				mostPlayed = { timeControl, total: stats.total };
			}
			if (!highestWinRate || winRate > highestWinRate.rate) {
				highestWinRate = { timeControl, rate: winRate };
			}
			if (!bestPerformance || performanceScore > bestPerformance.score) {
				bestPerformance = { timeControl, score: performanceScore };
			}
			if (!worstPerformance || performanceScore < worstPerformance.score) {
				worstPerformance = { timeControl, score: performanceScore };
			}
		}
		return {
			mostPlayed,
			highestWinRate,
			bestPerformance,
			worstPerformance
		};
	});

	const timeControlIcons = {
		hyperbullet: hyperBulletIcon,
		bullet: bulletIcon,
		blitz: blitzIcon,
		rapid: rapidIcon,
		classical: classicalIcon
	};
</script>

{#if data.gameStats?.all.total === 0}
	<div class="mb-8 mt-8 flex flex-col items-center gap-4">
		<p class="text-center text-xl font-semibold">You haven't played any games yet</p>
		<Button href="/">Play a quick game from the pool</Button>
	</div>
{:else}
	<div
		class="*:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card *:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 mb-4 grid grid-cols-1 gap-4 *:data-[slot=card]:bg-gradient-to-t"
	>
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Most played</Card.Description>
				<Card.Title class="@[250px]/card:text-3xl text-2xl font-semibold capitalize tabular-nums">
					{facts.mostPlayed?.timeControl}
				</Card.Title>
				<Card.Action>
					{@render timeControlIcons[facts.mostPlayed?.timeControl as keyof typeof timeControlIcons]('h-8 w-8')}
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 text-2xl font-medium">
					Games played: {facts.mostPlayed?.total}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Highest win rate</Card.Description>
				<Card.Title class="@[250px]/card:text-3xl text-2xl font-semibold capitalize tabular-nums"
					>{facts.highestWinRate?.timeControl}</Card.Title
				>
				<Card.Action>
					{@render timeControlIcons[facts.highestWinRate?.timeControl as keyof typeof timeControlIcons]('h-8 w-8')}
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 text-2xl font-medium">
					Win rate: {((facts.highestWinRate?.rate ?? 0) * 100).toFixed(1)}%
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Best performance</Card.Description>
				<Card.Title class="@[250px]/card:text-3xl text-2xl font-semibold capitalize tabular-nums">
					{facts.bestPerformance?.timeControl}
				</Card.Title>
				<Card.Action>
					{@render timeControlIcons[facts.bestPerformance?.timeControl as keyof typeof timeControlIcons]('h-8 w-8')}
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 text-2xl font-medium">
					Score: {facts.bestPerformance?.score}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Worst performance</Card.Description>
				<Card.Title class="@[250px]/card:text-3xl text-2xl font-semibold capitalize tabular-nums">
					{facts.worstPerformance?.timeControl}
				</Card.Title>
				<Card.Action>
					{@render timeControlIcons[facts.worstPerformance?.timeControl as keyof typeof timeControlIcons]('h-8 w-8')}
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 text-2xl font-medium">
					Score: {facts.worstPerformance?.score}
				</div>
			</Card.Footer>
		</Card.Root>
	</div>
{/if}

<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
	{#each Object.entries(chartDataByCategory ?? {}).sort((a, b) => a[1].weight - b[1].weight) as [cat, val]}
		<Card.Root class="flex flex-col">
			<Card.Header>
				<Card.Title class="text-center"><span class="capitalize">{cat}</span> results</Card.Title>
				<Card.Description class="text-center">
					{cat === 'all' ? `Game results across all time controls` : `Game results for ${cat} time control`}
				</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={val.data}
						key="result"
						value="count"
						c="color"
						innerRadius={60}
						padding={28}
						legend
						props={{ pie: { motion: 'tween' } }}
					>
						{#snippet tooltip()}
							<Chart.Tooltip hideLabel />
						{/snippet}
						{#snippet aboveMarks()}
							<Text
								value={val.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total games"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<Button class="cursor-pointer">More details</Button>
			</Card.Footer>
		</Card.Root>
	{/each}
</div> -->
