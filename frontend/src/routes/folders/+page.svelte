<script lang="ts">
	import { Button, Input, Card, Modal } from '$lib/components';
	import { folders, type Folder } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	let folderList = $state<Folder[]>([]);
	let loading = $state(true);
	
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
		
		const res = await folders.delete(id);
		if (res.success) {
			await loadFolders();
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
			<h1 class="page-title">Folders</h1>
			<p class="page-subtitle">{folderList.length} folders</p>
		</div>
		<Button onclick={() => showCreateModal = true}>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<line x1="12" y1="5" x2="12" y2="19"/>
				<line x1="5" y1="12" x2="19" y2="12"/>
			</svg>
			New Folder
		</Button>
	</header>
	
	{#if loading}
		<div class="loading">Loading...</div>
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
			{#each folderList as folder}
				<Card hoverable>
					<div class="folder-card">
						<div class="folder-icon">
							<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
							</svg>
						</div>
						<div class="folder-info">
							<h3 class="folder-name">{folder.name}</h3>
							<span class="folder-date">{formatDate(folder.created_at)}</span>
						</div>
						<button class="folder-delete" onclick={() => handleDeleteFolder(folder.id)} title="Delete">
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<polyline points="3 6 5 6 21 6"/>
								<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
							</svg>
						</button>
					</div>
				</Card>
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
	
	.folder-card {
		display: flex;
		align-items: center;
		gap: var(--space-4);
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
