<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index';
	import { CourierMessageType } from '$lib/gen/fluffly_openapi';
	import { typeIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = typeIcons.get(value as CourierMessageType);
	let color = $derived.by(() => {
		switch (value as CourierMessageType) {
			case CourierMessageType.email:
				return 'text-blue-400';
			case CourierMessageType.phone:
				return 'text-blue-400';
			default:
				return '';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	<Icon class="" />
	<span>{value ?? ''}</span>
</Badge>
