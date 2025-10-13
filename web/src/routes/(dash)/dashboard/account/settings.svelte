<script lang="ts">
	import type { PageProps } from './$types';
	import {
		FetchError,
		instanceOfSettingsFlow,
		isBrowserLocationChangeRequired,
		isGenericErrorResponse,
		RequiredError,
		ResponseError,
		type UpdateSettingsFlowWithProfileMethod
	} from '@ory/client-fetch';
	import { goto, invalidateAll } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Form from '$lib/components/ui/form';
	import { browser } from '$app/environment';
	import { config } from '$lib/kratos/config';
	import {
		isErrorIdSecurityCsrfViolation,
		isErrorIdSecurityIdentityMismatch,
		isErrorIdSessionInactive,
		isErrorIdSessionRefreshRequired
	} from '$lib/kratos/helpers';

	let {
		data,
		currentFlowForm = $bindable()
	}: PageProps & { currentFlowForm: 'settings' | 'password' | 'socials' | undefined } = $props();

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data.flow = null;
		if (browser) {
			goto(redirectUrl);
		}
		return;
	}

	const settingsFormSchema = v.object({
		csrf_token: v.pipe(v.string(), v.minLength(1, 'csrf_token is required')),
		method: v.literal('profile'),
		traits: v.object({
			first_name: v.string(),
			last_name: v.string(),
			email: v.pipe(v.string(), v.minLength(1, 'E-Mail is required'), v.email('E-Mail must be a valid email')),
			username: v.pipe(v.string(), v.minLength(1, 'Username is required')),
			avatar_url: v.string()
		})
	});

	type SettingsFormSchema = v.InferInput<typeof settingsFormSchema>;

	const initialSettingsForm: SettingsFormSchema = {
		method: 'profile',
		csrf_token: data.csrf ?? '',
		traits: {
			first_name: data.flow?.identity.traits['first_name'] ?? '',
			last_name: data.flow?.identity.traits['last_name'] ?? '',
			email: data.flow?.identity.traits['email'] ?? '',
			username: data.flow?.identity.traits['username'] ?? '',
			avatar_url: data.flow?.identity.traits['avatar_url'] ?? ''
		}
	};

	const supForm = superForm(initialSettingsForm, {
		id: 'account_settings',
		validators: valibot(settingsFormSchema),
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
			const body = form.data as UpdateSettingsFlowWithProfileMethod & { method: 'profile' };
			if (url) {
				try {
					currentFlowForm = 'settings';
					const settingsFlow = await kratos.updateSettingsFlow({
						flow: data.flow?.id ?? '',
						updateSettingsFlowBody: body
					});
					data = { ...data, flow: settingsFlow, csrf: data?.csrf ?? '' };
					toast.success('Account settings have been updated');
					invalidateAll();
					if (settingsFlow.continue_with) {
						for (const item of settingsFlow.continue_with) {
							switch (item.action) {
								case 'show_verification_ui':
									if (item?.flow?.id) {
										goto(item?.flow?.url as string);
									}
									break;
								// case 'redirect_browser_to':
								// 	window.location.href = item.redirect_browser_to;
								// 	break;
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
								if (instanceOfSettingsFlow(err)) {
									data = { ...data, flow: err, csrf: data.csrf ?? '' };
									const nodes = err.ui.nodes ?? [];
									const fieldErrors: ValidationErrors<SettingsFormSchema> = {};
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
							case 401:
							case 403:
							case 410: {
								if (isGenericErrorResponse(err)) {
									if (isErrorIdSessionRefreshRequired(err.error?.id)) {
										goto(`${config.routes.login.path}?refresh=true&return_to=${window.location.href}`);
									} else if (isErrorIdSecurityCsrfViolation(err.error?.id)) {
										handleFlowErrAction(config.routes.settings.path, err.error.message);
									} else if (isErrorIdSessionInactive(err.error?.id)) {
										handleFlowErrAction(
											config.routes.login.path + `?return_to=${encodeURIComponent(window.location.href)}`,
											err.error.message
										);
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
</script>

<Card.Root class="max-w-md">
	<Card.Header>
		<Card.Title>Account settings</Card.Title>
		<Card.Description>Change account settings</Card.Description>
	</Card.Header>

	<Card.Content>
		<div class="grid gap-4">
			<form method="POST" use:enhance class="grid gap-4">
				{#if currentFlowForm === 'settings'}
					{#each data?.flow?.ui?.messages ?? [] as msg}
						<Alert.Root variant={msg.type === '11184809' ? 'info' : msg.type} icon>
							<Alert.Title>{msg.type === 'error' ? 'Unable to change settings' : ''}</Alert.Title>
							<Alert.Description>{msg.text}</Alert.Description>
						</Alert.Root>
					{/each}
				{/if}

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
				<Form.Button class="w-full font-bold">Update settings</Form.Button>
			</form>
		</div>
	</Card.Content>
</Card.Root>
