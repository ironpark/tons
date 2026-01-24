<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { Label } from '$lib/components/ui/label';
	import Terminal from '@lucide/svelte/icons/terminal';
	import Server from '@lucide/svelte/icons/server';
	import Zap from '@lucide/svelte/icons/zap';
	import {
		EngineType,
		TerminalAgentType
	} from '$lib/bindings/github.com/ironpark/tons/internal/config/models';
	import {
		getEngineConfig,
		getSelectedTerminalAgent,
		getSelectedOllamaModel,
		terminalAgents,
		ollamaModels,
		setEngineType,
		setTerminalAgent,
		setOllamaModel
	} from './settings.svelte.ts';

	const engineConfig = $derived(getEngineConfig());
	const selectedTerminalAgent = $derived(getSelectedTerminalAgent());
	const selectedOllamaModel = $derived(getSelectedOllamaModel());
</script>

<div class="flex flex-col gap-6">
	<div>
		<h2 class="text-lg font-semibold">Engine</h2>
		<p class="text-sm text-muted-foreground">Configure translation engine</p>
	</div>

	<!-- Engine Selection -->
	<RadioGroup.Root
		value={engineConfig.type}
		onValueChange={(value) => setEngineType(value as EngineType)}
		class="flex flex-col gap-2"
	>
		<!-- Internal -->
		<Label
			class="flex cursor-pointer items-center gap-3 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
		>
			<RadioGroup.Item value={EngineType.EngineInternal} />
			<Zap class="size-4 text-muted-foreground" />
			<div class="flex flex-col">
				<span class="text-sm font-medium">Internal</span>
				<span class="text-xs text-muted-foreground">Built-in translation engine</span>
			</div>
		</Label>

		<!-- Terminal Agent -->
		<Label
			class="flex cursor-pointer items-center gap-3 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
		>
			<RadioGroup.Item value={EngineType.EngineTerminalAgent} />
			<Terminal class="size-4 text-muted-foreground" />
			<div class="flex flex-col">
				<span class="text-sm font-medium">Terminal Agent</span>
				<span class="text-xs text-muted-foreground">Use CLI AI assistants for translation</span>
			</div>
		</Label>

		<!-- Ollama -->
		<Label
			class="flex cursor-pointer items-center gap-3 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
		>
			<RadioGroup.Item value={EngineType.EngineOllama} />
			<Server class="size-4 text-muted-foreground" />
			<div class="flex flex-col">
				<span class="text-sm font-medium">Ollama</span>
				<span class="text-xs text-muted-foreground">Run AI models locally</span>
			</div>
		</Label>
	</RadioGroup.Root>

	<!-- Terminal Agent Options -->
	{#if engineConfig.type === EngineType.EngineTerminalAgent}
		<div class="flex flex-col gap-3 rounded-lg border border-border bg-muted/30 p-4">
			<Label class="text-sm font-medium">Select Terminal Agent</Label>
			<Select.Root
				type="single"
				value={engineConfig.terminalAgent.selected}
				onValueChange={(value) => setTerminalAgent(value as TerminalAgentType)}
			>
				<Select.Trigger class="w-full border-border bg-background hover:bg-accent/50">
					<span>{selectedTerminalAgent.label}</span>
				</Select.Trigger>
				<Select.Content>
					{#each terminalAgents as agent (agent.value)}
						<Select.Item value={agent.value} label={agent.label}>
							{agent.label}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	{/if}

	<!-- Ollama Options -->
	{#if engineConfig.type === EngineType.EngineOllama}
		<div class="flex flex-col gap-3 rounded-lg border border-border bg-muted/30 p-4">
			<Label class="text-sm font-medium">Select Model</Label>
			<Select.Root
				type="single"
				value={engineConfig.ollama.model}
				onValueChange={(value) => setOllamaModel(value)}
			>
				<Select.Trigger class="w-full border-border bg-background hover:bg-accent/50">
					<span>{selectedOllamaModel.label}</span>
				</Select.Trigger>
				<Select.Content>
					{#each ollamaModels as model (model.value)}
						<Select.Item value={model.value} label={model.label}>
							{model.label}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	{/if}
</div>
