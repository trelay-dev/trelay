<script lang="ts">
	import { Button, Input, Card, LinkRow, Modal, QRModal } from '$lib/components';
	import { links, folders, type Link, type CreateLinkRequest, type Folder } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	
	let linkList = $state<Link[]>([]);
	let folderList = $state<Folder[]>([]);
	let loading = $state(true);
	let search = $state('');
	let searchTimeout: ReturnType<typeof setTimeout>;
	
	// Create modal
	let showCreateModal = $state(false);
	let createLoading = $state(false);
	let createError = $state('');
	let newUrl = $state('');
	let newSlug = $state('');
	let newPassword = $state('');
	let newTtl = $state('');
	let newTags = $state('');
	let newFolderId = $state('');
	let isOneTime = $state(false);
	
	// Edit modal
	let showEditModal = $state(false);
	let editLink = $state<Link | null>(null);
	let editLoading = $state(false);
	let editError = $state('');
	let editUrl = $state('');
	let editTags = $state('');
	let editFolderId = $state('');
	
	// QR modal
	let showQRModal = $state(false);
	let qrSlug = $state('');
	
	// Stats modal
	let showStatsModal = $state(false);
	let statsLink = $state<Link | null>(null);
	let statsData = $state<{ total_clicks: number; clicks_by_day?: { date: string; clicks: number }[] } | null>(null);
	let statsLoading = $state(false);
	
	let baseUrl = $derived(browser ? window.location.origin : '');
	
	onMount(async () => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		
		await Promise.all([loadLinks(), loadFolders()]);
	});
	
	async function loadFolders() {
		try {
			const res = await folders.list();
			if (res.success && res.data) {
				folderList = res.data;
			}
		} catch (e) {
			console.error('Failed to load folders:', e);
		}
	}
	
	// Reactive search with debounce
	$effect(() => {
		const term = search;
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			loadLinks();
		}, 300);
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
			if (newFolderId) req.folder_id = parseInt(newFolderId);
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
		newFolderId = '';
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
	
	function handleQR(link: Link) {
		qrSlug = link.slug;
		showQRModal = true;
	}
	
	function handleEdit(link: Link) {
		editLink = link;
		editUrl = link.original_url;
		editTags = link.tags?.join(', ') || '';
		editFolderId = link.folder_id ? String(link.folder_id) : '';
		editError = '';
		showEditModal = true;
	}
	
	async function saveEdit() {
		if (!editLink) return;
		
		editLoading = true;
		editError = '';
		
		try {
			const req: Partial<CreateLinkRequest> = {};
			if (editUrl !== editLink.original_url) req.url = editUrl;
			if (editTags.trim()) {
				req.tags = editTags.split(',').map(t => t.trim()).filter(Boolean);
			}
			const newFolderIdNum = editFolderId ? parseInt(editFolderId) : null;
			const currentFolderId = editLink.folder_id || null;
			if (newFolderIdNum !== currentFolderId) {
				req.folder_id = newFolderIdNum as number | undefined;
			}
			
			const res = await links.update(editLink.slug, req);
			
			if (res.success) {
				showEditModal = false;
				editLink = null;
				await loadLinks();
			} else {
				editError = res.error?.message || 'Failed to update link';
			}
		} catch (e) {
			editError = 'Failed to update link';
		} finally {
			editLoading = false;
		}
	}
	
	async function handleStats(link: Link) {
		statsLink = link;
		statsData = null;
		statsLoading = true;
		showStatsModal = true;
		
		try {
			const res = await fetch(`/api/v1/stats/${link.slug}`, {
				headers: { 'X-API-Key': localStorage.getItem('trelay-api-key') || '' }
			});
			const data = await res.json();
			if (data.success) {
				statsData = data.data;
			}
		} catch (e) {
			console.error('Failed to load stats:', e);
		} finally {
			statsLoading = false;
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
		/>
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
					<LinkRow 
						{link} 
						ondelete={handleDeleteLink} 
						onedit={handleEdit}
						onqr={handleQR}
						onstats={handleStats}
					/>
				{/each}
			</div>
		{/if}
	</Card>
</div>

<QRModal 
	open={showQRModal} 
	slug={qrSlug} 
	{baseUrl}
	onclose={() => showQRModal = false} 
/>

<Modal
	open={showEditModal}
	title="Edit Link"
	onclose={() => { showEditModal = false; editLink = null; }}
>
	<form class="create-form" onsubmit={(e) => { e.preventDefault(); saveEdit(); }}>
		<div class="edit-slug">
			<span class="edit-slug-label">Slug</span>
			<span class="edit-slug-value">/{editLink?.slug}</span>
		</div>
		<Input
			type="url"
			label="URL"
			placeholder="https://example.com"
			bind:value={editUrl}
			error={editError}
		/>
		<Input
			type="text"
			label="Tags (comma-separated)"
			placeholder="project, docs"
			bind:value={editTags}
		/>
		<div class="select-wrapper">
			<label class="select-label">Folder</label>
			<select class="select-input" bind:value={editFolderId}>
				<option value="">No folder</option>
				{#each folderList as folder}
					<option value={folder.id}>{folder.name}</option>
				{/each}
			</select>
		</div>
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showEditModal = false; editLink = null; }}>Cancel</Button>
			<Button type="submit" loading={editLoading}>Save</Button>
		</div>
	</form>
</Modal>

<Modal
	open={showStatsModal}
	title={statsLink ? `Stats: /${statsLink.slug}` : 'Stats'}
	onclose={() => { showStatsModal = false; statsLink = null; }}
>
	<div class="stats-content">
		{#if statsLoading}
			<div class="stats-loading">Loading stats...</div>
		{:else if statsData}
			<div class="stats-hero">
				<span class="stats-hero-value">{statsData.total_clicks.toLocaleString()}</span>
				<span class="stats-hero-label">Total Clicks</span>
			</div>
			{#if statsData.clicks_by_day && statsData.clicks_by_day.length > 0}
				<div class="stats-chart">
					<h4 class="stats-chart-title">Last 7 Days</h4>
					<div class="stats-bars">
						{#each statsData.clicks_by_day.slice(-7) as day}
							<div class="stats-bar-item">
								<div class="stats-bar" style="height: {Math.max(4, (day.clicks / Math.max(...statsData.clicks_by_day!.map(d => d.clicks))) * 80)}px"></div>
								<span class="stats-bar-label">{day.date.slice(5)}</span>
							</div>
						{/each}
					</div>
				</div>
			{/if}
		{:else}
			<div class="stats-loading">No stats available</div>
		{/if}
	</div>
</Modal>

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
		<div class="select-wrapper">
			<label class="select-label">Folder (optional)</label>
			<select class="select-input" bind:value={newFolderId}>
				<option value="">No folder</option>
				{#each folderList as folder}
					<option value={folder.id}>{folder.name}</option>
				{/each}
			</select>
		</div>
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
	
	.select-wrapper {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.select-label {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.select-input {
		width: 100%;
		padding: var(--space-2) var(--space-3);
		font-size: var(--text-base);
		color: var(--text-primary);
		background: var(--bg-primary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: border-color var(--transition-fast);
	}
	
	.select-input:hover {
		border-color: var(--border-medium);
	}
	
	.select-input:focus {
		outline: none;
		border-color: var(--accent-primary);
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
	
	.edit-slug {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.edit-slug-label {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.edit-slug-value {
		font-size: var(--text-base);
		color: var(--text-primary);
		padding: var(--space-2) var(--space-3);
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
	}
	
	.stats-content {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}
	
	.stats-loading {
		text-align: center;
		color: var(--text-muted);
		padding: var(--space-8);
	}
	
	.stats-hero {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-1);
	}
	
	.stats-hero-value {
		font-size: var(--text-hero);
		font-weight: var(--font-bold);
		color: var(--text-primary);
		line-height: 1;
	}
	
	.stats-hero-label {
		font-size: var(--text-sm);
		color: var(--text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.stats-chart {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}
	
	.stats-chart-title {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.stats-bars {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: var(--space-2);
		height: 100px;
		padding-bottom: var(--space-6);
	}
	
	.stats-bar-item {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-1);
	}
	
	.stats-bar {
		width: 100%;
		max-width: 40px;
		background: var(--accent-primary);
		border-radius: var(--radius-sm) var(--radius-sm) 0 0;
		transition: height var(--transition-slow);
	}
	
	.stats-bar-label {
		font-size: var(--text-xs);
		color: var(--text-muted);
	}
</style>
