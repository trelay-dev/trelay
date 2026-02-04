<script lang="ts">
	import { Button, Input, Card, Modal, LinkRow } from '$lib/components';
	import { folders, links, type Folder, type Link } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	let folderList = $state<Folder[]>([]);
	let loading = $state(true);
	let selectedFolder = $state<Folder | null>(null);
	let folderLinks = $state<Link[]>([]);
	let linksLoading = $state(false);
	
	let showCreateModal = $state(false);
	let createLoading = $state(false);
	let createError = $state('');
	let newName = $state('');
	
	onMount(async () => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		
		await loadFolders();
	});
	
	async function loadFolders() {
		loading = true;
		try {
			const res = await folders.list();
			if (res.success && res.data) {
				folderList = res.data;
			}
		} catch (e) {
			console.error('Failed to load folders:', e);
		} finally {
			loading = false;
		}
	}
	
	async function selectFolder(folder: Folder) {
		selectedFolder = folder;
		linksLoading = true;
		try {
			const res = await links.list({ folder_id: folder.id });
			if (res.success && res.data) {
				folderLinks = res.data;
			}
		} catch (e) {
			console.error('Failed to load folder links:', e);
		} finally {
			linksLoading = false;
		}
	}
	
	function clearSelection() {
		selectedFolder = null;
		folderLinks = [];
	}
	
	async function handleCreateFolder() {
		if (!newName.trim()) {
			createError = 'Name is required';
			return;
		}
		
		createLoading = true;
		createError = '';
		
		try {
			const res = await folders.create(newName);
			
			if (res.success) {
				showCreateModal = false;
				newName = '';
				await loadFolders();
			} else {
				createError = res.error?.message || 'Failed to create folder';
			}
		} catch (e) {
			createError = 'Failed to create folder';
		} finally {
			createLoading = false;
		}
	}
	
	async function handleDeleteFolder(id: number) {
		if (!confirm('Delete this folder?')) return;
		
		try {
			const res = await folders.delete(id);
			if (res.success) {
				if (selectedFolder?.id === id) {
					clearSelection();
				}
				folderList = folderList.filter(f => f.id !== id);
			}
		} catch (e) {
			console.error('Failed to delete folder:', e);
		}
	}
	
	async function handleDeleteLink(slug: string) {
		if (!confirm('Delete this link?')) return;
		
		try {
			const res = await links.delete(slug);
			if (res.success) {
				folderLinks = folderLinks.filter(l => l.slug !== slug);
			}
		} catch (e) {
			console.error('Failed to delete link:', e);
		}
	}
	
	function formatDate(date: string) {
		return new Date(date).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Folders - Trelay</title>
</svelte:head>

<div class="folders-page container">
	<header class="page-header">
		<div class="page-header-content">
			{#if selectedFolder}
				<button class="back-btn" onclick={clearSelection}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<line x1="19" y1="12" x2="5" y2="12"/>
						<polyline points="12 19 5 12 12 5"/>
					</svg>
					Back
				</button>
				<h1 class="page-title">{selectedFolder.name}</h1>
				<p class="page-subtitle">{folderLinks.length} links in this folder</p>
			{:else}
				<h1 class="page-title">Folders</h1>
				<p class="page-subtitle">{folderList.length} folders</p>
			{/if}
		</div>
		{#if !selectedFolder}
			<Button onclick={() => showCreateModal = true}>
				<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<line x1="12" y1="5" x2="12" y2="19"/>
					<line x1="5" y1="12" x2="19" y2="12"/>
				</svg>
				New Folder
			</Button>
		{/if}
	</header>
	
	{#if loading}
		<div class="loading">Loading...</div>
	{:else if selectedFolder}
		<Card padding="none">
			{#if linksLoading}
				<div class="loading">Loading links...</div>
			{:else if folderLinks.length === 0}
				<div class="empty-state">
					<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
						<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
						<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
					</svg>
					<p>No links in this folder</p>
				</div>
			{:else}
				<div class="links-list">
					{#each folderLinks as link (link.id)}
						<LinkRow {link} ondelete={handleDeleteLink} />
					{/each}
				</div>
			{/if}
		</Card>
	{:else if folderList.length === 0}
		<Card>
			<div class="empty-state">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
					<path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
				</svg>
				<p>No folders yet</p>
				<Button variant="secondary" onclick={() => showCreateModal = true}>Create your first folder</Button>
			</div>
		</Card>
	{:else}
		<div class="folders-grid">
			{#each folderList as folder (folder.id)}
				<div class="folder-card-wrapper">
					<button class="folder-card" onclick={() => selectFolder(folder)}>
						<div class="folder-icon">
							<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
							</svg>
						</div>
						<div class="folder-info">
							<h3 class="folder-name">{folder.name}</h3>
							<span class="folder-date">{formatDate(folder.created_at)}</span>
						</div>
					</button>
					<button class="folder-delete" onclick={() => handleDeleteFolder(folder.id)} title="Delete">
						<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<polyline points="3 6 5 6 21 6"/>
							<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
						</svg>
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>

<Modal
	open={showCreateModal}
	title="Create Folder"
	onclose={() => { showCreateModal = false; createError = ''; }}
>
	<form class="create-form" onsubmit={(e) => { e.preventDefault(); handleCreateFolder(); }}>
		<Input
			type="text"
			label="Name"
			placeholder="My Folder"
			bind:value={newName}
			error={createError}
		/>
		<div class="form-actions">
			<Button variant="secondary" onclick={() => showCreateModal = false}>Cancel</Button>
			<Button type="submit" loading={createLoading}>Create</Button>
		</div>
	</form>
</Modal>

<style>
	.folders-page {
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
	
	.folders-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: var(--space-4);
	}
	
	.folder-card-wrapper {
		display: flex;
		align-items: center;
		padding: var(--space-4);
		background: var(--card-bg);
		border: 1px solid var(--card-border);
		border-radius: var(--radius-lg);
		transition: box-shadow var(--transition-fast), border-color var(--transition-fast);
	}
	
	.folder-card-wrapper:hover {
		border-color: var(--border-medium);
		box-shadow: var(--shadow-md);
	}
	
	.back-btn {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-1) 0;
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-tertiary);
		background: none;
		border: none;
		cursor: pointer;
		transition: color var(--transition-fast);
	}
	
	.back-btn:hover {
		color: var(--accent-primary);
	}
	
	.links-list {
		max-height: 500px;
		overflow-y: auto;
	}
	
	.folder-card {
		display: flex;
		align-items: center;
		gap: var(--space-4);
		flex: 1;
		background: none;
		border: none;
		cursor: pointer;
		text-align: left;
	}
	
	.folder-icon {
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
	
	.folder-info {
		flex: 1;
		min-width: 0;
	}
	
	.folder-name {
		font-size: var(--text-base);
		font-weight: var(--font-medium);
		color: var(--text-primary);
		margin-bottom: var(--space-1);
	}
	
	.folder-date {
		font-size: var(--text-sm);
		color: var(--text-muted);
	}
	
	.folder-delete {
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
	
	.folder-delete:hover {
		background: rgba(239, 68, 68, 0.1);
		color: var(--color-error);
	}
	
	.loading {
		padding: var(--space-12);
		text-align: center;
		color: var(--text-muted);
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-4);
		padding: var(--space-8);
		text-align: center;
		color: var(--text-muted);
	}
	
	.empty-state svg {
		color: var(--text-muted);
	}
	
	.create-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-3);
		margin-top: var(--space-2);
	}
</style>
