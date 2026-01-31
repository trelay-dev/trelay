<script lang="ts">
	import type { Snippet } from 'svelte';
	
	interface Props {
		padding?: 'none' | 'sm' | 'md' | 'lg';
		hoverable?: boolean;
		children: Snippet;
	}
	
	let {
		padding = 'md',
		hoverable = false,
		children
	}: Props = $props();
</script>

<div class="card padding-{padding}" class:hoverable>
	{@render children()}
</div>

<style>
	.card {
		background: var(--card-bg);
		border: 1px solid var(--card-border);
		border-radius: var(--radius-lg);
		transition: box-shadow var(--transition-fast), border-color var(--transition-fast);
	}
	
	.hoverable:hover {
		border-color: var(--border-medium);
		box-shadow: var(--shadow-md);
	}
	
	.padding-none { padding: 0; }
	.padding-sm { padding: var(--space-3); }
	.padding-md { padding: var(--space-4); }
	.padding-lg { padding: var(--space-6); }
</style>
