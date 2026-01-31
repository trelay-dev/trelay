<script lang="ts">
	import type { Link } from '$lib/utils/api';
	
	interface Props {
		link: Link;
		baseUrl?: string;
		ondelete?: (slug: string) => void;
		onedit?: (link: Link) => void;
	}
	
	let { link, baseUrl = '', ondelete, onedit }: Props = $props();
	
	let copied = $state(false);
	
	function copyLink() {
		const url = `${baseUrl}/${link.slug}`;
		navigator.clipboard.writeText(url);
		copied = true;
		setTimeout(() => copied = false, 2000);
	}
	
	function formatDate(date: string) {
		return new Date(date).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
	
	function truncateUrl(url: string, max = 50) {
		if (url.length <= max) return url;
		return url.slice(0, max) + '...';
	}
</script>

<div class="link-row">
	<div class="link-main">
		<div class="link-slug-row">
			<span class="link-slug mono">/{link.slug}</span>
			{#if link.has_password}
				<span class="badge badge-warning">Protected</span>
			{/if}
			{#if link.is_one_time}
				<span class="badge badge-info">One-time</span>
			{/if}
		</div>
		<a href={link.original_url} target="_blank" rel="noopener" class="link-url truncate">
			{truncateUrl(link.original_url)}
		</a>
	</div>
	
	<div class="link-stats">
		<span class="click-count">
			<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
				<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
			</svg>
			{link.click_count.toLocaleString()}
		</span>
		<span class="link-date">{formatDate(link.created_at)}</span>
	</div>
	
	<div class="link-actions">
		<button class="action-btn" onclick={copyLink} title="Copy link">
			{#if copied}
				<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<polyline points="20 6 9 17 4 12"/>
				</svg>
			{:else}
				<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
					<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
				</svg>
			{/if}
		</button>
		{#if onedit}
			<button class="action-btn" onclick={() => onedit(link)} title="Edit">
				<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
					<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
				</svg>
			</button>
		{/if}
		{#if ondelete}
			<button class="action-btn action-btn-danger" onclick={() => ondelete(link.slug)} title="Delete">
				<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<polyline points="3 6 5 6 21 6"/>
					<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
				</svg>
			</button>
		{/if}
	</div>
</div>

<style>
	.link-row {
		display: grid;
		grid-template-columns: 1fr auto auto;
		align-items: center;
		gap: var(--space-4);
		padding: var(--space-4);
		border-bottom: 1px solid var(--border-light);
		transition: background var(--transition-fast);
	}
	
	.link-row:hover {
		background: var(--bg-hover);
	}
	
	.link-row:last-child {
		border-bottom: none;
	}
	
	.link-main {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
		min-width: 0;
	}
	
	.link-slug-row {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}
	
	.link-slug {
		font-size: var(--text-base);
		font-weight: var(--font-medium);
		color: var(--text-primary);
	}
	
	.link-url {
		font-size: var(--text-sm);
		color: var(--text-tertiary);
	}
	
	.link-url:hover {
		color: var(--accent-primary);
	}
	
	.badge {
		padding: var(--space-1) var(--space-2);
		font-size: var(--text-xs);
		font-weight: var(--font-medium);
		border-radius: var(--radius-sm);
	}
	
	.badge-warning {
		background: rgba(245, 158, 11, 0.1);
		color: var(--color-warning);
	}
	
	.badge-info {
		background: rgba(59, 130, 246, 0.1);
		color: var(--color-info);
	}
	
	.link-stats {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: var(--space-1);
	}
	
	.click-count {
		display: flex;
		align-items: center;
		gap: var(--space-1);
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.link-date {
		font-size: var(--text-xs);
		color: var(--text-muted);
	}
	
	.link-actions {
		display: flex;
		align-items: center;
		gap: var(--space-1);
	}
	
	.action-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		color: var(--text-tertiary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.action-btn:hover {
		background: var(--bg-tertiary);
		color: var(--text-primary);
	}
	
	.action-btn-danger:hover {
		background: rgba(239, 68, 68, 0.1);
		color: var(--color-error);
	}
	
	@media (max-width: 640px) {
		.link-row {
			grid-template-columns: 1fr;
			gap: var(--space-3);
		}
		
		.link-stats {
			flex-direction: row;
			align-items: center;
			justify-content: flex-start;
			gap: var(--space-4);
		}
		
		.link-actions {
			justify-content: flex-start;
		}
	}
</style>
