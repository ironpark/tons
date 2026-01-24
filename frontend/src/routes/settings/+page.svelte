<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Settings from '@lucide/svelte/icons/settings';
	import Sliders from '@lucide/svelte/icons/sliders-horizontal';
	import Cpu from '@lucide/svelte/icons/cpu';
	import MessageSquareText from '@lucide/svelte/icons/message-square-text';
	import { getActiveSection, setActiveSection, loadConfig, sections } from './settings.svelte.ts';
	import GeneralSection from './GeneralSection.svelte';
	import EngineSection from './EngineSection.svelte';
	import PromptSection from './PromptSection.svelte';

	const sectionIcons = {
		general: Sliders,
		engine: Cpu,
		prompt: MessageSquareText
	};

	const activeSection = $derived(getActiveSection());

	onMount(() => {
		loadConfig();
	});
</script>

<div
	class="relative flex h-screen flex-col overflow-hidden bg-background font-sans text-foreground"
>
	<!-- Background grid -->
	<div
		class="pointer-events-none absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.015)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.015)_1px,transparent_1px)] bg-[size:60px_60px]"
	></div>

	<!-- Header -->
	<header
		class="relative z-10 flex items-center gap-3 border-b border-border bg-background/80 py-2.5 pr-4 pl-20 backdrop-blur-xl [-webkit-app-region:drag]"
	>
		<div class="[-webkit-app-region:no-drag]">
			<Button
				variant="ghost"
				size="icon"
				href="/"
				class="text-muted-foreground hover:text-foreground"
			>
				<ArrowLeft class="size-[18px]" />
			</Button>
		</div>
		<div class="flex items-center gap-2">
			<div
				class="flex h-9 w-9 items-center justify-center rounded-md bg-accent text-accent-foreground"
			>
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
				{@const Icon = sectionIcons[section.id as keyof typeof sectionIcons]}
				<button
					onclick={() => setActiveSection(section.id)}
					class="flex items-center gap-2 rounded-lg px-3 py-2 text-left text-sm font-medium transition-colors {activeSection ===
					section.id
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'}"
				>
					<Icon class="size-4" />
					{section.label}
				</button>
			{/each}
		</nav>

		<!-- Content -->
		<div class="flex-1 overflow-y-auto p-6">
			<div class="h-full">
				{#if activeSection === 'general'}
					<GeneralSection />
				{:else if activeSection === 'engine'}
					<EngineSection />
				{:else if activeSection === 'prompt'}
					<PromptSection />
				{/if}
			</div>
		</div>
	</main>
</div>
