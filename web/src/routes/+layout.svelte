<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from '$lib/components/ui/sonner/index';
	import ConfirmationDialog from '$lib/components/confirmation-dialog/confirmation-dialog.svelte';
	import '../app.css';
	import '@uppy/core/css/style.min.css';
	import '@uppy/dashboard/css/style.min.css';
	import '@uppy/image-editor/css/style.min.css';

	let { children } = $props();
</script>

<ConfirmationDialog />
<ModeWatcher />
<Toaster position="top-center" richColors />
{@render children()}
