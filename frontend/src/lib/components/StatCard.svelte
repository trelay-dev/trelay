<script lang="ts">
	interface Props {
		label: string;
		value: string | number;
		change?: number;
		icon?: 'link' | 'click' | 'folder' | 'chart';
	}
	
	let { label, value, change, icon }: Props = $props();
</script>

<div class="stat-card">
	<div class="stat-icon">
		{#if icon === 'link'}
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
				<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
				<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
			</svg>
		{:else if icon === 'click'}
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
				<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
			</svg>
		{:else if icon === 'folder'}
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
				<path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
			</svg>
		{:else if icon === 'chart'}
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
				<line x1="18" y1="20" x2="18" y2="10"/>
				<line x1="12" y1="20" x2="12" y2="4"/>
				<line x1="6" y1="20" x2="6" y2="14"/>
			</svg>
		{/if}
	</div>
	<div class="stat-content">
		<span class="stat-label">{label}</span>
		<span class="stat-value">{value}</span>
		{#if change !== undefined}
			<span class="stat-change" class:positive={change > 0} class:negative={change < 0}>
				{change > 0 ? '+' : ''}{change}%
			</span>
		{/if}
	</div>
</div>

<style>
	.stat-card {
		display: flex;
		align-items: flex-start;
		gap: var(--space-4);
		padding: var(--space-5);
		background: var(--card-bg);
		border: 1px solid var(--card-border);
		border-radius: var(--radius-lg);
	}
	
	.stat-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 44px;
		height: 44px;
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
		color: var(--accent-primary);
		flex-shrink: 0;
	}
	
	.stat-content {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.stat-label {
		font-size: var(--text-sm);
		color: var(--text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.025em;
	}
	
	.stat-value {
		font-size: var(--text-3xl);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
		line-height: 1;
	}
	
	.stat-change {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
	}
	
	.stat-change.positive {
		color: var(--color-success);
	}
	
	.stat-change.negative {
		color: var(--color-error);
	}
</style>
