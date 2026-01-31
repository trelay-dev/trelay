<script lang="ts">
	interface Props {
		type?: 'text' | 'password' | 'email' | 'url' | 'number' | 'search';
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		readonly?: boolean;
		error?: string;
		label?: string;
		id?: string;
		name?: string;
		oninput?: (e: Event) => void;
		onkeydown?: (e: KeyboardEvent) => void;
	}
	
	let {
		type = 'text',
		value = $bindable(''),
		placeholder = '',
		disabled = false,
		readonly = false,
		error = '',
		label = '',
		id = '',
		name = '',
		oninput,
		onkeydown
	}: Props = $props();
</script>

<div class="input-wrapper">
	{#if label}
		<label for={id} class="label">{label}</label>
	{/if}
	<input
		{type}
		{id}
		{name}
		{placeholder}
		{disabled}
		{readonly}
		bind:value
		class="input"
		class:has-error={!!error}
		{oninput}
		{onkeydown}
	/>
	{#if error}
		<span class="error">{error}</span>
	{/if}
</div>

<style>
	.input-wrapper {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.label {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.input {
		height: 40px;
		padding: 0 var(--space-3);
		font-family: var(--font-sans);
		font-size: var(--text-base);
		color: var(--text-primary);
		background: var(--input-bg);
		border: 1px solid var(--input-border);
		border-radius: var(--radius-md);
		transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
	}
	
	.input::placeholder {
		color: var(--text-muted);
	}
	
	.input:focus {
		outline: none;
		border-color: var(--input-focus);
		box-shadow: 0 0 0 3px rgba(233, 30, 99, 0.1);
	}
	
	.input:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	
	.input.has-error {
		border-color: var(--color-error);
	}
	
	.input.has-error:focus {
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
	}
	
	.error {
		font-size: var(--text-xs);
		color: var(--color-error);
	}
</style>
