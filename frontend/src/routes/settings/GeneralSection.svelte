<script lang="ts">
	import * as Select from '$lib/components/ui/select';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { Label } from '$lib/components/ui/label';
	import Sun from '@lucide/svelte/icons/sun';
	import Moon from '@lucide/svelte/icons/moon';
	import Monitor from '@lucide/svelte/icons/monitor';
	import Globe from '@lucide/svelte/icons/globe';
	import { Theme } from '$lib/bindings/github.com/ironpark/tons/internal/config/models';
	import {
		getGeneralConfig,
		getSelectedLanguage,
		languages,
		setTheme,
		setLanguage
	} from './settings.svelte.ts';

	const generalConfig = $derived(getGeneralConfig());
	const selectedLanguage = $derived(getSelectedLanguage());
</script>

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
		<RadioGroup.Root
			value={generalConfig.theme}
			onValueChange={(value) => setTheme(value as Theme)}
			class="flex gap-2"
		>
			<Label
				class="flex flex-1 cursor-pointer flex-col items-center gap-2 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
			>
				<RadioGroup.Item value={Theme.ThemeLight} class="sr-only" />
				<Sun class="size-5" />
				<span class="text-sm font-medium">Light</span>
			</Label>
			<Label
				class="flex flex-1 cursor-pointer flex-col items-center gap-2 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
			>
				<RadioGroup.Item value={Theme.ThemeDark} class="sr-only" />
				<Moon class="size-5" />
				<span class="text-sm font-medium">Dark</span>
			</Label>
			<Label
				class="flex flex-1 cursor-pointer flex-col items-center gap-2 rounded-lg border border-border p-3 transition-colors hover:bg-accent/50 has-[[data-state=checked]]:border-primary has-[[data-state=checked]]:bg-accent"
			>
				<RadioGroup.Item value={Theme.ThemeSystem} class="sr-only" />
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
		<Select.Root
			type="single"
			value={generalConfig.language}
			onValueChange={(value) => setLanguage(value)}
		>
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
