<script lang="ts">
	import './layout.css';
	import { page } from '$app/state';
	import { locales, localizeHref } from '$lib/paraglide/runtime';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { fly, fade } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
{#key page.url.pathname}
	<div
		class="page-wrapper"
		in:fly={{ x: 20, duration: 200, delay: 100, easing: cubicOut }}
		out:fade={{ duration: 100 }}
	>
		{@render children()}
	</div>
{/key}
<div style="display:none">
	{#each locales as locale (locale)}
		<a href={localizeHref(page.url.pathname, { locale })}>
			{locale}
		</a>
	{/each}
</div>
