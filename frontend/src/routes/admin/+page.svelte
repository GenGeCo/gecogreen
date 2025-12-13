<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { isAuthenticated, isAdmin, currentUser } from '$lib/stores/auth';

	interface ImageReview {
		id: string;
		user_id: string;
		product_id?: string;
		image_url: string;
		image_type: string;
		detected_text: string;
		detected_phone: boolean;
		detected_email: boolean;
		detected_url: boolean;
		confidence_score: number;
		status: string;
		created_at: string;
	}

	let reviews: ImageReview[] = [];
	let stats = { PENDING: 0, APPROVED: 0, REJECTED: 0 };
	let loading = true;
	let error = '';
	let success = '';
	let page = 1;
	let totalPages = 1;

	// Redirect if not admin
	$: if ($isAuthenticated !== undefined && !$isAuthenticated) {
		goto('/login');
	}
	$: if ($isAdmin !== undefined && !$isAdmin) {
		goto('/');
	}

	async function loadReviews() {
		loading = true;
		error = '';
		try {
			const token = localStorage.getItem('access_token');
			const response = await fetch(`/api/v1/admin/reviews?page=${page}&per_page=20`, {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data = await response.json();
			if (!response.ok) throw new Error(data.error);
			reviews = data.reviews || [];
			totalPages = data.total_pages || 1;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore caricamento';
		}
		loading = false;
	}

	async function loadStats() {
		try {
			const token = localStorage.getItem('access_token');
			const response = await fetch('/api/v1/admin/reviews/stats', {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data = await response.json();
			if (response.ok) {
				stats = data;
			}
		} catch (e) {
			console.error('Error loading stats:', e);
		}
	}

	async function approveReview(id: string) {
		try {
			const token = localStorage.getItem('access_token');
			const response = await fetch(`/api/v1/admin/reviews/${id}/approve`, {
				method: 'POST',
				headers: { Authorization: `Bearer ${token}` }
			});
			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error);
			}
			reviews = reviews.filter(r => r.id !== id);
			stats.PENDING--;
			stats.APPROVED++;
			success = 'Immagine approvata';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore approvazione';
		}
	}

	async function rejectReview(id: string) {
		const reason = prompt('Motivo del rifiuto (opzionale):');
		try {
			const token = localStorage.getItem('access_token');
			const response = await fetch(`/api/v1/admin/reviews/${id}/reject`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ reason: reason || '' })
			});
			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error);
			}
			reviews = reviews.filter(r => r.id !== id);
			stats.PENDING--;
			stats.REJECTED++;
			success = 'Immagine rifiutata';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore rifiuto';
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString('it-IT');
	}

	onMount(() => {
		if ($isAdmin) {
			loadReviews();
			loadStats();
		}
	});

	$: if ($isAdmin) {
		loadReviews();
		loadStats();
	}
</script>

<svelte:head>
	<title>Admin - Moderazione Immagini - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<h1 class="text-3xl font-bold mb-6">Moderazione Immagini</h1>

	{#if error}
		<div class="alert alert-error mb-6">
			<span>{error}</span>
			<button class="btn btn-sm btn-ghost" on:click={() => error = ''}>✕</button>
		</div>
	{/if}

	{#if success}
		<div class="alert alert-success mb-6">
			<span>{success}</span>
			<button class="btn btn-sm btn-ghost" on:click={() => success = ''}>✕</button>
		</div>
	{/if}

	<!-- Stats -->
	<div class="grid grid-cols-3 gap-4 mb-8">
		<div class="stat bg-warning/20 rounded-lg">
			<div class="stat-title">In Attesa</div>
			<div class="stat-value text-warning">{stats.PENDING}</div>
		</div>
		<div class="stat bg-success/20 rounded-lg">
			<div class="stat-title">Approvate</div>
			<div class="stat-value text-success">{stats.APPROVED}</div>
		</div>
		<div class="stat bg-error/20 rounded-lg">
			<div class="stat-title">Rifiutate</div>
			<div class="stat-value text-error">{stats.REJECTED}</div>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if reviews.length === 0}
		<div class="text-center py-12 bg-base-100 rounded-lg shadow">
			<p class="text-lg text-base-content/70">Nessuna immagine da moderare</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each reviews as review}
				<div class="card bg-base-100 shadow">
					<figure class="relative">
						<img
							src={review.image_url}
							alt="Immagine da moderare"
							class="w-full h-48 object-cover"
						/>
						<div class="absolute top-2 right-2">
							<span class="badge badge-{review.image_type === 'PRODUCT' ? 'primary' : 'secondary'}">
								{review.image_type}
							</span>
						</div>
					</figure>
					<div class="card-body">
						<!-- Detection Info -->
						<div class="space-y-2 text-sm">
							{#if review.detected_phone}
								<div class="flex items-center gap-2 text-error">
									<span>📞</span> Telefono rilevato
								</div>
							{/if}
							{#if review.detected_email}
								<div class="flex items-center gap-2 text-error">
									<span>📧</span> Email rilevata
								</div>
							{/if}
							{#if review.detected_url}
								<div class="flex items-center gap-2 text-error">
									<span>🔗</span> URL rilevato
								</div>
							{/if}
							{#if review.detected_text}
								<details class="collapse collapse-arrow bg-base-200 rounded-lg">
									<summary class="collapse-title text-sm font-medium py-2 min-h-0">
										Testo rilevato
									</summary>
									<div class="collapse-content">
										<pre class="text-xs whitespace-pre-wrap">{review.detected_text}</pre>
									</div>
								</details>
							{/if}
						</div>

						<div class="text-xs text-base-content/60 mt-2">
							{formatDate(review.created_at)}
						</div>

						<div class="card-actions justify-end mt-4">
							<button
								class="btn btn-error btn-sm"
								on:click={() => rejectReview(review.id)}
							>
								Rifiuta
							</button>
							<button
								class="btn btn-success btn-sm"
								on:click={() => approveReview(review.id)}
							>
								Approva
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="flex justify-center gap-2 mt-8">
				<button
					class="btn btn-sm"
					disabled={page <= 1}
					on:click={() => { page--; loadReviews(); }}
				>
					« Precedente
				</button>
				<span class="btn btn-sm btn-ghost">Pagina {page} di {totalPages}</span>
				<button
					class="btn btn-sm"
					disabled={page >= totalPages}
					on:click={() => { page++; loadReviews(); }}
				>
					Successiva »
				</button>
			</div>
		{/if}
	{/if}
</div>
