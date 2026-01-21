<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';
	import X from '@lucide/svelte/icons/x';
	import Copy from '@lucide/svelte/icons/copy';

	interface Props {
		label: string;
		value: string;
		placeholder?: string;
		readonly?: boolean;
		loading?: boolean;
		onClear?: () => void;
		onCopy?: () => void;
	}

	let {
		label,
		value = $bindable(),
		placeholder = '',
		readonly = false,
		loading = false,
		onClear,
		onCopy
	}: Props = $props();

	function handleCopy() {
		if (value && onCopy) {
			navigator.clipboard.writeText(value);
			onCopy();
		} else if (value) {
			navigator.clipboard.writeText(value);
		}
	}
</script>

<div class="flex flex-col overflow-hidden rounded-md border border-border bg-surface transition-colors duration-150 focus-within:border-accent">
	<div class="flex min-h-11 items-center justify-between border-b border-border px-3 py-2">
		<span class="text-xs font-medium uppercase tracking-wide text-text-muted">{label}</span>
		{#if readonly}
			<Button
				variant="ghost"
				size="icon-sm"
				class="h-7 w-7 text-text-muted hover:bg-success/10 hover:text-success {!value || loading ? 'opacity-0 pointer-events-none' : ''}"
				onclick={handleCopy}
				disabled={!value || loading}
			>
				<Copy class="size-3.5" />
			</Button>
		{:else}
			<Button
				variant="ghost"
				size="icon-sm"
				class="h-7 w-7 text-text-muted hover:bg-red-500/10 hover:text-red-500 {!value ? 'opacity-0 pointer-events-none' : ''}"
				onclick={onClear}
				disabled={!value}
			>
				<X class="size-3.5" />
			</Button>
		{/if}
	</div>

	{#if readonly}
		<div
			class="flex-1 overflow-y-auto whitespace-pre-wrap break-words bg-transparent p-3 font-mono text-base leading-relaxed text-text"
			class:flex={loading}
			class:items-center={loading}
			class:justify-center={loading}
		>
			{#if loading}
				<div class="flex gap-1.5">
					<span class="h-2 w-2 animate-bounce rounded-sm bg-accent [animation-delay:-0.32s]"></span>
					<span class="h-2 w-2 animate-bounce rounded-sm bg-accent [animation-delay:-0.16s]"></span>
					<span class="h-2 w-2 animate-bounce rounded-sm bg-accent"></span>
				</div>
			{:else if value}
				{value}
			{:else}
				<span class="text-text-muted">{placeholder}</span>
			{/if}
		</div>
	{:else}
		<Textarea
			bind:value
			class="min-h-36 flex-1 resize-none rounded-none border-none bg-transparent p-3 font-mono text-base leading-relaxed shadow-none focus-visible:ring-0"
			{placeholder}
			spellcheck={false}
		/>
	{/if}

	<div class="min-h-8 border-t border-border px-3 py-1.5">
		{#if value && !loading}
			<span class="font-mono text-xs text-text-muted">{value.length}{readonly ? ' characters' : ' / 5000'}</span>
		{/if}
	</div>
</div>
