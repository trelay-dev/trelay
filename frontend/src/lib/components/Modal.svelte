<script lang="ts">
	import type { Snippet } from 'svelte';
	
	interface Props {
		open: boolean;
		title?: string;
		onclose: () => void;
		children: Snippet;
		footer?: Snippet;
	}
	
	let { open, title = '', onclose, children, footer }: Props = $props();
	
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose();
	}
	
	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) onclose();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
	<div class="modal-backdrop" onclick={handleBackdropClick}>
		<div class="modal" role="dialog" aria-modal="true">
			{#if title}
				<div class="modal-header">
					<h3 class="modal-title">{title}</h3>
					<button class="close-btn" onclick={onclose} aria-label="Close">
						<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<line x1="18" y1="6" x2="6" y2="18"/>
							<line x1="6" y1="6" x2="18" y2="18"/>
						</svg>
					</button>
				</div>
			{/if}
			<div class="modal-body">
				{@render children()}
			</div>
			{#if footer}
				<div class="modal-footer">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		z-index: var(--z-modal);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-4);
		background: rgba(0, 0, 0, 0.5);
		animation: fadeIn 150ms ease;
	}
	
	.modal {
		width: 100%;
		max-width: 480px;
		max-height: 90vh;
		overflow-y: auto;
		background: var(--card-bg);
		border: 1px solid var(--card-border);
		border-radius: var(--radius-xl);
		box-shadow: var(--shadow-lg);
		animation: slideUp 150ms ease;
	}
	
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-4) var(--space-5);
		border-bottom: 1px solid var(--border-light);
	}
	
	.modal-title {
		font-size: var(--text-lg);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.close-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		background: transparent;
		border: none;
		border-radius: var(--radius-md);
		color: var(--text-tertiary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.close-btn:hover {
		background: var(--bg-hover);
		color: var(--text-primary);
	}
	
	.modal-body {
		padding: var(--space-5);
	}
	
	.modal-footer {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-3);
		padding: var(--space-4) var(--space-5);
		border-top: 1px solid var(--border-light);
	}
	
	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
	
	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
