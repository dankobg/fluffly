<script lang="ts">
	import { type UiNode } from '@ory/client-fetch';
	import type { PageProps } from './$types';
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';

	let { data, socialsAction }: PageProps & { socialsAction?: 'link' | 'unlink' } = $props();

	const filterBy = (n: UiNode, action: 'link' | 'unlink') =>
		n.group === 'oidc' && n.type === 'input' && n.attributes.node_type === 'input' && n.attributes.name === action;

	let providersToLink = $derived(data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'link')) ?? []);
	let providersToUnlink = $derived(data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'unlink')) ?? []);

	$effect(() => {
		const val = window.sessionStorage.getItem('socialsAction') as 'link' | 'unlink' | undefined;
		if (val) {
			socialsAction = val;
			toast.success(`Your account has been ${val}ed`);
		}
		return () => {
			socialsAction = undefined;
			sessionStorage.removeItem('socialsAction');
		};
	});
</script>

<svelte:window
	onbeforeunload={() => {
		if (socialsAction) {
			sessionStorage.removeItem('socialsAction');
		}
	}}
/>

{#if providersToLink.length > 0}
	<Card.Root class="max-w-md">
		<Card.Header>
			<Card.Title>Social providers</Card.Title>
			<Card.Description>Link/Unlink auth social providers</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-4">
			{#if socialsAction === 'link' || socialsAction === 'unlink'}
				{#each data?.flow?.ui?.messages ?? [] as msg}
					<Alert.Root variant={msg.type === '11184809' ? 'info' : msg.type} icon>
						<Alert.Title>{msg.type === 'error' ? `Unable to ${socialsAction} account` : ''}</Alert.Title>
						<Alert.Description>{msg.text}</Alert.Description>
					</Alert.Root>
				{/each}
			{/if}

			{#each providersToLink as provider}
				{#if provider.attributes.node_type === 'input'}
					<form
						action={data.flow?.ui.action}
						method="post"
						encType="application/x-www-form-urlencoded"
						class="w-full space-y-6"
					>
						<input type="hidden" name="link" value={provider.attributes.value} readonly required />
						<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

						<div class="flex w-full justify-start gap-4">
							<img
								class="ml-4 inline-flex h-6 w-6 object-cover"
								src="/images/providers/{provider.attributes.value}.svg"
								alt={provider.attributes.value}
							/>
							Link {provider.attributes.value} account
							<Switch
								type="submit"
								onclick={() => {
									window.sessionStorage.setItem('socialsAction', 'link');
								}}
								class="ml-auto"
							/>
						</div>
					</form>
				{/if}
			{/each}
			{#each providersToUnlink as provider}
				{#if provider.attributes.node_type === 'input'}
					<form
						action={data.flow?.ui.action}
						method="post"
						encType="application/x-www-form-urlencoded"
						class="w-full space-y-6"
					>
						<input type="hidden" name="unlink" value={provider.attributes.value} readonly required />
						<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

						<div class="flex w-full justify-between gap-4">
							<img
								class="ml-4 inline-flex h-6 w-6 object-cover"
								src="/images/providers/{provider.attributes.value}.svg"
								alt={provider.attributes.value}
							/>
							Unlink {provider.attributes.value} account
							<Switch
								type="submit"
								checked
								onclick={() => {
									window.sessionStorage.setItem('socialsAction', 'unlink');
								}}
								class="ml-auto"
							/>
						</div>
					</form>
				{/if}
			{/each}
		</Card.Content>
	</Card.Root>
{/if}
