<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { CourierMessageTemplateType } from '$lib/gen/fluffly_openapi';
	import { templateTypeIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = templateTypeIcons.get(value as CourierMessageTemplateType);
	let color = $derived.by(() => {
		switch (value as CourierMessageTemplateType) {
			case CourierMessageTemplateType.recovery_valid:
			case CourierMessageTemplateType.recovery_code_valid:
			case CourierMessageTemplateType.verification_valid:
			case CourierMessageTemplateType.verification_code_valid:
			case CourierMessageTemplateType.login_code_valid:
			case CourierMessageTemplateType.registration_code_valid:
				return 'text-green-400';
			case CourierMessageTemplateType.verification_code_invalid:
			case CourierMessageTemplateType.verification_invalid:
			case CourierMessageTemplateType.recovery_code_invalid:
			case CourierMessageTemplateType.recovery_invalid:
				return 'text-red-400';
			case CourierMessageTemplateType.stub:
				return 'text-purple-400';
			default:
				return '';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	<Icon class="" />
	<span>{value ?? ''}</span>
</Badge>
