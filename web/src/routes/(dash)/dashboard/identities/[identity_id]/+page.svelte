<script lang="ts">
	import type { PageProps } from './$types';
	import * as Table from '$lib/components/ui/table/index';
	import { stateIcons } from '../identities-data-grid/data';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import Button from '$lib/components/ui/button/button.svelte';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { fluffly } from '$lib/fluffly/client';
	import { goto } from '$app/navigation';
	import { IdentityState } from '$lib/gen/fluffly_openapi';
	import type { CustomTraits } from '$lib/kratos/service';

	let { data, params }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let StateIcon = $derived(data.identity?.state && stateIcons.get(data.identity.state));
	let stateIconClasses = $derived.by(() => {
		switch (data.identity?.state as IdentityState) {
			case IdentityState.active:
				return 'text-green-400';
			case IdentityState.inactive:
				return 'text-red-400';
			default:
				return '';
		}
	});

	async function onConfirmDeleteIdentity() {
		if (!data.identity) {
			return;
		}
		try {
			await fluffly.DELETE('/identities/{id}', {
				params: {
					path: { id: data.identity.id }
				}
			});
			goto('/dashboard/identities');
		} catch (error) {
			console.log('err', error);
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeleteIdentityClick() {
		confirmation.openDialog({
			description: deleteIdentityDescriptionSnippet,
			onConfirm: onConfirmDeleteIdentity,
			destructive: true
		});
	}
</script>

{#snippet deleteIdentityDescriptionSnippet()}
	{@const email = (data?.identity?.traits as CustomTraits)?.['email']}
	This action cannot be undone. This will delete the identity <strong>{email}</strong> completely.
{/snippet}

{#if data.identity}
	<section class="mb-6 gap-4">
		<p class="mb-6 text-lg">Actions</p>
		<div class="flex gap-4">
			<Button href="/dashboard/identities/{data.identity.id}/sessions">View sessions</Button>
			<Button href="/dashboard/identities/{data.identity.id}/edit">Edit identity</Button>
			<Button variant="destructive" onclick={onDeleteIdentityClick}>Delete identity</Button>
		</div>
	</section>

	<h1 class="mb-6 text-2xl font-bold">Identity</h1>
	<p class="mb-6 text-lg">Details</p>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.identity.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">E-Mail</span>
			<span class="font-medium">{(data.identity.traits as CustomTraits)['email'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">First name</span>
			<span class="font-medium">{(data.identity.traits as CustomTraits)['first_name'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Last name</span>
			<span class="font-medium">{(data.identity.traits as CustomTraits)['last_name'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Avatar URL</span>
			<span class="font-medium">{(data.identity.traits as CustomTraits)['avatar_url'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Schema ID</span>
			<span class="font-medium">{data.identity.schema_id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Schema URL</span>
			<span class="font-medium">{data.identity.schema_url}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">State</span>
			<span class="flex gap-2 font-medium">{data.identity.state} <StateIcon class={stateIconClasses} /></span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">State changed time</span>
			<time class="font-medium"
				>{data.identity.state_changed_at && fmt.format(new Date(data.identity.state_changed_at))}</time
			>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{data.identity.created_at && fmt.format(new Date(data.identity.created_at))}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{data.identity.updated_at && fmt.format(new Date(data.identity.updated_at))}</time>
		</div>
	</div>

	<p class="mt-8 text-lg">Credentials</p>
	<Table.Root>
		<Table.Caption>A list of credentials</Table.Caption>
		<Table.Header>
			<Table.Row>
				<Table.Head>Type</Table.Head>
				<Table.Head>Version</Table.Head>
				<Table.Head>Config</Table.Head>
				<Table.Head>Identifiers</Table.Head>
				<Table.Head>Created time</Table.Head>
				<Table.Head>Update time</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each Object.values(data.identity.credentials ?? {}) as credential}
				<Table.Row>
					<Table.Cell class="font-medium">{credential.type}</Table.Cell>
					<Table.Cell class="font-medium">{credential.version}</Table.Cell>
					<Table.Cell class="font-medium"><pre>{JSON.stringify(credential.config, null, 2)}</pre></Table.Cell>
					<Table.Cell class="font-medium">{credential.identifiers?.join(', ')}</Table.Cell>
					<Table.Cell>{credential.created_at && fmt.format(new Date(credential.created_at))}</Table.Cell>
					<Table.Cell>{credential.updated_at && fmt.format(new Date(credential.updated_at))}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>

	{#if data.identity.recovery_addresses && data.identity.recovery_addresses.length > 0}
		<p class="mt-8 text-lg">Recovery addresses</p>
		<Table.Root>
			<Table.Caption>A list of recovery addresses</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Value</Table.Head>
					<Table.Head>Via</Table.Head>
					<Table.Head>Created time</Table.Head>
					<Table.Head>Update time</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.identity.recovery_addresses as recAddr (recAddr)}
					<Table.Row>
						<Table.Cell class="font-medium">{recAddr.id}</Table.Cell>
						<Table.Cell>{recAddr.value}</Table.Cell>
						<Table.Cell>{recAddr.via}</Table.Cell>
						<Table.Cell>{recAddr.created_at && fmt.format(new Date(recAddr.created_at))}</Table.Cell>
						<Table.Cell>{recAddr.updated_at && fmt.format(new Date(recAddr.updated_at))}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}

	{#if data.identity.verifiable_addresses && data.identity.verifiable_addresses.length > 0}
		<p class="mt-8 text-lg">Verifiable addresses:</p>
		<Table.Root>
			<Table.Caption>A list of verifiable addresses</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Value</Table.Head>
					<Table.Head>Via</Table.Head>
					<Table.Head>Status</Table.Head>
					<Table.Head>Verfiied</Table.Head>
					<Table.Head>Verified time</Table.Head>
					<Table.Head>Created time</Table.Head>
					<Table.Head>Update time</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.identity.verifiable_addresses as verAddr (verAddr)}
					<Table.Row>
						<Table.Cell class="font-medium">{verAddr.id}</Table.Cell>
						<Table.Cell>{verAddr.value}</Table.Cell>
						<Table.Cell>{verAddr.via}</Table.Cell>
						<Table.Cell>{verAddr.status}</Table.Cell>
						<Table.Cell>
							<div class="flex gap-2">
								{verAddr.verified}
								{#if verAddr.verified}
									<IconCheck class="text-green-400" />
								{:else}
									<IconX class="text-red-400" />
								{/if}
							</div>
						</Table.Cell>
						<Table.Cell>{verAddr.verified_at && fmt.format(new Date(verAddr.verified_at))}</Table.Cell>
						<Table.Cell>{verAddr.created_at && fmt.format(new Date(verAddr.created_at))}</Table.Cell>
						<Table.Cell>{verAddr.updated_at && fmt.format(new Date(verAddr.updated_at))}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
{/if}
