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
	import PasswordInput from '$lib/components/password-input/password-input.svelte';
	import {
		CreateIdentityBodyState,
		IdentityState,
		VerifiableIdentityAddressVia,
		type components
	} from '$lib/gen/fluffly_openapi';
	import type { CustomTraits } from '$lib/kratos/service';

	let { data }: PageProps = $props();

	let createdIdentity: components['schemas']['Identity'] | null = $state(null);

	const createIdentitySchema = v.object({
		schemaId: v.pipe(v.string(), v.minLength(1, 'Schema id is required')),
		state: v.enum(CreateIdentityBodyState),
		verified: v.boolean(),
		traits: v.object({
			first_name: v.string(),
			last_name: v.string(),
			email: v.pipe(v.string(), v.minLength(1, 'E-Mail is required'), v.email('E-Mail must be a valid email')),
			avatar_url: v.string()
		}),
		password: v.union([v.literal(''), v.pipe(v.string(), v.minLength(8, 'Password must have min. 8 characters'))]),
		hashedPassword: v.string()
	});

	type CreateIdentitySchema = v.InferInput<typeof createIdentitySchema>;

	const initialIdentitySchema: CreateIdentitySchema = {
		schemaId: 'customer',
		state: CreateIdentityBodyState.active,
		verified: false,
		traits: {
			first_name: '',
			last_name: '',
			email: '',
			avatar_url: ''
		},
		password: '',
		hashedPassword: ''
	};

	const supForm = superForm(initialIdentitySchema, {
		id: 'create_identity',
		validators: valibot(createIdentitySchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		resetForm: true,
		async onUpdate({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}
			try {
				const body: components['schemas']['CreateIdentityBody'] = {
					schema_id: form.data.schemaId,
					state: form.data.state,
					traits: form.data.traits as unknown as Record<string, never>
					// metadataAdmin?: any | null;
					// metadataPublic?: any | null;
					// recoveryAddresses?: Array<RecoveryIdentityAddress>;
				};
				if (form.data.verified) {
					body.verifiable_addresses = [
						{
							verified: true,
							via: VerifiableIdentityAddressVia.email,
							value: form.data.traits.email,
							status: 'completed'
						}
					];
				}
				if (form.data.password || form.data.hashedPassword) {
					body.credentials ??= {};
					body.credentials.password = {
						config: {
							...(form.data.password && { password: form.data.password }),
							...(form.data.hashedPassword && { hashedPassword: form.data.hashedPassword })
						}
					};
				}
				const identityRes = await fluffly.POST('/identities', {
					body
				});
				if (identityRes.data) {
					createdIdentity = identityRes.data;
					toast.success('identity created successfuly');
				}
				if (identityRes.error) {
					toast.error('failed to create identity');
				}
			} catch (error) {
				console.log('err', error);
			}
		}
	});

	const { form, enhance } = supForm;
</script>

<Card.Root class="max-w-6xl">
	<Card.Header>
		<Card.Title>Create identity</Card.Title>
		<Card.Description>Create new identity</Card.Description>
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
									<Select.Root type="single" bind:value={$form.schemaId}>
										<Select.Trigger {...props}>
											{$form.schemaId ? $form.schemaId : 'Select a schema id'}
										</Select.Trigger>
										<Select.Content>
											{#each data.schemasResult?.data ?? [] as schema (schema.id)}
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
									<Select.Root type="single" bind:value={$form.state}>
										<Select.Trigger {...props}>
											{$form.state ? $form.state : 'Select a state'}
										</Select.Trigger>
										<Select.Content>
											{#each Object.values(IdentityState) as x (x)}
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
					<div class="grid items-center gap-2">
						<Form.Field form={supForm} name="verified">
							<Form.Control>
								{#snippet children({ props })}
									<div class="flex items-center gap-2">
										<Checkbox {...props} bind:checked={$form.verified} />
										<div class="space-y-1 leading-none">
											<Form.Label>Verified email</Form.Label>
											<Form.Description />
										</div>
									</div>
								{/snippet}
							</Form.Control>
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
							<Form.Description>
								Set this only if you want to create user with the password. It's advised that you create the user and
								and send them recovery link or they can start recovery flow themselves.
							</Form.Description>
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
				<Form.Button class="mt-8 font-bold">Create identity</Form.Button>
			</form>
		</div>
	</Card.Content>
</Card.Root>

{#if createdIdentity}
	<Card.Root class="max-w-6xl">
		<Card.Header>
			<Card.Title>View created identity</Card.Title>
		</Card.Header>

		<Card.Content>
			<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">ID</span>
					<a class="font-medium text-sky-500 underline" href="/dashboard/identities/{createdIdentity.id}">
						{createdIdentity.id}
					</a>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">E-Mail</span>
					<span class="font-medium">{(createdIdentity.traits as CustomTraits)['email'] ?? ''}</span>
				</div>
				{#if (createdIdentity.traits as CustomTraits)['first_name']}
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">First name</span>
						<span class="font-medium">{(createdIdentity.traits as CustomTraits)['first_name'] ?? ''}</span>
					</div>
				{/if}
				{#if (createdIdentity.traits as CustomTraits)['last_name']}
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Last name</span>
						<span class="font-medium">{(createdIdentity.traits as CustomTraits)['last_name'] ?? ''}</span>
					</div>
				{/if}
				{#if (createdIdentity.traits as CustomTraits)['avatar_url']}
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Avatar URL</span>
						<span class="font-medium">{(createdIdentity.traits as CustomTraits)['avatar_url'] ?? ''}</span>
					</div>
				{/if}
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Schema ID</span>
					<span class="font-medium">{createdIdentity.schema_id}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">State</span>
					<span class="font-medium">{createdIdentity.state}</span>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
{/if}
