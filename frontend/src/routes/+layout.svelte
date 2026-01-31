<script lang="ts">
	import '$lib/styles/global.css';
	import { Navbar } from '$lib/components';
	import { theme } from '$lib/stores/theme';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	
	let { children } = $props();
	
	onMount(() => {
		if (browser) {
			document.documentElement.setAttribute('data-theme', $theme);
		}
	});
	
	$effect(() => {
		if (browser) {
			document.documentElement.setAttribute('data-theme', $theme);
		}
	});
</script>

<div class="app">
	<Navbar />
	<main class="main">
		{@render children()}
	</main>
</div>

<style>
	.app {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}
	
	.main {
		flex: 1;
		padding: var(--space-6) 0;
	}
</style>
