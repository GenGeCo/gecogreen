<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Product, type Category } from '$lib/api';

	let products: Product[] = [];
	let categories: Category[] = [];
	let loading = true;
	let error = '';
	let totalPages = 1;
	let currentPage = 1;

	// Filters
	let search = '';
	let categoryId = '';
	let minPrice = '';
	let maxPrice = '';
	let city = '';
	let sortBy = 'created_at';
	let sortOrder = 'desc';

	async function loadCategories() {
		try {
			categories = await api.getCategories();
		} catch (e) {
			console.error('Errore caricamento categorie:', e);
		}
	}

	async function loadProducts() {
		loading = true;
		error = '';
		try {
			const result = await api.getProducts({
				page: currentPage,
				per_page: 12,
				search: search || undefined,
				category_id: categoryId || undefined,
				min_price: minPrice ? parseFloat(minPrice) : undefined,
				max_price: maxPrice ? parseFloat(maxPrice) : undefined,
				city: city || undefined,
				sort_by: sortBy,
				sort_order: sortOrder
			});
			products = result.products || [];
			totalPages = result.total_pages;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel caricamento';
		}
		loading = false;
	}

	function applyFilters() {
		currentPage = 1;
		loadProducts();
	}

	function clearFilters() {
		search = '';
		categoryId = '';
		minPrice = '';
		maxPrice = '';
		city = '';
		sortBy = 'created_at';
		sortOrder = 'desc';
		currentPage = 1;
		loadProducts();
	}

	function goToPage(page: number) {
		currentPage = page;
		loadProducts();
	}

	function formatPrice(price: number): string {
		return new Intl.NumberFormat('it-IT', {
			style: 'currency',
			currency: 'EUR'
		}).format(price);
	}

	function getDiscountPercent(original: number, current: number): number {
		return Math.round(((original - current) / original) * 100);
	}

	onMount(() => {
		loadCategories();
		loadProducts();
	});
</script>

<svelte:head>
	<title>Prodotti - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<h1 class="text-3xl font-bold mb-6">Prodotti Antispreco</h1>

	<!-- Filters -->
	<div class="card bg-base-100 shadow mb-6">
		<div class="card-body">
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4">
				<div class="form-control">
					<input
						type="text"
						placeholder="Cerca prodotti..."
						class="input input-bordered"
						bind:value={search}
						on:keypress={(e) => e.key === 'Enter' && applyFilters()}
					/>
				</div>
				<div class="form-control">
					<select class="select select-bordered" bind:value={categoryId} on:change={applyFilters}>
						<option value="">Tutte le categorie</option>
						{#each categories as cat}
							<option value={cat.id}>{cat.name}</option>
						{/each}
					</select>
				</div>
				<div class="form-control">
					<input
						type="number"
						placeholder="Prezzo min"
						class="input input-bordered"
						bind:value={minPrice}
					/>
				</div>
				<div class="form-control">
					<input
						type="number"
						placeholder="Prezzo max"
						class="input input-bordered"
						bind:value={maxPrice}
					/>
				</div>
				<div class="form-control">
					<input
						type="text"
						placeholder="Città"
						class="input input-bordered"
						bind:value={city}
					/>
				</div>
				<div class="form-control">
					<select class="select select-bordered" bind:value={sortBy} on:change={applyFilters}>
						<option value="created_at">Più recenti</option>
						<option value="price">Prezzo</option>
						<option value="expiry_date">Scadenza</option>
					</select>
				</div>
			</div>
			<div class="flex gap-2 mt-4">
				<button class="btn btn-primary" on:click={applyFilters}>Applica Filtri</button>
				<button class="btn btn-ghost" on:click={clearFilters}>Cancella</button>
			</div>
		</div>
	</div>

	<!-- Error -->
	{#if error}
		<div class="alert alert-error mb-6">
			<span>{error}</span>
		</div>
	{/if}

	<!-- Loading -->
	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if products.length === 0}
		<div class="text-center py-12">
			<p class="text-lg text-base-content/70">Nessun prodotto trovato</p>
			<button class="btn btn-primary mt-4" on:click={clearFilters}>Rimuovi filtri</button>
		</div>
	{:else}
		<!-- Products Grid -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
			{#each products as product}
				<a href="/products/{product.id}" class="card bg-base-100 shadow-lg hover:shadow-xl transition-shadow">
					<figure class="relative h-48 bg-base-200">
						{#if product.images && product.images.length > 0}
							<img
								src={product.images[0]}
								alt={product.title}
								class="w-full h-full object-cover"
							/>
						{:else}
							<div class="flex items-center justify-center h-full text-base-content/30">
								<svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
							</div>
						{/if}
						{#if product.original_price && product.original_price > product.price}
							<div class="absolute top-2 right-2 badge badge-error">
								-{getDiscountPercent(product.original_price, product.price)}%
							</div>
						{/if}
						{#if product.listing_type === 'GIFT'}
							<div class="absolute top-2 left-2 badge badge-success">Gratis</div>
						{/if}
					</figure>
					<div class="card-body p-4">
						<h2 class="card-title text-base line-clamp-2">{product.title}</h2>
						<p class="text-sm text-base-content/70 line-clamp-2">{product.description}</p>
						<div class="flex items-center justify-between mt-2">
							<div>
								{#if product.listing_type === 'GIFT'}
									<span class="text-lg font-bold text-success">Gratis</span>
								{:else}
									<span class="text-lg font-bold text-primary">{formatPrice(product.price)}</span>
									{#if product.original_price && product.original_price > product.price}
										<span class="text-sm line-through text-base-content/50 ml-2">
											{formatPrice(product.original_price)}
										</span>
									{/if}
								{/if}
							</div>
							{#if product.quantity_available > 0}
								<span class="badge badge-outline">{product.quantity_available} disp.</span>
							{:else}
								<span class="badge badge-error">Esaurito</span>
							{/if}
						</div>
						<div class="flex items-center gap-2 mt-2 text-xs text-base-content/60">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
							<span>{product.city}</span>
						</div>
					</div>
				</a>
			{/each}
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="flex justify-center mt-8">
				<div class="join">
					<button
						class="join-item btn"
						disabled={currentPage === 1}
						on:click={() => goToPage(currentPage - 1)}
					>
						«
					</button>
					{#each Array(totalPages).fill(0).map((_, i) => i + 1) as page}
						{#if page === 1 || page === totalPages || (page >= currentPage - 1 && page <= currentPage + 1)}
							<button
								class="join-item btn"
								class:btn-active={page === currentPage}
								on:click={() => goToPage(page)}
							>
								{page}
							</button>
						{:else if page === currentPage - 2 || page === currentPage + 2}
							<button class="join-item btn btn-disabled">...</button>
						{/if}
					{/each}
					<button
						class="join-item btn"
						disabled={currentPage === totalPages}
						on:click={() => goToPage(currentPage + 1)}
					>
						»
					</button>
				</div>
			</div>
		{/if}
	{/if}
</div>
