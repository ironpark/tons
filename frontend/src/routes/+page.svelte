<script lang="ts">
	import TranslatePanel from '$lib/components/TranslatePanel.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import Settings from '@lucide/svelte/icons/settings';
	import ArrowLeftRight from '@lucide/svelte/icons/arrow-left-right';
	import Languages from '@lucide/svelte/icons/languages';

	let sourceText = $state('');
	let translatedText = $state('');
	let sourceLangValue = $state('en');
	let targetLangValue = $state('ko');
	let isTranslating = $state(false);

	const languages = [
		{ value: 'en', label: 'English', flag: '🇺🇸' },
		{ value: 'ko', label: '한국어', flag: '🇰🇷' },
		{ value: 'ja', label: '日本語', flag: '🇯🇵' },
		{ value: 'zh', label: '中文', flag: '🇨🇳' },
		{ value: 'es', label: 'Español', flag: '🇪🇸' },
		{ value: 'fr', label: 'Français', flag: '🇫🇷' },
		{ value: 'de', label: 'Deutsch', flag: '🇩🇪' },
		{ value: 'pt', label: 'Português', flag: '🇧🇷' },
		{ value: 'ru', label: 'Русский', flag: '🇷🇺' },
		{ value: 'ar', label: 'العربية', flag: '🇸🇦' }
	];

	const sourceLang = $derived(languages.find((l) => l.value === sourceLangValue) ?? languages[0]);
	const targetLang = $derived(languages.find((l) => l.value === targetLangValue) ?? languages[1]);

	function swapLanguages() {
		const tempLang = sourceLangValue;
		const tempText = sourceText;
		sourceLangValue = targetLangValue;
		targetLangValue = tempLang;
		sourceText = translatedText;
		translatedText = tempText;
	}

	function handleTranslate() {
		if (!sourceText.trim()) return;
		isTranslating = true;
		setTimeout(() => {
			translatedText = `[${targetLangValue.toUpperCase()}] ${sourceText}`;
			isTranslating = false;
		}, 600);
	}

	function clearAll() {
		sourceText = '';
		translatedText = '';
	}

	$effect(() => {
		if (sourceText) {
			const timeout = setTimeout(handleTranslate, 500);
			return () => clearTimeout(timeout);
		} else {
			translatedText = '';
		}
	});
</script>

<div class="relative flex h-screen flex-col overflow-hidden bg-bg font-sans text-text">
	<!-- Background grid -->
	<div class="pointer-events-none absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.015)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.015)_1px,transparent_1px)] bg-[size:60px_60px]"></div>

	<!-- Header -->
	<header class="relative z-10 flex items-center justify-between border-b border-border bg-bg/80 py-2.5 pr-4 pl-20 backdrop-blur-xl [-webkit-app-region:drag]">
		<div class="flex items-center gap-2">
			<div class="flex h-9 w-9 items-center justify-center rounded-md bg-accent text-white">
				<Languages class="size-5" />
			</div>
			<h1 class="text-xl font-semibold tracking-tight text-text">Tons</h1>
		</div>
		<div class="[-webkit-app-region:no-drag]">
			<Button variant="ghost" size="icon" href="/settings" class="text-text-muted hover:text-text">
				<Settings class="size-[18px]" />
			</Button>
		</div>
	</header>

	<!-- Main Content -->
	<main class="relative z-10 flex flex-1 flex-col gap-3 overflow-hidden p-3">
		<!-- Language Selector Bar -->
		<div class="flex items-center justify-center gap-2">
			<Select.Root type="single" bind:value={sourceLangValue}>
				<Select.Trigger class="min-w-40 bg-surface border-border hover:bg-surface-elevated">
					<span class="flex items-center gap-2">
						<span>{sourceLang.flag}</span>
						<span>{sourceLang.label}</span>
					</span>
				</Select.Trigger>
				<Select.Content>
					{#each languages as lang (lang.value)}
						<Select.Item value={lang.value} label={lang.label}>
							<span class="flex items-center gap-2">
								<span>{lang.flag}</span>
								<span>{lang.label}</span>
							</span>
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>

			<Button variant="ghost" size="icon" onclick={swapLanguages} class="text-text-muted hover:text-text">
				<ArrowLeftRight class="size-[18px]" />
			</Button>

			<Select.Root type="single" bind:value={targetLangValue}>
				<Select.Trigger class="min-w-40 bg-surface border-border hover:bg-surface-elevated">
					<span class="flex items-center gap-2">
						<span>{targetLang.flag}</span>
						<span>{targetLang.label}</span>
					</span>
				</Select.Trigger>
				<Select.Content>
					{#each languages as lang (lang.value)}
						<Select.Item value={lang.value} label={lang.label}>
							<span class="flex items-center gap-2">
								<span>{lang.flag}</span>
								<span>{lang.label}</span>
							</span>
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<!-- Translation Panels -->
		<div class="grid min-h-0 flex-1 grid-cols-2 gap-3">
			<TranslatePanel
				label={sourceLang.label}
				bind:value={sourceText}
				placeholder="Enter text to translate..."
				onClear={clearAll}
			/>
			<TranslatePanel
				label={targetLang.label}
				value={translatedText}
				placeholder="Translation will appear here..."
				readonly
				loading={isTranslating}
			/>
		</div>
	</main>
</div>
