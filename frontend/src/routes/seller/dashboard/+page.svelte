<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api, type Product } from '$lib/api';
	import { isAuthenticated } from '$lib/stores/auth';

	let products: Product[] = [];
	let loading = true;
	let error = '';
	let stats = {
		total: 0,
		active: 0,
		sold: 0,
		draft: 0
	};

	// Redirect if not authenticated
	$: if ($isAuthenticated !== undefined && !$isAuthenticated) {
		goto('/login');
	}

	async function loadProducts() {
		loading = true;
		error = '';
		try {
			const result = await api.getMyProducts({ per_page: 100 });
			products = result.products || [];

			// Calculate stats
			stats.total = products.length;
			stats.active = products.filter(p => p.status === 'ACTIVE').length;
			stats.sold = products.filter(p => p.status === 'SOLD').length;
			stats.draft = products.filter(p => p.status === 'DRAFT').length;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel caricamento';
		}
		loading = false;
	}

	async function deleteProduct(id: string, title: string) {
		if (!confirm(`Sei sicuro di voler eliminare "${title}"?`)) return;

		try {
			await api.deleteProduct(id);
			products = products.filter(p => p.id !== id);
			stats.total--;
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Errore eliminazione');
		}
	}

	function formatPrice(price: number): string {
		return new Intl.NumberFormat('it-IT', {
			style: 'currency',
			currency: 'EUR'
		}).format(price);
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('it-IT');
	}

	function getStatusBadge(status: string): { class: string; text: string } {
		switch (status) {
			case 'ACTIVE': return { class: 'badge-success', text: 'Attivo' };
			case 'SOLD': return { class: 'badge-info', text: 'Venduto' };
			case 'DRAFT': return { class: 'badge-warning', text: 'Bozza' };
			case 'EXPIRED': return { class: 'badge-error', text: 'Scaduto' };
			default: return { class: 'badge-ghost', text: status };
		}
	}

	onMount(loadProducts);
</script>

<svelte:head>
	<title>Dashboard Venditore - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-bold">Dashboard Venditore</h1>
		<a href="/seller/products/new" class="btn btn-primary">
			+ Nuovo Prodotto
		</a>
	</div>

	<!-- Stats -->
	<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Totale Prodotti</div>
			<div class="stat-value text-primary">{stats.total}</div>
		</div>
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Attivi</div>
			<div class="stat-value text-success">{stats.active}</div>
		</div>
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Venduti</div>
			<div class="stat-value text-info">{stats.sold}</div>
		</div>
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Bozze</div>
			<div class="stat-value text-warning">{stats.draft}</div>
		</div>
	</div>

	{#if error}
		<div class="alert alert-error mb-6">
			<span>{error}</span>
		</div>
	{/if}

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if products.length === 0}
		<div class="text-center py-12 bg-base-100 rounded-lg shadow">
			<p class="text-lg text-base-content/70 mb-4">Non hai ancora prodotti</p>
			<a href="/seller/products/new" class="btn btn-primary">
				Pubblica il tuo primo prodotto
			</a>
		</div>
	{:else}
		<div class="overflow-x-auto bg-base-100 rounded-lg shadow">
			<table class="table">
				<thead>
					<tr>
						<th>Prodotto</th>
						<th>Prezzo</th>
						<th>Disponibilità</th>
						<th>Stato</th>
						<th>Data</th>
						<th>Azioni</th>
					</tr>
				</thead>
				<tbody>
					{#each products as product}
						<tr>
							<td>
								<div class="flex items-center gap-3">
									<div class="avatar">
										<div class="w-12 h-12 rounded bg-base-200">
											{#if product.images && product.images.length > 0}
												<img src={product.images[0]} alt={product.title} />
											{/if}
										</div>
									</div>
									<div>
										<div class="font-bold line-clamp-1">{product.title}</div>
										<div class="text-sm text-base-content/60">{product.city}</div>
									</div>
								</div>
							</td>
							<td>
								{#if product.listing_type === 'GIFT'}
									<span class="text-success">Gratis</span>
								{:else}
									{formatPrice(product.price)}
								{/if}
							</td>
							<td>{product.quantity_available} / {product.quantity}</td>
							<td>
								<span class="badge {getStatusBadge(product.status).class}">
									{getStatusBadge(product.status).text}
								</span>
							</td>
							<td>{formatDate(product.created_at)}</td>
							<td>
								<div class="flex gap-2">
									<a href="/products/{product.id}" class="btn btn-ghost btn-xs">
										Vedi
									</a>
									<a href="/seller/products/{product.id}/edit" class="btn btn-ghost btn-xs">
										Modifica
									</a>
									<button
										class="btn btn-ghost btn-xs text-error"
										on:click={() => deleteProduct(product.id, product.title)}
									>
										Elimina
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
