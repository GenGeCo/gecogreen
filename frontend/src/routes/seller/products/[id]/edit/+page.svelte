<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api, type CreateProductRequest, type Location, type Product } from '$lib/api';
	import { isAuthenticated, currentUser } from '$lib/stores/auth';

	let loading = true;
	let saving = false;
	let deleting = false;
	let error = '';
	let product: Product | null = null;

	// Form data
	let title = '';
	let description = '';
	let price = 0;
	let originalPrice: number | undefined = undefined;
	let quantity = 1;
	let quantityUnit: import('$lib/api').QuantityUnit = 'PIECE';
	let quantityUnitCustom = '';
	let listingType: 'SALE' | 'GIFT' = 'SALE';
	let expiryDate = '';
	let city = '';
	let province = '';

	// Shipping options (checkboxes)
	let canShip = false;
	let canPickup = true;
	let shippingCost = 0;
	let selectedLocationIds: string[] = [];
	let userLocations: Location[] = [];

	// Dutch Auction
	let isDutchAuction = false;
	let dutchStartPrice: number | undefined = undefined;
	let dutchDecreaseAmount: number | undefined = undefined;
	let dutchDecreaseHours: number | undefined = 24;
	let dutchMinPrice: number | undefined = undefined;

	// Redirect if not authenticated
	$: if ($isAuthenticated !== undefined && !$isAuthenticated) {
		goto('/login');
	}

	// Load product and locations
	onMount(async () => {
		const productId = $page.params.id;

		try {
			// Load product and user locations in parallel
			const [productData, profile] = await Promise.all([
				api.getProduct(productId),
				api.getProfile()
			]);

			product = productData;
			userLocations = profile.locations || [];

			// Check if user owns this product
			if (product.seller_id !== $currentUser?.id) {
				error = 'Non hai i permessi per modificare questo prodotto';
				loading = false;
				return;
			}

			// Populate form fields
			title = product.title;
			description = product.description;
			price = product.price;
			originalPrice = product.original_price;
			quantity = product.quantity;
			quantityUnit = product.quantity_unit;
			quantityUnitCustom = product.quantity_unit_custom || '';
			listingType = product.listing_type;
			city = product.city || '';
			province = product.province || '';
			shippingCost = product.shipping_cost || 0;

			// Parse expiry date (convert from ISO to YYYY-MM-DD for input)
			if (product.expiry_date) {
				expiryDate = product.expiry_date.split('T')[0];
			}

			// Determine shipping options from shipping_method
			if (product.shipping_method === 'BOTH') {
				canShip = true;
				canPickup = true;
			} else if (product.shipping_method === 'SELLER_SHIPS') {
				canShip = true;
				canPickup = false;
			} else {
				canShip = false;
				canPickup = true;
			}

			// Set selected pickup locations
			selectedLocationIds = product.pickup_location_ids || [];

			// Dutch auction settings
			isDutchAuction = product.is_dutch_auction;
			dutchStartPrice = product.dutch_start_price;
			dutchDecreaseAmount = product.dutch_decrease_amount;
			dutchDecreaseHours = product.dutch_decrease_hours || 24;
			dutchMinPrice = product.dutch_min_price;

		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel caricamento del prodotto';
		}

		loading = false;
	});

	function toggleLocation(locationId: string) {
		if (selectedLocationIds.includes(locationId)) {
			selectedLocationIds = selectedLocationIds.filter(id => id !== locationId);
		} else {
			selectedLocationIds = [...selectedLocationIds, locationId];
		}
	}

	async function handleSubmit() {
		error = '';

		if (!title.trim()) {
			error = 'Il titolo è obbligatorio';
			return;
		}
		if (!description.trim()) {
			error = 'La descrizione è obbligatoria';
			return;
		}
		if (listingType === 'SALE' && price <= 0) {
			error = 'Il prezzo deve essere maggiore di 0';
			return;
		}
		if (quantity < 1) {
			error = 'La quantità deve essere almeno 1';
			return;
		}
		if (!canShip && !canPickup) {
			error = 'Seleziona almeno un metodo di consegna (spedizione o ritiro)';
			return;
		}
		if (canPickup && selectedLocationIds.length === 0 && userLocations.length > 0) {
			error = 'Seleziona almeno una sede per il ritiro';
			return;
		}

		saving = true;

		// Determine shipping method based on checkboxes
		let shippingMethod: 'PICKUP' | 'SELLER_SHIPS' | 'BOTH' = 'PICKUP';
		if (canShip && canPickup) {
			shippingMethod = 'BOTH';
		} else if (canShip) {
			shippingMethod = 'SELLER_SHIPS';
		}

		try {
			// Convert date to ISO format for backend (Go expects RFC3339)
			const expiryDateISO = expiryDate ? new Date(expiryDate + 'T23:59:59Z').toISOString() : undefined;

			const productData: Partial<CreateProductRequest> = {
				title: title.trim(),
				description: description.trim(),
				price: listingType === 'GIFT' ? 0 : (isDutchAuction ? (dutchStartPrice || price) : price),
				original_price: originalPrice,
				quantity,
				quantity_unit: quantityUnit,
				quantity_unit_custom: quantityUnit === 'CUSTOM' ? quantityUnitCustom : undefined,
				listing_type: listingType,
				shipping_method: shippingMethod,
				shipping_cost: canShip ? shippingCost : 0,
				expiry_date: expiryDateISO,
				is_dutch_auction: isDutchAuction,
				dutch_start_price: isDutchAuction ? dutchStartPrice : undefined,
				dutch_decrease_amount: isDutchAuction ? dutchDecreaseAmount : undefined,
				dutch_decrease_hours: isDutchAuction ? dutchDecreaseHours : undefined,
				dutch_min_price: isDutchAuction ? dutchMinPrice : undefined,
				city: city || undefined,
				province: province || undefined,
				pickup_location_ids: canPickup ? selectedLocationIds : undefined
			};

			await api.updateProduct($page.params.id, productData);

			goto(`/products/${$page.params.id}`);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel salvataggio';
		}

		saving = false;
	}

	async function handleDelete() {
		if (!confirm('Sei sicuro di voler eliminare questo prodotto? Questa azione non può essere annullata.')) {
			return;
		}
		deleting = true;
		error = '';
		try {
			await api.deleteProduct($page.params.id);
			goto('/seller/dashboard');
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nell\'eliminazione';
			deleting = false;
		}
	}
</script>

<svelte:head>
	<title>Modifica Prodotto - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-3xl">
	<div class="breadcrumbs text-sm mb-6">
		<ul>
			<li><a href="/seller/dashboard">Dashboard</a></li>
			<li><a href="/products/{$page.params.id}">Prodotto</a></li>
			<li>Modifica</li>
		</ul>
	</div>

	<h1 class="text-3xl font-bold mb-6">Modifica Prodotto</h1>

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if error && !product}
		<div class="alert alert-error">
			<span>{error}</span>
		</div>
		<a href="/seller/dashboard" class="btn btn-primary mt-4">Torna alla Dashboard</a>
	{:else}
		{#if error}
			<div class="alert alert-error mb-6">
				<span>{error}</span>
			</div>
		{/if}

		<form on:submit|preventDefault={handleSubmit} class="space-y-6">
			<!-- Basic Info -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Informazioni Base</h2>

					<div class="form-control">
						<label class="label" for="title">
							<span class="label-text">Titolo *</span>
						</label>
						<input
							type="text"
							id="title"
							bind:value={title}
							class="input input-bordered"
							placeholder="es. Frutta mista vicino scadenza"
							maxlength="200"
							required
						/>
					</div>

					<div class="form-control">
						<label class="label" for="description">
							<span class="label-text">Descrizione *</span>
						</label>
						<textarea
							id="description"
							bind:value={description}
							class="textarea textarea-bordered h-32"
							placeholder="Descrivi il prodotto, le condizioni, perché lo vendi..."
							required
						></textarea>
					</div>

					{#if product?.images && product.images.length > 0}
						<div class="form-control">
							<label class="label">
								<span class="label-text">Immagini attuali</span>
							</label>
							<div class="flex gap-2 flex-wrap">
								{#each product.images as img, i}
									<img src={img} alt="Immagine {i + 1}" class="w-20 h-20 object-cover rounded-lg border" />
								{/each}
							</div>
							<label class="label">
								<span class="label-text-alt">Per modificare le immagini, elimina il prodotto e creane uno nuovo</span>
							</label>
						</div>
					{/if}
				</div>
			</div>

			<!-- Expiry Information -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Scadenza</h2>

					<div class="form-control">
						<label class="label" for="expiryDate">
							<span class="label-text">Data di Scadenza</span>
						</label>
						<input
							type="date"
							id="expiryDate"
							bind:value={expiryDate}
							class="input input-bordered w-48"
						/>
					</div>
				</div>
			</div>

			<!-- Pricing -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Prezzo e Quantità</h2>

					<div class="form-control">
						<label class="label">
							<span class="label-text">Tipo di Inserzione</span>
						</label>
						<div class="flex gap-4">
							<label class="label cursor-pointer gap-2">
								<input
									type="radio"
									name="listingType"
									class="radio radio-primary"
									value="SALE"
									bind:group={listingType}
								/>
								<span class="label-text">Vendita</span>
							</label>
							<label class="label cursor-pointer gap-2">
								<input
									type="radio"
									name="listingType"
									class="radio radio-primary"
									value="GIFT"
									bind:group={listingType}
								/>
								<span class="label-text">Regalo</span>
							</label>
						</div>
					</div>

					{#if listingType === 'SALE'}
						<div class="grid grid-cols-2 gap-4">
							<div class="form-control">
								<label class="label" for="price">
									<span class="label-text">Prezzo *</span>
								</label>
								<input
									type="number"
									id="price"
									bind:value={price}
									class="input input-bordered"
									min="0.01"
									step="0.01"
									required
								/>
							</div>
							<div class="form-control">
								<label class="label" for="originalPrice">
									<span class="label-text">Prezzo Originale</span>
								</label>
								<input
									type="number"
									id="originalPrice"
									bind:value={originalPrice}
									class="input input-bordered"
									min="0"
									step="0.01"
									placeholder="Opzionale"
								/>
							</div>
						</div>
					{/if}

					<div class="grid grid-cols-3 gap-4">
						<div class="form-control">
							<label class="label" for="quantity">
								<span class="label-text">Quantità *</span>
							</label>
							<input
								type="number"
								id="quantity"
								bind:value={quantity}
								class="input input-bordered"
								min="1"
								step="0.1"
								required
							/>
						</div>
						<div class="form-control">
							<label class="label" for="quantityUnit">
								<span class="label-text">Unità *</span>
							</label>
							<select id="quantityUnit" bind:value={quantityUnit} class="select select-bordered">
								<option value="PIECE">Pezzi</option>
								<option value="KG">Kg</option>
								<option value="G">Grammi</option>
								<option value="L">Litri</option>
								<option value="ML">Millilitri</option>
								<option value="CUSTOM">Altro...</option>
							</select>
						</div>
						{#if quantityUnit === 'CUSTOM'}
							<div class="form-control">
								<label class="label" for="quantityUnitCustom">
									<span class="label-text">Specifica *</span>
								</label>
								<input
									type="text"
									id="quantityUnitCustom"
									bind:value={quantityUnitCustom}
									class="input input-bordered"
									placeholder="es. confezioni"
									required
								/>
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Dutch Auction (Asta al Contrario) -->
			{#if listingType === 'SALE'}
				<div class="card bg-base-100 shadow">
					<div class="card-body">
						<div class="flex items-center justify-between">
							<div>
								<h2 class="card-title text-lg">Asta al Contrario (Olandese)</h2>
								<p class="text-sm text-base-content/70 mt-1">
									Il prezzo inizia alto e scende automaticamente fino ad un minimo
								</p>
							</div>
							<input
								type="checkbox"
								bind:checked={isDutchAuction}
								class="toggle toggle-primary"
							/>
						</div>

						{#if isDutchAuction}
							<div class="grid grid-cols-2 gap-4 mt-4">
								<div class="form-control">
									<label class="label" for="dutchStartPrice">
										<span class="label-text">Prezzo Iniziale (€) *</span>
									</label>
									<input
										type="number"
										id="dutchStartPrice"
										bind:value={dutchStartPrice}
										class="input input-bordered"
										min="0.01"
										step="0.01"
										placeholder="es. 10.00"
										required
									/>
								</div>
								<div class="form-control">
									<label class="label" for="dutchMinPrice">
										<span class="label-text">Prezzo Minimo (€) *</span>
									</label>
									<input
										type="number"
										id="dutchMinPrice"
										bind:value={dutchMinPrice}
										class="input input-bordered"
										min="0.01"
										step="0.01"
										placeholder="es. 2.00"
										required
									/>
								</div>
								<div class="form-control">
									<label class="label" for="dutchDecreaseAmount">
										<span class="label-text">Riduzione Prezzo (€) *</span>
									</label>
									<input
										type="number"
										id="dutchDecreaseAmount"
										bind:value={dutchDecreaseAmount}
										class="input input-bordered"
										min="0.01"
										step="0.01"
										placeholder="es. 1.00"
										required
									/>
									<label class="label">
										<span class="label-text-alt">Quanto scende il prezzo ad ogni intervallo</span>
									</label>
								</div>
								<div class="form-control">
									<label class="label" for="dutchDecreaseHours">
										<span class="label-text">Ogni Quante Ore *</span>
									</label>
									<input
										type="number"
										id="dutchDecreaseHours"
										bind:value={dutchDecreaseHours}
										class="input input-bordered"
										min="1"
										step="1"
										placeholder="es. 24"
										required
									/>
									<label class="label">
										<span class="label-text-alt">Frequenza riduzione prezzo (in ore)</span>
									</label>
								</div>
							</div>

							<div class="alert alert-info mt-4">
								<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
								</svg>
								<div class="text-sm">
									<p><strong>Esempio:</strong> Prezzo iniziale €{dutchStartPrice || 10}, scende di €{dutchDecreaseAmount || 1} ogni {dutchDecreaseHours || 24} ore fino a €{dutchMinPrice || 2}</p>
									<p class="mt-1">Il prezzo scenderà automaticamente fino a quando qualcuno acquista o si raggiunge il minimo.</p>
								</div>
							</div>
						{/if}
					</div>
				</div>
			{/if}

			<!-- Shipping & Pickup -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Consegna</h2>

					<div class="space-y-4">
						<!-- Shipping option -->
						<div class="form-control">
							<label class="label cursor-pointer justify-start gap-3">
								<input
									type="checkbox"
									bind:checked={canShip}
									class="checkbox checkbox-primary"
								/>
								<div>
									<span class="label-text font-medium">Spedizione disponibile</span>
									<p class="text-xs text-base-content/60">Puoi spedire il prodotto all'acquirente</p>
								</div>
							</label>
						</div>

						{#if canShip}
							<div class="form-control ml-8">
								<label class="label" for="shippingCost">
									<span class="label-text">Costo Spedizione (€)</span>
								</label>
								<input
									type="number"
									id="shippingCost"
									bind:value={shippingCost}
									class="input input-bordered w-32"
									min="0"
									step="0.01"
									placeholder="0.00"
								/>
								<label class="label">
									<span class="label-text-alt">0 = spedizione gratuita</span>
								</label>
							</div>
						{/if}

						<!-- Pickup option -->
						<div class="form-control">
							<label class="label cursor-pointer justify-start gap-3">
								<input
									type="checkbox"
									bind:checked={canPickup}
									class="checkbox checkbox-primary"
								/>
								<div>
									<span class="label-text font-medium">Ritiro in sede</span>
									<p class="text-xs text-base-content/60">L'acquirente può ritirare il prodotto</p>
								</div>
							</label>
						</div>

						{#if canPickup}
							{#if userLocations.length > 0}
								<div class="ml-8 space-y-2">
									<label class="label">
										<span class="label-text">Seleziona le sedi per il ritiro</span>
									</label>
									{#each userLocations as loc}
										<label class="flex items-center gap-3 cursor-pointer p-2 rounded hover:bg-base-200">
											<input
												type="checkbox"
												checked={selectedLocationIds.includes(loc.id)}
												on:change={() => toggleLocation(loc.id)}
												class="checkbox checkbox-sm"
											/>
											<div class="flex-1">
												<span class="font-medium">{loc.name}</span>
												<span class="text-sm text-base-content/60 ml-2">
													{loc.address_city}{loc.address_province ? `, ${loc.address_province}` : ''}
												</span>
											</div>
										</label>
									{/each}
								</div>
							{:else}
								<div class="alert alert-warning ml-8">
									<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
									</svg>
									<div>
										<span>Non hai ancora aggiunto sedi di ritiro.</span>
										<a href="/profile" class="link link-primary">Aggiungi una sede</a>
									</div>
								</div>
							{/if}
						{/if}

						<!-- Digital Freight Forwarders (Coming Soon) -->
						<div class="form-control opacity-60">
							<label class="label cursor-not-allowed justify-start gap-3">
								<input
									type="checkbox"
									disabled
									class="checkbox"
								/>
								<div>
									<span class="label-text font-medium">
										Spedizione Industriale
										<span class="badge badge-sm badge-info ml-2">Coming Soon</span>
									</span>
									<p class="text-xs text-base-content/60">Spedizionieri per bancali e carichi industriali (Digital Freight Forwarders)</p>
								</div>
							</label>
						</div>
					</div>

					<div class="divider"></div>

					<div class="grid grid-cols-2 gap-4">
						<div class="form-control">
							<label class="label" for="city">
								<span class="label-text">Città</span>
							</label>
							<input
								type="text"
								id="city"
								bind:value={city}
								class="input input-bordered"
								placeholder="es. Milano"
							/>
							<label class="label">
								<span class="label-text-alt">Visibile agli acquirenti prima dell'acquisto</span>
							</label>
						</div>
						<div class="form-control">
							<label class="label" for="province">
								<span class="label-text">Provincia</span>
							</label>
							<input
								type="text"
								id="province"
								bind:value={province}
								class="input input-bordered"
								placeholder="es. MI"
								maxlength="2"
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Submit -->
			<div class="flex gap-4">
				<button
					type="submit"
					class="btn btn-primary flex-1"
					disabled={saving}
				>
					{#if saving}
						<span class="loading loading-spinner"></span>
						Salvataggio...
					{:else}
						Salva Modifiche
					{/if}
				</button>
				<a href="/products/{$page.params.id}" class="btn btn-ghost">Annulla</a>
			</div>

			<!-- Delete -->
			<div class="divider"></div>
			<div class="flex justify-end">
				<button
					type="button"
					class="btn btn-error btn-outline"
					on:click={handleDelete}
					disabled={deleting}
				>
					{#if deleting}
						<span class="loading loading-spinner loading-sm"></span>
						Eliminazione...
					{:else}
						Elimina Prodotto
					{/if}
				</button>
			</div>
		</form>
	{/if}
</div>
