<script lang="ts">
	interface DataPoint {
		label: string;
		value: number;
	}

	interface Props {
		data: DataPoint[];
		/** Height of the plot area (bars + grid), in px */
		plotHeight?: number;
	}

	let { data, plotHeight = 160 }: Props = $props();

	/** Upper bound for the Y scale (nice round number ≥ max data value). */
	function computeAxisMax(m: number): number {
		if (m <= 0) return 1;
		if (m <= 4) return 4;
		const exp = Math.floor(Math.log10(m));
		const f = m / 10 ** exp;
		const nf = f <= 1 ? 1 : f <= 2 ? 2 : f <= 5 ? 5 : 10;
		return nf * 10 ** exp;
	}

	let axisMax = $derived.by(() => {
		const vals = data.map((d) => d.value);
		if (vals.length === 0) return 1;
		return computeAxisMax(Math.max(...vals));
	});

	const DIVISIONS = 4;

	let yTicks = $derived.by(() => {
		const am = axisMax;
		return Array.from({ length: DIVISIONS + 1 }, (_, i) => {
			const raw = (am * i) / DIVISIONS;
			return {
				value: raw,
				label: Number.isInteger(raw) ? String(raw) : raw.toFixed(1),
				pctFromBottom: i / DIVISIONS
			};
		});
	});

	function barHeightPct(value: number): number {
		if (axisMax <= 0) return 0;
		return (value / axisMax) * 100;
	}

	let axisMaxLabel = $derived(
		Number.isInteger(axisMax) ? String(axisMax) : axisMax.toFixed(1)
	);
</script>

<div class="chart">
	<div class="chart-layout" style="--plot-h: {plotHeight}px">
		<div class="y-axis" aria-hidden="true">
			{#each [...yTicks].reverse() as tick}
				<span class="y-label">{tick.label}</span>
			{/each}
		</div>
		<div class="plot-stack">
			<div class="plot">
				{#each yTicks as tick}
					<div
						class="grid-line"
						style="bottom: {tick.pctFromBottom * 100}%"
					></div>
				{/each}
				<div class="bars-row">
					{#each data as point}
						<div class="bar-column">
							<div class="bar-track">
								{#if point.value > 0}
									<span class="bar-value">{point.value}</span>
								{:else}
									<span class="bar-value bar-value--muted">0</span>
								{/if}
								<div
									class="bar"
									class:bar-visible={point.value > 0}
									style="height: {barHeightPct(point.value)}%"
									title="{point.label}: {point.value} clicks"
								></div>
							</div>
						</div>
					{/each}
				</div>
			</div>
			<div class="x-labels-row">
				{#each data as point}
					<span class="x-label">{point.label}</span>
				{/each}
			</div>
		</div>
	</div>
	<p class="chart-legend">Clicks per day (Y-axis 0–{axisMaxLabel})</p>
</div>

<style>
	.chart {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		width: 100%;
	}

	.chart-layout {
		display: grid;
		grid-template-columns: 2.25rem minmax(0, 1fr);
		gap: var(--space-2);
		align-items: start;
	}

	.y-axis {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		height: var(--plot-h);
		flex-shrink: 0;
		align-items: flex-end;
		align-self: start;
		box-sizing: border-box;
	}

	.y-label {
		font-size: var(--text-xs);
		font-variant-numeric: tabular-nums;
		line-height: 1;
		color: var(--text-muted);
	}

	.plot-stack {
		display: flex;
		flex-direction: column;
		min-width: 0;
		gap: var(--space-2);
	}

	.plot {
		position: relative;
		height: var(--plot-h);
		border-radius: var(--radius-sm);
		background: var(--bg-secondary);
	}

	.grid-line {
		position: absolute;
		left: 0;
		right: 0;
		height: 1px;
		margin-bottom: -0.5px;
		background: var(--border-light);
		pointer-events: none;
	}

	.bars-row {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: stretch;
		gap: var(--space-1);
		padding: var(--space-2);
		box-sizing: border-box;
		z-index: 1;
	}

	.bar-column {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		min-width: 0;
		min-height: 0;
		height: 100%;
	}

	.bar-track {
		flex: 1;
		width: 100%;
		max-width: 36px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: flex-end;
		min-height: 0;
		gap: var(--space-1);
	}

	.bar-value {
		flex-shrink: 0;
		font-size: var(--text-xs);
		font-weight: var(--font-medium);
		font-variant-numeric: tabular-nums;
		line-height: 1;
		color: var(--text-secondary);
	}

	.bar-value--muted {
		font-weight: var(--font-normal);
		color: var(--text-muted);
		opacity: 0.75;
	}

	.bar {
		width: 100%;
		min-height: 0;
		background: var(--accent-primary);
		border-radius: var(--radius-sm) var(--radius-sm) 0 0;
		transition: height var(--transition-slow);
	}

	.bar-visible {
		min-height: 2px;
	}

	.bar:hover {
		background: var(--accent-primary-hover);
	}

	.x-labels-row {
		display: flex;
		gap: var(--space-1);
		padding: 0 var(--space-2);
		box-sizing: border-box;
	}

	.x-label {
		flex: 1;
		min-width: 0;
		font-size: var(--text-xs);
		color: var(--text-muted);
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.chart-legend {
		margin: 0;
		font-size: var(--text-xs);
		color: var(--text-muted);
	}
</style>
