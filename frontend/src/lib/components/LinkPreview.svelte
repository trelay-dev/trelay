<script lang="ts">
	import { preview, type LinkPreview } from '$lib/utils/api';
	import { onMount } from 'svelte';
	
	interface Props {
		url: string;
	}
	
	let { url }: Props = $props();
	
	let previewData = $state<LinkPreview | null>(null);
	let loading = $state(true);
	let error = $state(false);
	
	onMount(async () => {
		await loadPreview();
	});
	
	async function loadPreview() {
		loading = true;
		error = false;
		try {
			const res = await preview.fetch(url);
			if (res.success && res.data) {
				previewData = res.data;
			} else {
				error = true;
			}
		} catch (e) {
			error = true;
		} finally {
			loading = false;
		}
	}
</script>

<div class="link-preview">
	{#if loading}
		<div class="preview-loading">
			<div class="preview-spinner"></div>
			<span>Loading preview...</span>
		</div>
	{:else if error || !previewData}
		<div class="preview-error">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
				<circle cx="12" cy="12" r="10"/>
				<line x1="12" y1="8" x2="12" y2="12"/>
				<line x1="12" y1="16" x2="12.01" y2="16"/>
			</svg>
			<span>Preview unavailable</span>
		</div>
	{:else}
		<div class="preview-card">
			{#if previewData.image_url}
				<div class="preview-image">
					<img src={previewData.image_url} alt="" loading="lazy" />
				</div>
			{/if}
			<div class="preview-content">
				{#if previewData.title}
					<h4 class="preview-title">{previewData.title}</h4>
				{/if}
				{#if previewData.description}
					<p class="preview-description">{previewData.description}</p>
				{/if}
				<span class="preview-url">{url}</span>
			</div>
		</div>
	{/if}
</div>

<style>
	.link-preview {
		margin-top: var(--space-4);
	}
	
	.preview-loading, .preview-error {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-4);
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
		color: var(--text-muted);
		font-size: var(--text-sm);
	}
	
	.preview-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid var(--border-light);
		border-top-color: var(--accent-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
	
	.preview-card {
		display: flex;
		flex-direction: column;
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		overflow: hidden;
		background: var(--bg-secondary);
	}
	
	.preview-image {
		width: 100%;
		height: 140px;
		overflow: hidden;
		background: var(--bg-tertiary);
	}
	
	.preview-image img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	
	.preview-content {
		padding: var(--space-3);
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.preview-title {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-primary);
		line-height: 1.3;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	
	.preview-description {
		font-size: var(--text-xs);
		color: var(--text-secondary);
		line-height: 1.4;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	
	.preview-url {
		font-size: var(--text-xs);
		color: var(--text-muted);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
