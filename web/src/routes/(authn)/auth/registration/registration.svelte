<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import {
		type UpdateRegistrationFlowWithPasswordMethod,
		instanceOfRegistrationFlow,
		ResponseError,
		isGenericErrorResponse,
		isBrowserLocationChangeRequired,
		FetchError,
		RequiredError
	} from '@ory/client-fetch';
	import { goto } from '$app/navigation';
	import { config } from '$lib/kratos/config';
	import { kratos } from '$lib/kratos/client';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import {
		isErrorIdSecurityCsrfViolation,
		isErrorIdSecurityIdentityMismatch,
		isErrorIdSelfServiceFlowExpired,
		isErrorIdSessionAlreadyAvailable,
		providers
	} from '$lib/kratos/helpers';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Form from '$lib/components/ui/form';
	import { browser } from '$app/environment';
	import PasswordInput from '$lib/components/password-input/password-input.svelte';

	let {
		data,
		capture = $bindable(),
		restore = $bindable()
	}: PageProps & { capture: Snapshot['capture']; restore: Snapshot['restore'] } = $props();

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data = { ...data, flow: null, csrf: data?.csrf ?? '' };
		if (browser) {
			goto(redirectUrl);
		}
		return;
	}

	export const registrationFormSchema = v.object({
		csrf_token: v.pipe(v.string(), v.minLength(1, 'csrf_token is required')),
		method: v.literal('password'),
		password: v.pipe(v.string(), v.minLength(8, 'Password must have min. 8 characters')),
		traits: v.object({
			first_name: v.string(),
			last_name: v.string(),
			email: v.pipe(v.string(), v.minLength(1, 'E-Mail is required'), v.email('E-Mail must be a valid email')),
			username: v.pipe(v.string(), v.minLength(1, 'Username is required')),
			avatar_url: v.string()
		}),
		transient_payload: v.optional(v.object({}))
	});

	type RegistrationFormSchema = v.InferInput<typeof registrationFormSchema>;

	const initialRegistrationForm: RegistrationFormSchema = {
		password: '',
		method: 'password',
		csrf_token: data.csrf ?? '',
		traits: {
			first_name: '',
			last_name: '',
			email: '',
			username: '',
			avatar_url: ''
		},
		transient_payload: {}
	};

	const supForm = superForm(initialRegistrationForm, {
		id: 'auth_registration',
		validators: valibot(registrationFormSchema),
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
			const url = data.flow?.ui.action;
			const body = form.data as UpdateRegistrationFlowWithPasswordMethod & { method: 'password' };
			if (url) {
				try {
					const successfulRegistration = await kratos.updateRegistrationFlow({
						flow: data.flow?.id ?? '',
						updateRegistrationFlowBody: body
					});
					if (successfulRegistration.continue_with) {
						for (const item of successfulRegistration.continue_with) {
							switch (item.action) {
								case 'show_verification_ui':
									if (item?.flow?.id) {
										goto(item?.flow?.url as string);
									}
									break;
							}
						}
					}
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}
					if (error instanceof ResponseError) {
						const err = await error.response.json();
						switch (error.response.status) {
							case 400: {
								if (instanceOfRegistrationFlow(err)) {
									data = { ...data, flow: err, csrf: data.csrf ?? '' };
									const nodes = err.ui.nodes ?? [];
									const fieldErrors: ValidationErrors<RegistrationFormSchema> = {};
									for (const node of nodes) {
										const errMsgs: string[] = [];
										if (node.attributes.node_type === 'input') {
											for (const msg of node?.messages ?? []) {
												errMsgs.push(msg.text);
												const fieldName = node.attributes.name;
												set(fieldErrors, fieldName, errMsgs);
											}
										}
									}
									errors.set(fieldErrors);
								}
								break;
							}
							case 410: {
								if (isGenericErrorResponse(err)) {
									if (isErrorIdSessionAlreadyAvailable(err.error?.id)) {
										goto('/');
									} else if (isErrorIdSelfServiceFlowExpired(err.error?.id)) {
										if (browser) {
											goto(`${config.routes.registration.path}?return_to=${window.location.href}`);
										}
									} else if (isErrorIdSecurityCsrfViolation(err.error?.id)) {
										handleFlowErrAction(config.routes.login.path, err.error.message);
									} else if (isErrorIdSecurityIdentityMismatch(err.error?.id)) {
										goto('/');
									}
								}
								break;
							}
							case 422: {
								if (isBrowserLocationChangeRequired(err)) {
									window.location.href = err.redirect_browser_to || '/';
								}
								break;
							}
							case 500:
								console.error('unexpected server error');
								break;
							default:
								break;
						}
						return;
					}
					if (error instanceof FetchError) {
						console.error('fetch error: ', error.cause);
						return;
					}
					if (error instanceof RequiredError) {
						console.error('required error: ', error.field);
						return;
					}
					console.error('unexpected error');
				}
			}
		}
	});

	const { form, enhance, errors } = supForm;
	capture = supForm.capture;
	restore = supForm.restore;
</script>

<section class="grid h-[calc(100vh-4rem)] place-content-center gap-4">
	<a href="/" class="justify-self-center">
		<img class="mr-2 h-8 w-8" src="/images/logo.svg" alt="logo" />
	</a>

	<Card.Root class="mx-auto max-w-md">
		<Card.Header>
			<Card.Title class="text-center text-2xl">Register</Card.Title>
			<Card.Description>Create a new account to expore your fluffy friends</Card.Description>
		</Card.Header>

		<Card.Content>
			<div class="grid gap-4">
				<form method="POST" use:enhance class="grid gap-4">
					{#each data?.flow?.ui?.messages ?? [] as msg}
						<Alert.Root variant={msg.type === '11184809' ? 'info' : msg.type} icon>
							<Alert.Title>{msg.type === 'error' ? 'Unable to sign up' : ''}</Alert.Title>
							<Alert.Description>{msg.text}</Alert.Description>
						</Alert.Root>
					{/each}

					<div class="grid grid-cols-2 gap-4">
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
						<Form.Field form={supForm} name="password">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Password</Form.Label>
									<PasswordInput {...props} bind:value={$form.password} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<Form.Button class="w-full font-bold">Create an account</Form.Button>
				</form>

				<section>
					<div class="inline-flex w-full items-center justify-center">
						<hr class="my-8 h-px w-full border-0 bg-gray-200 dark:bg-gray-700" />
						<span
							class="bg-card dark:bg-card absolute left-1/2 -translate-x-1/2 px-3 font-medium text-gray-900 dark:text-white"
						>
							or signup with
						</span>
					</div>

					<div class="align-center flex flex-wrap justify-between">
						{#each providers as provider}
							<form
								action={data.flow?.ui.action}
								method="post"
								encType="application/x-www-form-urlencoded"
								data-provider={provider.name}
							>
								<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
								<input type="hidden" name="provider" value={provider.name} readonly required />

								<Tooltip.Provider>
									<Tooltip.Root delayDuration={100}>
										<Tooltip.Trigger type="submit" class="hover:bg-primary/10 rounded">
											<img
												class="h-12 w-12 object-cover"
												src="/images/providers/{provider.name}.svg"
												alt="signup with {provider.label}"
											/>
										</Tooltip.Trigger>
										<Tooltip.Content>
											<span>{provider.label}</span>
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</form>
						{/each}
					</div>
				</section>
			</div>
			<div class="mt-4 text-center text-sm">
				Already have an account?
				<a href={config.routes.login.path} class="underline">Log in</a>
			</div>
		</Card.Content>
	</Card.Root>
</section>
