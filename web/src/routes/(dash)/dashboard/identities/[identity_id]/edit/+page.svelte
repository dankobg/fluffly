<script lang="ts">
	import type { PageProps } from './$types';
	import { superForm } from 'sveltekit-superforms/client';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import * as Select from '$lib/components/ui/select/index.js';
	import { fluffly } from '$lib/fluffly/client';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { invalidate } from '$app/navigation';
	import PasswordInput from '$lib/components/password-input/password-input.svelte';
	import { IdentityState, JsonPatchOp, UpdateIdentityBodyState, type components } from '$lib/gen/fluffly_openapi';
	import type { CustomTraits } from '$lib/kratos/service';

	let { data, params }: PageProps = $props();

	let updateIdentity: components['schemas']['Identity'] | null = $state(null);

	const updateIdentitySchema = v.object({
		schemaId: v.pipe(v.string(), v.minLength(1, 'Schema id is required')),
		// state: v.union([v.literal('active'), v.literal('inactive')]),
		state: v.enum(UpdateIdentityBodyState),
		traits: v.object({
			first_name: v.string(),
			last_name: v.string(),
			email: v.pipe(v.string(), v.minLength(1, 'E-Mail is required'), v.email('E-Mail must be a valid email')),
			username: v.pipe(v.string(), v.minLength(1, 'Username is required')),
			avatar_url: v.string()
		}),
		password: v.union([v.literal(''), v.pipe(v.string(), v.minLength(8, 'Password must have min. 8 characters'))]),
		hashedPassword: v.string()
	});

	type UpdateIdentitySchema = v.InferInput<typeof updateIdentitySchema>;

	const initialIdentitySchema: UpdateIdentitySchema = {
		schemaId: data?.identity?.schema_id ?? '',
		state: (data?.identity?.state as unknown as UpdateIdentityBodyState) ?? UpdateIdentityBodyState.active,
		traits: {
			first_name: (data?.identity?.traits as CustomTraits)?.['first_name'] ?? '',
			last_name: (data?.identity?.traits as CustomTraits)?.['last_name'] ?? '',
			email: (data?.identity?.traits as CustomTraits)?.['email'] ?? '',
			username: (data?.identity?.traits as CustomTraits)?.['username'] ?? '',
			avatar_url: (data?.identity?.traits as CustomTraits)?.['avatar_url'] ?? ''
		},
		password: '',
		hashedPassword: ''
	};

	const supForm = superForm(initialIdentitySchema, {
		id: 'update_identity_form',
		validators: valibot(updateIdentitySchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		resetForm: false,
		async onUpdated({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}
			if (!data.identity) {
				return;
			}
			try {
				const body: components['schemas']['UpdateIdentityBody'] = {
					schema_id: form.data.schemaId,
					state: form.data.state,
					traits: form.data.traits as any
					// metadataAdmin?: any | null;
					// metadataPublic?: any | null;
				};
				if (form.data.password || form.data.hashedPassword) {
					body.credentials ??= {};
					body.credentials.password = {
						config: {
							...(form.data.password && { password: form.data.password }),
							...(form.data.hashedPassword && { hashedPassword: form.data.hashedPassword })
						}
					};
				}
				const identityRes = await fluffly.PUT('/identities/{id}', {
					body,
					params: {
						path: { id: data.identity.id }
					}
				});
				if (identityRes.data) {
					updateIdentity = identityRes.data;
					toast.success('identity updated successfuly');
				}
				if (identityRes.error) {
					toast.error('failed to update identity');
				}
			} catch (error) {
				console.log('err', error);
			}
		}
	});

	const { form, enhance, errors } = supForm;

	async function onSetUnverifiedAddressClick() {
		if (!data.identity) {
			return;
		}
		const index = (data.identity.verifiable_addresses ?? []).findIndex(
			x => x.value === (data?.identity?.traits as CustomTraits)?.['email']
		);
		if (index === -1) {
			return;
		}
		const updatedAt = new Date().toISOString();
		const jsonPatch: components['schemas']['JsonPatch'][] = [
			{
				op: JsonPatchOp.replace,
				path: `/verifiable_addresses/${index}/verified`,
				value: false
			},
			{
				op: JsonPatchOp.replace,
				path: `/verifiable_addresses/${index}/status`,
				value: 'pending'
			},
			{
				op: JsonPatchOp.replace,
				path: `/verifiable_addresses/${index}/updated_at`,
				value: updatedAt
			}
		];
		try {
			await fluffly.PATCH('/identities/{id}', {
				body: jsonPatch,
				params: {
					path: { id: data.identity.id }
				}
			});
			toast.success('email set to unverified');
			invalidate(`data:identity-${data.identity.id}`);
		} catch (error) {
			console.log('err', error);
			toast.error('set email to unverified failed');
		}
	}

	async function onSetVerifiedAddressClick() {
		if (!data.identity) {
			return;
		}
		const index = (data.identity.verifiable_addresses ?? []).findIndex(
			x => x.value === (data?.identity?.traits as CustomTraits)?.['email']
		);
		if (index === -1) {
			return;
		}
		const updatedAt = new Date().toISOString();
		const jsonPatch: components['schemas']['JsonPatch'][] = [
			{
				op: JsonPatchOp.replace,
				path: `/verifiable_addresses/${index}/verified`,
				value: true
			},
			{
				op: JsonPatchOp.replace,
				path: `/verifiable_addresses/${index}/status`,
				value: 'completed'
			},
			{
				op: JsonPatchOp.replace,
				path: `/verifiable_addresses/${index}/updated_at`,
				value: updatedAt
			}
		];
		try {
			await fluffly.PATCH('/identities/{id}', {
				body: jsonPatch,
				params: {
					path: { id: data.identity.id }
				}
			});
			toast.success('email set to verified');
			invalidate(`data:identity-${data.identity.id}`);
		} catch (error) {
			console.log('err', error);
			toast.error('set email to verified failed');
		}
	}
</script>

<Card.Root class="max-w-6xl">
	<Card.Header>
		<Card.Title>Update identity</Card.Title>
		<Card.Description
			>Update identity details for {(data?.identity?.traits as CustomTraits)?.['email']}</Card.Description
		>
	</Card.Header>

	<Card.Content>
		<div class="grid gap-4">
			<form method="POST" use:enhance>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div class="grid gap-2">
						<Form.Field form={supForm} name="schemaId">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Schema ID</Form.Label>
									<Select.Root type="single" bind:value={$form.schemaId} name="schemaId">
										<Select.Trigger {...props}>
											{$form.schemaId ? $form.schemaId : 'Select a schema id'}
										</Select.Trigger>
										<Select.Content>
											{#each data.schemas ?? [] as schema}
												<Select.Item value={schema?.id ?? ''} label={schema?.id ?? ''} />
											{/each}
										</Select.Content>
									</Select.Root>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="state">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>State</Form.Label>
									<Select.Root type="single" bind:value={$form.state} name="state">
										<Select.Trigger {...props}>
											{$form.state ? $form.state : 'Select a state'}
										</Select.Trigger>
										<Select.Content>
											{#each Object.values(IdentityState) as x}
												<Select.Item value={x} label={x} />
											{/each}
										</Select.Content>
									</Select.Root>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.first_name">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>First name</Form.Label>
									<Input {...props} bind:value={$form.traits.first_name} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.last_name">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Last name</Form.Label>
									<Input {...props} bind:value={$form.traits.last_name} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.username">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Username</Form.Label>
									<Input {...props} bind:value={$form.traits.username} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.email">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>E-Mail</Form.Label>
									<Input type="email" {...props} bind:value={$form.traits.email} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.avatar_url">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Avatar URL</Form.Label>
									<Input {...props} bind:value={$form.traits.avatar_url} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="password">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Password</Form.Label>
									<PasswordInput {...props} bind:value={$form.password} />
								{/snippet}
							</Form.Control>
							<Form.Description>Set this only if you want to update the user password.</Form.Description>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="hashedPassword">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Hashed password</Form.Label>
									<PasswordInput {...props} bind:value={$form.hashedPassword} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
				</div>
				<Form.Button class="mt-8 font-bold">Update identity</Form.Button>
			</form>
		</div>
	</Card.Content>
</Card.Root>

<Card.Root class="max-w-6xl">
	<Card.Header>
		<Card.Title>Update verified address</Card.Title>
		<Card.Description>
			Update verified address state for: {(data?.identity?.traits as CustomTraits)?.['email']}
		</Card.Description>
	</Card.Header>

	<Card.Content>
		{@const verified =
			data.identity?.verifiable_addresses?.find(x => x.value === (data.identity?.traits as CustomTraits)?.['email'])
				?.verified ?? false}
		{#if verified}
			<Button onclick={onSetUnverifiedAddressClick}>Set unverified</Button>
		{:else}
			<Button onclick={onSetVerifiedAddressClick}>Set verified</Button>
		{/if}
	</Card.Content>
</Card.Root>
