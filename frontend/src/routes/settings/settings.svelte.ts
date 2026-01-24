import * as SettingService from '$lib/bindings/github.com/ironpark/tons/internal/services/settingservice';
import {
	GeneralConfig,
	EngineConfig,
	PromptConfig,
	Theme,
	EngineType,
	TerminalAgentType
} from '$lib/bindings/github.com/ironpark/tons/internal/config/models';

// State
let generalConfig = $state(new GeneralConfig({ theme: Theme.ThemeSystem, language: 'system' }));
let engineConfig = $state(new EngineConfig({ type: EngineType.EngineInternal }));
let promptConfig = $state(new PromptConfig());
let activeSection = $state('general');

// Options
export const languages = [
	{ value: 'system', label: 'System Default' },
	{ value: 'en', label: 'English' },
	{ value: 'ko', label: '한국어' },
	{ value: 'ja', label: '日本語' }
];

export const terminalAgents = [
	{ value: TerminalAgentType.AgentClaudeCode, label: 'Claude Code' },
	{ value: TerminalAgentType.AgentGeminiCLI, label: 'Gemini CLI' },
	{ value: TerminalAgentType.AgentCodex, label: 'Codex' }
];

export const ollamaModels = [
	{ value: 'llama3.2', label: 'Llama 3.2' },
	{ value: 'llama3.1', label: 'Llama 3.1' },
	{ value: 'mistral', label: 'Mistral' },
	{ value: 'codellama', label: 'Code Llama' },
	{ value: 'gemma2', label: 'Gemma 2' }
];

export const sections = [
	{ id: 'general', label: 'General' },
	{ id: 'engine', label: 'Engine' },
	{ id: 'prompt', label: 'Prompt' }
];

// Default prompt
export const defaultPrompt = `Translate the following text from {{source_lang}} to {{target_lang}}.
Keep the original formatting and tone.
Only return the translated text without any explanations.

Text to translate:
{{text}}`;

// Getters for reactive state
export function getGeneralConfig() {
	return generalConfig;
}

export function getEngineConfig() {
	return engineConfig;
}

export function getPromptConfig() {
	return promptConfig;
}

export function getActiveSection() {
	return activeSection;
}

// Derived values
export function getSelectedLanguage() {
	return languages.find((l) => l.value === generalConfig.language) ?? languages[0];
}

export function getSelectedTerminalAgent() {
	return terminalAgents.find((a) => a.value === engineConfig.terminalAgent.selected) ?? terminalAgents[0];
}

export function getSelectedOllamaModel() {
	return ollamaModels.find((m) => m.value === engineConfig.ollama.model) ?? ollamaModels[0];
}

// Load config
export async function loadConfig() {
	const config = await SettingService.GetCurrentConfig();
	if (config) {
		generalConfig = config.general;
		engineConfig = config.engine;
		promptConfig = config.prompt;
	}
}

// Save handlers
export async function saveGeneralConfig() {
	await SettingService.UpdateGeneralConfig(generalConfig);
}

export async function saveEngineConfig() {
	await SettingService.UpdateEngineConfig(engineConfig);
}

export async function savePromptConfig() {
	await SettingService.UpdatePromptConfig(promptConfig);
}

// Setters
export function setActiveSection(section: string) {
	activeSection = section;
}

export function setTheme(theme: Theme) {
	generalConfig = { ...generalConfig, theme };
	saveGeneralConfig();
}

export function setLanguage(language: string) {
	generalConfig = { ...generalConfig, language };
	saveGeneralConfig();
}

export function setEngineType(type: EngineType) {
	engineConfig = { ...engineConfig, type };
	saveEngineConfig();
}

export function setTerminalAgent(selected: TerminalAgentType) {
	engineConfig = {
		...engineConfig,
		terminalAgent: { ...engineConfig.terminalAgent, selected }
	};
	saveEngineConfig();
}

export function setOllamaModel(model: string) {
	engineConfig = {
		...engineConfig,
		ollama: { ...engineConfig.ollama, model }
	};
	saveEngineConfig();
}

export function setPromptTemplate(template: string) {
	promptConfig = { ...promptConfig, template };
}

export function resetPrompt() {
	promptConfig = { ...promptConfig, template: defaultPrompt };
	savePromptConfig();
}
