<script lang="ts">
	interface DataPoint {
		label: string;
		value: number;
	}
	
	interface Props {
		data: DataPoint[];
		height?: number;
	}
	
	let { data, height = 200 }: Props = $props();
	
	let maxValue = $derived(Math.max(...data.map(d => d.value), 1));
</script>

<div class="chart" style="height: {height}px">
	<div class="chart-bars">
		{#each data as point, i}
			<div class="bar-wrapper">
				<div 
					class="bar" 
					style="height: {(point.value / maxValue) * 100}%"
					title="{point.label}: {point.value}"
				></div>
				<span class="bar-label">{point.label}</span>
			</div>
		{/each}
	</div>
</div>

<style>
	.chart {
		display: flex;
		flex-direction: column;
	}
	
	.chart-bars {
		display: flex;
		align-items: flex-end;
		gap: var(--space-1);
		height: 100%;
		padding-bottom: var(--space-6);
	}
	
	.bar-wrapper {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		height: 100%;
	}
	
	.bar {
		width: 100%;
		max-width: 32px;
		min-height: 2px;
		background: var(--accent-primary);
		border-radius: var(--radius-sm) var(--radius-sm) 0 0;
		transition: height var(--transition-slow);
	}
	
	.bar:hover {
		background: var(--accent-primary-hover);
	}
	
	.bar-label {
		position: absolute;
		bottom: 0;
		font-size: var(--text-xs);
		color: var(--text-muted);
		transform: translateY(100%);
		padding-top: var(--space-2);
	}
	
	.bar-wrapper {
		position: relative;
	}
</style>
