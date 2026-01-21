<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Settings from '@lucide/svelte/icons/settings';
	import Sun from '@lucide/svelte/icons/sun';
	import Moon from '@lucide/svelte/icons/moon';
	import Monitor from '@lucide/svelte/icons/monitor';
	import Globe from '@lucide/svelte/icons/globe';
	import Cpu from '@lucide/svelte/icons/cpu';
	import Terminal from '@lucide/svelte/icons/terminal';
	import Server from '@lucide/svelte/icons/server';
	import Zap from '@lucide/svelte/icons/zap';
	import Sliders from '@lucide/svelte/icons/sliders-horizontal';
	import MessageSquareText from '@lucide/svelte/icons/message-square-text';
	import RotateCcw from '@lucide/svelte/icons/rotate-ccw';

	// Navigation
	let activeSection = $state('general');
	const sections = [
		{ id: 'general', label: 'General', icon: Sliders },
		{ id: 'engine', label: 'Engine', icon: Cpu },
		{ id: 'prompt', label: 'Prompt', icon: MessageSquareText }
	];

	// Theme settings
	let theme = $state('system');

	// Language settings
	let language = $state('system');
	const languages = [
		{ value: 'system', label: 'System Default' },
		{ value: 'en', label: 'English' },
		{ value: 'ko', label: '한국어' },
		{ value: 'ja', label: '日本語' }
	];

	// Translation Engine settings
	let engine = $state('internal');

	// Terminal Agent options
	let terminalAgent = $state('claude-code');
	const terminalAgents = [
		{ value: 'claude-code', label: 'Claude Code' },
		{ value: 'gemini-cli', label: 'Gemini CLI' },
		{ value: 'codex', label: 'Codex' }
	];

	// Ollama options
	let ollamaModel = $state('llama3.2');
	const ollamaModels = [
		{ value: 'llama3.2', label: 'Llama 3.2' },
		{ value: 'llama3.1', label: 'Llama 3.1' },
		{ value: 'mistral', label: 'Mistral' },
		{ value: 'codellama', label: 'Code Llama' },
		{ value: 'gemma2', label: 'Gemma 2' }
	];

	const selectedLanguage = $derived(languages.find((l) => l.value === language) ?? languages[0]);
	const selectedTerminalAgent = $derived(terminalAgents.find((a) => a.value === terminalAgent) ?? terminalAgents[0]);
	const selectedOllamaModel = $derived(ollamaModels.find((m) => m.value === ollamaModel) ?? ollamaModels[0]);

	// Prompt settings
	const defaultPrompt = `Translate the following text from {{source_lang}} to {{target_lang}}.
Keep the original formatting and tone.
Only return the translated text without any explanations.

Text to translate:
{{text}}`;
	let prompt = $state(defaultPrompt);

	function resetPrompt() {
		prompt = defaultPrompt;
	}
</script>

<div class="relative flex h-screen flex-col overflow-hidden bg-background font-sans text-foreground">
	<!-- Background grid -->
	<div class="pointer-events-none absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.015)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.015)_1px,transparent_1px)] bg-[size:60px_60px]"></div>

	<!-- Header -->
	<header class="relative z-10 flex items-center gap-3 border-b border-border bg-background/80 py-2.5 pr-4 pl-20 backdrop-blur-xl [-webkit-app-region:drag]">
		<div class="[-webkit-app-region:no-drag]">
			<Button variant="ghost" size="icon" href="/" class="text-muted-foreground hover:text-foreground">
				<ArrowLeft class="size-[18px]" />
			</Button>
		</div>
		<div class="flex items-center gap-2">
			<div class="flex h-9 w-9 items-center justify-center rounded-md bg-accent text-accent-foreground">
				<Settings class="size-5" />
			</div>
			<h1 class="text-xl font-semibold tracking-tight">Settings</h1>
		</div>
	</header>

	<!-- Main Content -->
	<main class="relative z-10 flex flex-1 overflow-hidden">
		<!-- Sidebar -->
		<nav class="flex w-48 flex-col gap-1 border-r border-border bg-background/50 p-3">
			{#each sections as section (section.id)}
				<button
					onclick={() => (activeSection = section.id)}
					class="flex items-center gap-2 rounded-lg px-3 py-2 text-left text-sm font-medium transition-colors {activeSection === section.id
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'}"
				>
					<section.icon class="size-4" />
					{section.label}
				</button>
			{/each}
		</nav>

		<!-- Content -->
		<div class="flex-1 overflow-y-auto p-6">
			<div class="h-full">
				{#if activeSection === 'general'}
					<div class="flex flex-col gap-6">
						<div>
							<h2 class="text-lg font-semibold">General</h2>
							<p class="text-sm text-muted-foreground">Customize your app preferences</p>
						</div>

						<!-- Theme -->
						<div class="flex flex-col gap-3">
							<div class="flex items-center gap-2">
								<Sun class="size-4 text-muted-foreground" />
								<Label class="text-sm font-medium">Theme</Label>
							</div>
							<RadioGroup.Root bind:value={theme} class="flex gap-2">
								<Label
									class="flex flex-1 cursor-pointer flex-col items-center gap-2 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
								>
									<RadioGroup.Item value="light" class="sr-only" />
									<Sun class="size-5" />
									<span class="text-sm font-medium">Light</span>
								</Label>
								<Label
									class="flex flex-1 cursor-pointer flex-col items-center gap-2 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
								>
									<RadioGroup.Item value="dark" class="sr-only" />
									<Moon class="size-5" />
									<span class="text-sm font-medium">Dark</span>
								</Label>
								<Label
									class="flex flex-1 cursor-pointer flex-col items-center gap-2 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
								>
									<RadioGroup.Item value="system" class="sr-only" />
									<Monitor class="size-5" />
									<span class="text-sm font-medium">System</span>
								</Label>
							</RadioGroup.Root>
						</div>

						<!-- Language -->
						<div class="flex flex-col gap-3">
							<div class="flex items-center gap-2">
								<Globe class="size-4 text-muted-foreground" />
								<Label class="text-sm font-medium">Language</Label>
							</div>
							<Select.Root type="single" bind:value={language}>
								<Select.Trigger class="w-full border-border bg-background hover:bg-accent/50">
									<span>{selectedLanguage.label}</span>
								</Select.Trigger>
								<Select.Content>
									{#each languages as lang (lang.value)}
										<Select.Item value={lang.value} label={lang.label}>
											{lang.label}
										</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
						</div>
					</div>
				{:else if activeSection === 'engine'}
					<div class="flex flex-col gap-6">
						<div>
							<h2 class="text-lg font-semibold">Engine</h2>
							<p class="text-sm text-muted-foreground">Configure translation engine</p>
						</div>

						<!-- Engine Selection -->
						<RadioGroup.Root bind:value={engine} class="flex flex-col gap-2">
							<!-- Internal -->
							<Label
								class="flex cursor-pointer items-center gap-3 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
							>
								<RadioGroup.Item value="internal" />
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
								<RadioGroup.Item value="terminal-agent" />
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
								<RadioGroup.Item value="ollama" />
								<Server class="size-4 text-muted-foreground" />
								<div class="flex flex-col">
									<span class="text-sm font-medium">Ollama</span>
									<span class="text-xs text-muted-foreground">Run AI models locally</span>
								</div>
							</Label>
						</RadioGroup.Root>

						<!-- Terminal Agent Options -->
						{#if engine === 'terminal-agent'}
							<div class="flex flex-col gap-3 rounded-lg border border-border bg-muted/30 p-4">
								<Label class="text-sm font-medium">Select Terminal Agent</Label>
								<Select.Root type="single" bind:value={terminalAgent}>
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
						{#if engine === 'ollama'}
							<div class="flex flex-col gap-3 rounded-lg border border-border bg-muted/30 p-4">
								<Label class="text-sm font-medium">Select Model</Label>
								<Select.Root type="single" bind:value={ollamaModel}>
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
				{:else if activeSection === 'prompt'}
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

						<!-- Prompt Editor -->
						<div class="flex flex-col gap-3">
							<div class="flex items-center gap-2">
								<MessageSquareText class="size-4 text-muted-foreground" />
								<Label class="text-sm font-medium">Prompt Template</Label>
							</div>
							<Textarea
								bind:value={prompt}
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
				{/if}
			</div>
		</div>
	</main>
</div>
