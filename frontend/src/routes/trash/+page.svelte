<script lang="ts">
	import { Button, Card, LinkRow } from '$lib/components';
	import { links, type Link } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	let deletedLinks = $state<Link[]>([]);
	let loading = $state(true);
	
	onMount(async () => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		
		await loadDeletedLinks();
	});
	
	async function loadDeletedLinks() {
		loading = true;
		try {
			const res = await links.list({ only_deleted: true });
			if (res.success && res.data) {
				deletedLinks = res.data;
			}
		} catch (e) {
			console.error('Failed to load deleted links:', e);
		} finally {
			loading = false;
		}
	}
	
	async function handleRestore(slug: string) {
		try {
			const res = await links.restore(slug);
			if (res.success) {
				deletedLinks = deletedLinks.filter(l => l.slug !== slug);
			}
		} catch (e) {
			console.error('Failed to restore link:', e);
		}
	}
	
	async function handlePermanentDelete(slug: string) {
		if (!confirm('Permanently delete this link? This cannot be undone.')) return;
		
		const res = await links.delete(slug, true);
		if (res.success) {
			await loadDeletedLinks();
		}
	}
	
	async function emptyTrash() {
		if (!confirm(`Permanently delete all ${deletedLinks.length} links? This cannot be undone.`)) return;
		
		const slugs = deletedLinks.map(l => l.slug);
		const res = await links.bulkDelete(slugs, true);
		if (res.success) {
			await loadDeletedLinks();
		}
	}
</script>

<svelte:head>
	<title>Trash - Trelay</title>
</svelte:head>

<div class="trash-page container">
	<header class="page-header">
		<div class="page-header-content">
			<h1 class="page-title">Trash</h1>
			<p class="page-subtitle">{deletedLinks.length} deleted links</p>
		</div>
		{#if deletedLinks.length > 0}
			<Button variant="danger" onclick={emptyTrash}>
				<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<polyline points="3 6 5 6 21 6"/>
					<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
				</svg>
				Empty Trash
			</Button>
		{/if}
	</header>
	
	<Card padding="none">
		{#if loading}
			<div class="loading">Loading...</div>
		{:else if deletedLinks.length === 0}
			<div class="empty-state">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
					<polyline points="3 6 5 6 21 6"/>
					<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
				</svg>
				<p>Trash is empty</p>
			</div>
		{:else}
			<div class="links-list">
				{#each deletedLinks as link}
					<div class="trash-item">
						<div class="trash-link-info">
							<span class="link-slug">/{link.slug}</span>
							<span class="link-url">{link.original_url}</span>
							{#if link.deleted_at}
								<span class="deleted-date">Deleted {new Date(link.deleted_at).toLocaleDateString()}</span>
							{/if}
						</div>
						<div class="trash-actions">
							<Button variant="secondary" size="sm" onclick={() => handleRestore(link.slug)}>
								<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
									<polyline points="1 4 1 10 7 10"/>
									<path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/>
								</svg>
								Restore
							</Button>
							<Button variant="danger" size="sm" onclick={() => handlePermanentDelete(link.slug)}>
								<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
									<line x1="18" y1="6" x2="6" y2="18"/>
									<line x1="6" y1="6" x2="18" y2="18"/>
								</svg>
								Delete
							</Button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</Card>
</div>

<style>
	.trash-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}
	
	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-4);
		flex-wrap: wrap;
	}
	
	.page-title {
		font-size: var(--text-3xl);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.page-subtitle {
		font-size: var(--text-base);
		color: var(--text-tertiary);
		margin-top: var(--space-1);
	}
	
	.links-list {
		max-height: 600px;
		overflow-y: auto;
	}
	
	.loading, .empty-state {
		padding: var(--space-12);
		text-align: center;
		color: var(--text-muted);
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-4);
	}
	
	.empty-state svg {
		color: var(--text-muted);
	}
	
	.trash-item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-4);
		padding: var(--space-4) var(--space-5);
		border-bottom: 1px solid var(--border-light);
	}
	
	.trash-item:last-child {
		border-bottom: none;
	}
	
	.trash-link-info {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
		min-width: 0;
		flex: 1;
	}
	
	.link-slug {
		font-weight: var(--font-medium);
		color: var(--text-primary);
	}
	
	.link-url {
		font-size: var(--text-sm);
		color: var(--text-secondary);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.deleted-date {
		font-size: var(--text-xs);
		color: var(--text-muted);
	}
	
	.trash-actions {
		display: flex;
		gap: var(--space-2);
		flex-shrink: 0;
	}
	
	@media (max-width: 640px) {
		.trash-item {
			flex-direction: column;
			align-items: flex-start;
		}
		
		.trash-actions {
			width: 100%;
			justify-content: flex-end;
		}
	}
</style>
