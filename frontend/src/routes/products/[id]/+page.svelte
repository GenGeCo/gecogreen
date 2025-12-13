<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api, type Product } from '$lib/api';
	import { isAuthenticated, currentUser } from '$lib/stores/auth';

	let product: Product | null = null;
	let loading = true;
	let error = '';
	let selectedImage = 0;
	let quantity = 1;

	$: productId = $page.params.id;

	async function loadProduct() {
		loading = true;
		error = '';
		try {
			product = await api.getProduct(productId);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel caricamento';
		}
		loading = false;
	}

	function formatPrice(price: number): string {
		return new Intl.NumberFormat('it-IT', {
			style: 'currency',
			currency: 'EUR'
		}).format(price);
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('it-IT', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		});
	}

	function getDiscountPercent(original: number, current: number): number {
		return Math.round(((original - current) / original) * 100);
	}

	function getShippingLabel(method: string): string {
		switch (method) {
			case 'PICKUP': return 'Ritiro in sede';
			case 'SELLER_SHIPS': return 'Spedizione a carico venditore';
			case 'BUYER_ARRANGES': return 'Spedizione a carico acquirente';
			case 'PLATFORM_MANAGED': return 'Gestito dalla piattaforma';
			default: return method;
		}
	}

	function handleAddToCart() {
		// TODO: Implement cart functionality
		alert(`Aggiunto ${quantity} x ${product?.title} al carrello`);
	}

	onMount(loadProduct);
</script>

<svelte:head>
	<title>{product?.title || 'Prodotto'} - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<span>{error}</span>
		</div>
	{:else if product}
		<div class="breadcrumbs text-sm mb-6">
			<ul>
				<li><a href="/">Home</a></li>
				<li><a href="/products">Prodotti</a></li>
				<li>{product.title}</li>
			</ul>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Images -->
			<div class="space-y-4">
				<div class="aspect-square bg-base-200 rounded-lg overflow-hidden">
					{#if product.images && product.images.length > 0}
						<img
							src={product.images[selectedImage]}
							alt={product.title}
							class="w-full h-full object-cover"
						/>
					{:else}
						<div class="flex items-center justify-center h-full text-base-content/30">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
							</svg>
						</div>
					{/if}
				</div>

				{#if product.images && product.images.length > 1}
					<div class="flex gap-2 overflow-x-auto">
						{#each product.images as img, i}
							<button
								class="w-20 h-20 rounded-lg overflow-hidden border-2 flex-shrink-0 transition-colors"
								class:border-primary={i === selectedImage}
								class:border-transparent={i !== selectedImage}
								on:click={() => selectedImage = i}
							>
								<img src={img} alt="Thumbnail {i + 1}" class="w-full h-full object-cover" />
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Product Info -->
			<div class="space-y-6">
				<div>
					{#if product.listing_type === 'GIFT'}
						<div class="badge badge-success badge-lg mb-2">Regalo</div>
					{/if}
					<h1 class="text-3xl font-bold">{product.title}</h1>
				</div>

				<!-- Price -->
				<div class="bg-base-200 p-4 rounded-lg">
					{#if product.listing_type === 'GIFT'}
						<span class="text-3xl font-bold text-success">Gratis</span>
					{:else}
						<div class="flex items-baseline gap-3">
							<span class="text-3xl font-bold text-primary">{formatPrice(product.price)}</span>
							{#if product.original_price && product.original_price > product.price}
								<span class="text-xl line-through text-base-content/50">
									{formatPrice(product.original_price)}
								</span>
								<span class="badge badge-error badge-lg">
									-{getDiscountPercent(product.original_price, product.price)}%
								</span>
							{/if}
						</div>
					{/if}
					{#if product.shipping_cost > 0}
						<p class="text-sm text-base-content/70 mt-2">
							+ {formatPrice(product.shipping_cost)} spedizione
						</p>
					{/if}
				</div>

				<!-- Stock & Expiry -->
				<div class="flex gap-4">
					{#if product.quantity_available > 0}
						<div class="badge badge-outline badge-lg">
							{product.quantity_available} disponibili
						</div>
					{:else}
						<div class="badge badge-error badge-lg">Esaurito</div>
					{/if}
					{#if product.expiry_date}
						<div class="badge badge-warning badge-lg">
							Scade: {formatDate(product.expiry_date)}
						</div>
					{/if}
				</div>

				<!-- Description -->
				<div>
					<h2 class="text-lg font-semibold mb-2">Descrizione</h2>
					<p class="text-base-content/80 whitespace-pre-wrap">{product.description}</p>
				</div>

				<!-- Details -->
				<div class="divider"></div>
				<div class="grid grid-cols-2 gap-4 text-sm">
					<div>
						<span class="text-base-content/60">Spedizione</span>
						<p class="font-medium">{getShippingLabel(product.shipping_method)}</p>
					</div>
					<div>
						<span class="text-base-content/60">Località</span>
						<p class="font-medium">{product.city}, {product.province}</p>
					</div>
					<div>
						<span class="text-base-content/60">Venditore</span>
						<p class="font-medium">
							{#if product.seller}
								{product.seller.first_name} {product.seller.last_name}
							{:else}
								Venditore
							{/if}
						</p>
					</div>
					<div>
						<span class="text-base-content/60">Pubblicato</span>
						<p class="font-medium">{formatDate(product.created_at)}</p>
					</div>
				</div>

				<!-- Actions -->
				{#if product.quantity_available > 0}
					<div class="divider"></div>
					<div class="flex gap-4 items-center">
						<div class="form-control w-24">
							<label class="label">
								<span class="label-text">Quantità</span>
							</label>
							<input
								type="number"
								min="1"
								max={product.quantity_available}
								bind:value={quantity}
								class="input input-bordered"
							/>
						</div>
						<div class="flex-1">
							{#if $isAuthenticated}
								{#if $currentUser?.id === product.seller_id}
									<a href="/seller/products/{product.id}/edit" class="btn btn-primary w-full">
										Modifica Prodotto
									</a>
								{:else}
									<button class="btn btn-primary btn-lg w-full" on:click={handleAddToCart}>
										Aggiungi al Carrello
									</button>
								{/if}
							{:else}
								<a href="/login" class="btn btn-primary btn-lg w-full">
									Accedi per acquistare
								</a>
							{/if}
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
