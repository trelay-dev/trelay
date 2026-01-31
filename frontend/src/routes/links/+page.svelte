<script lang="ts">
	import { Button, Input, Card, LinkRow, Modal } from '$lib/components';
	import { links, type Link, type CreateLinkRequest } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	let linkList = $state<Link[]>([]);
	let loading = $state(true);
	let search = $state('');
	
	let showCreateModal = $state(false);
	let createLoading = $state(false);
	let createError = $state('');
	let newUrl = $state('');
	let newSlug = $state('');
	let newPassword = $state('');
	let newTtl = $state('');
	let newTags = $state('');
	let isOneTime = $state(false);
	
	onMount(async () => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		
		await loadLinks();
	});
	
	async function loadLinks() {
		loading = true;
		try {
			const res = await links.list({ search: search || undefined });
			if (res.success && res.data) {
				linkList = res.data;
			}
		} catch (e) {
			console.error('Failed to load links:', e);
		} finally {
			loading = false;
		}
	}
	
	async function handleSearch() {
		await loadLinks();
	}
	
	async function handleCreateLink() {
		if (!newUrl.trim()) {
			createError = 'URL is required';
			return;
		}
		
		createLoading = true;
		createError = '';
		
		try {
			const req: CreateLinkRequest = { url: newUrl };
			if (newSlug.trim()) req.slug = newSlug;
			if (newPassword.trim()) req.password = newPassword;
			if (newTtl && parseInt(newTtl) > 0) req.ttl_hours = parseInt(newTtl);
			if (newTags.trim()) req.tags = newTags.split(',').map(t => t.trim()).filter(Boolean);
			if (isOneTime) req.is_one_time = true;
			
			const res = await links.create(req);
			
			if (res.success) {
				resetCreateForm();
				showCreateModal = false;
				await loadLinks();
			} else {
				createError = res.error?.message || 'Failed to create link';
			}
		} catch (e) {
			createError = 'Failed to create link';
		} finally {
			createLoading = false;
		}
	}
	
	function resetCreateForm() {
		newUrl = '';
		newSlug = '';
		newPassword = '';
		newTtl = '';
		newTags = '';
		isOneTime = false;
		createError = '';
	}
	
	async function handleDeleteLink(slug: string) {
		if (!confirm('Delete this link?')) return;
		
		const res = await links.delete(slug);
		if (res.success) {
			await loadLinks();
		}
	}
</script>

<svelte:head>
	<title>Links - Trelay</title>
</svelte:head>

<div class="links-page container">
	<header class="page-header">
		<div class="page-header-content">
			<h1 class="page-title">Links</h1>
			<p class="page-subtitle">{linkList.length} total links</p>
		</div>
		<Button onclick={() => showCreateModal = true}>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<line x1="12" y1="5" x2="12" y2="19"/>
				<line x1="5" y1="12" x2="19" y2="12"/>
			</svg>
			New Link
		</Button>
	</header>
	
	<div class="search-bar">
		<Input
			type="search"
			placeholder="Search links..."
			bind:value={search}
			onkeydown={(e) => e.key === 'Enter' && handleSearch()}
		/>
		<Button variant="secondary" onclick={handleSearch}>Search</Button>
	</div>
	
	<Card padding="none">
		{#if loading}
			<div class="loading">Loading...</div>
		{:else if linkList.length === 0}
			<div class="empty-state">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
					<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
					<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
				</svg>
				<p>No links found</p>
				<Button variant="secondary" onclick={() => showCreateModal = true}>Create your first link</Button>
			</div>
		{:else}
			<div class="links-list">
				{#each linkList as link}
					<LinkRow {link} ondelete={handleDeleteLink} />
				{/each}
			</div>
		{/if}
	</Card>
</div>

<Modal
	open={showCreateModal}
	title="Create Link"
	onclose={() => { showCreateModal = false; resetCreateForm(); }}
>
	<form class="create-form" onsubmit={(e) => { e.preventDefault(); handleCreateLink(); }}>
		<Input
			type="url"
			label="URL"
			placeholder="https://example.com"
			bind:value={newUrl}
			error={createError}
		/>
		<Input
			type="text"
			label="Custom Slug (optional)"
			placeholder="my-link"
			bind:value={newSlug}
		/>
		<Input
			type="password"
			label="Password (optional)"
			placeholder="Password protect this link"
			bind:value={newPassword}
		/>
		<Input
			type="number"
			label="Expiry (hours, optional)"
			placeholder="24"
			bind:value={newTtl}
		/>
		<Input
			type="text"
			label="Tags (comma-separated)"
			placeholder="project, docs"
			bind:value={newTags}
		/>
		<label class="checkbox-label">
			<input type="checkbox" bind:checked={isOneTime} />
			<span>One-time link (burns after first access)</span>
		</label>
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showCreateModal = false; resetCreateForm(); }}>Cancel</Button>
			<Button type="submit" loading={createLoading}>Create</Button>
		</div>
	</form>
</Modal>

<style>
	.links-page {
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
	
	.search-bar {
		display: flex;
		gap: var(--space-3);
	}
	
	.search-bar :global(.input-wrapper) {
		flex: 1;
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
	
	.create-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-sm);
		color: var(--text-secondary);
		cursor: pointer;
	}
	
	.checkbox-label input {
		width: 16px;
		height: 16px;
		accent-color: var(--accent-primary);
	}
	
	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-3);
		margin-top: var(--space-2);
	}
</style>
