<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';
	import MessageSquareText from '@lucide/svelte/icons/message-square-text';
	import Bot from '@lucide/svelte/icons/bot';
	import RotateCcw from '@lucide/svelte/icons/rotate-ccw';
	import {
		getPromptConfig,
		setPromptTemplate,
		setSystemPrompt,
		savePromptConfig,
		resetPrompt
	} from './settings.svelte.ts';

	const promptConfig = $derived(getPromptConfig());
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-start justify-between">
		<div>
			<h2 class="text-lg font-semibold">Prompt</h2>
			<p class="text-sm text-muted-foreground">Customize the translation prompt template</p>
		</div>
		<Button variant="outline" size="sm" onclick={resetPrompt} class="gap-1.5">
			<RotateCcw class="size-3.5" />
			Reset
		</Button>
	</div>

	<!-- System Prompt Editor -->
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<Bot class="size-4 text-muted-foreground" />
			<Label class="text-sm font-medium">System Prompt</Label>
		</div>
		<Textarea
			value={promptConfig.systemPrompt}
			oninput={(e) => setSystemPrompt(e.currentTarget.value)}
			onblur={savePromptConfig}
			placeholder="Enter the system prompt for the AI..."
			class="min-h-[100px] resize-none border-border bg-background font-mono text-sm"
		/>
		<p class="text-xs text-muted-foreground">
			The system prompt sets the AI's behavior and role. Used with terminal agents like Claude Code.
		</p>
	</div>

	<!-- Prompt Template Editor -->
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<MessageSquareText class="size-4 text-muted-foreground" />
			<Label class="text-sm font-medium">Prompt Template</Label>
		</div>
		<Textarea
			value={promptConfig.template}
			oninput={(e) => setPromptTemplate(e.currentTarget.value)}
			onblur={savePromptConfig}
			placeholder="Enter your translation prompt..."
			class="min-h-[200px] resize-none border-border bg-background font-mono text-sm"
		/>
	</div>

	<!-- Variables Help -->
	<div class="rounded-lg border border-border bg-muted/30 p-4">
		<Label class="text-sm font-medium">Available Variables</Label>
		<div class="mt-2 flex flex-col gap-1.5 text-sm text-muted-foreground">
			<div class="flex items-center gap-2">
				<code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{'{{source_lang}}'}</code>
				<span>Source language name</span>
			</div>
			<div class="flex items-center gap-2">
				<code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{'{{target_lang}}'}</code>
				<span>Target language name</span>
			</div>
			<div class="flex items-center gap-2">
				<code class="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{'{{text}}'}</code>
				<span>Text to translate</span>
			</div>
		</div>
	</div>
</div>
