<script lang="ts">
	import type { Snippet } from 'svelte';
	
	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		loading?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: (e: MouseEvent) => void;
		children: Snippet;
	}
	
	let {
		variant = 'primary',
		size = 'md',
		disabled = false,
		loading = false,
		type = 'button',
		onclick,
		children
	}: Props = $props();
</script>

<button
	{type}
	class="btn btn-{variant} btn-{size}"
	disabled={disabled || loading}
	onclick={onclick}
>
	{#if loading}
		<span class="spinner"></span>
	{/if}
	{@render children()}
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		font-family: var(--font-sans);
		font-weight: var(--font-medium);
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
		white-space: nowrap;
	}
	
	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	
	/* Sizes */
	.btn-sm {
		height: 32px;
		padding: 0 var(--space-3);
		font-size: var(--text-sm);
	}
	
	.btn-md {
		height: 38px;
		padding: 0 var(--space-4);
		font-size: var(--text-sm);
	}
	
	.btn-lg {
		height: 44px;
		padding: 0 var(--space-6);
		font-size: var(--text-base);
	}
	
	/* Variants */
	.btn-primary {
		background: var(--accent-primary);
		color: white;
		border-color: var(--accent-primary);
	}
	
	.btn-primary:hover:not(:disabled) {
		background: var(--accent-primary-hover);
		border-color: var(--accent-primary-hover);
	}
	
	.btn-secondary {
		background: var(--bg-primary);
		color: var(--text-primary);
		border-color: var(--border-medium);
	}
	
	.btn-secondary:hover:not(:disabled) {
		background: var(--bg-hover);
		border-color: var(--border-dark);
	}
	
	.btn-ghost {
		background: transparent;
		color: var(--text-secondary);
		border-color: transparent;
	}
	
	.btn-ghost:hover:not(:disabled) {
		background: var(--bg-hover);
		color: var(--text-primary);
	}
	
	.btn-danger {
		background: var(--color-error);
		color: white;
		border-color: var(--color-error);
	}
	
	.btn-danger:hover:not(:disabled) {
		background: #dc2626;
		border-color: #dc2626;
	}
	
	/* Spinner */
	.spinner {
		width: 14px;
		height: 14px;
		border: 2px solid currentColor;
		border-top-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
