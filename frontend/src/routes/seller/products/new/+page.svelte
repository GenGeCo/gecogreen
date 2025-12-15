<script lang="ts">
	import { goto } from '$app/navigation';
	import { api, type CreateProductRequest } from '$lib/api';
	import { isAuthenticated, currentUser } from '$lib/stores/auth';

	let loading = false;
	let error = '';
	let productImageFiles: FileList | null = null;
	let expiryPhotoFile: FileList | null = null;
	let uploadingImages = false;

	// Form data
	let title = '';
	let description = '';
	let price = 0;
	let originalPrice: number | undefined = undefined;
	let quantity = 1;
	let quantityUnit: import('$lib/api').QuantityUnit = 'PIECE';
	let quantityUnitCustom = '';
	let listingType: 'SALE' | 'GIFT' = 'SALE';
	let shippingMethod: 'PICKUP' | 'SELLER_SHIPS' | 'BUYER_ARRANGES' | 'DIGITAL_FORWARDERS' = 'PICKUP';
	let shippingCost = 0;
	let expiryDate = '';
	let city = '';
	let province = '';

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

	// Pre-fill city from user profile
	$: if ($currentUser && !city) {
		city = $currentUser.city || '';
		province = $currentUser.province || '';
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

		loading = true;

		try {
			const productData: CreateProductRequest = {
				title: title.trim(),
				description: description.trim(),
				price: listingType === 'GIFT' ? 0 : (isDutchAuction ? (dutchStartPrice || price) : price),
				original_price: originalPrice,
				quantity,
				quantity_unit: quantityUnit,
				quantity_unit_custom: quantityUnit === 'CUSTOM' ? quantityUnitCustom : undefined,
				listing_type: listingType,
				shipping_method: shippingMethod,
				shipping_cost: shippingMethod === 'PICKUP' || shippingMethod === 'DIGITAL_FORWARDERS' ? 0 : shippingCost,
				expiry_date: expiryDate || undefined,
				is_dutch_auction: isDutchAuction,
				dutch_start_price: isDutchAuction ? dutchStartPrice : undefined,
				dutch_decrease_amount: isDutchAuction ? dutchDecreaseAmount : undefined,
				dutch_decrease_hours: isDutchAuction ? dutchDecreaseHours : undefined,
				dutch_min_price: isDutchAuction ? dutchMinPrice : undefined,
				city: city || undefined,
				province: province || undefined
			};

			const product = await api.createProduct(productData);

			uploadingImages = true;

			// Upload product images if any
			if (productImageFiles && productImageFiles.length > 0) {
				for (let i = 0; i < productImageFiles.length && i < 5; i++) {
					try {
						await api.uploadProductImage(product.id, productImageFiles[i]);
					} catch (e) {
						console.error('Error uploading product image:', e);
					}
				}
			}

			// Upload expiry photo if present
			if (expiryPhotoFile && expiryPhotoFile.length > 0) {
				try {
					await api.uploadExpiryPhoto(product.id, expiryPhotoFile[0]);
				} catch (e) {
					console.error('Error uploading expiry photo:', e);
				}
			}

			goto('/seller/dashboard');
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nella creazione';
		}

		loading = false;
		uploadingImages = false;
	}
</script>

<svelte:head>
	<title>Nuovo Prodotto - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-3xl">
	<div class="breadcrumbs text-sm mb-6">
		<ul>
			<li><a href="/seller/dashboard">Dashboard</a></li>
			<li>Nuovo Prodotto</li>
		</ul>
	</div>

	<h1 class="text-3xl font-bold mb-6">Pubblica un Prodotto</h1>

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

				<div class="form-control">
					<label class="label" for="productImages">
						<span class="label-text">Foto Prodotto (max 5)</span>
					</label>
					<input
						type="file"
						id="productImages"
						bind:files={productImageFiles}
						accept="image/*"
						multiple
						class="file-input file-input-bordered"
					/>
					<label class="label">
						<span class="label-text-alt">Foto generali del prodotto. JPG, PNG o WebP. Max 5MB per foto.</span>
					</label>
				</div>
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

				<div class="form-control">
					<label class="label" for="expiryPhoto">
						<span class="label-text">Foto Scadenza</span>
					</label>
					<input
						type="file"
						id="expiryPhoto"
						bind:files={expiryPhotoFile}
						accept="image/*"
						class="file-input file-input-bordered"
					/>
					<label class="label">
						<span class="label-text-alt">Foto che mostra chiaramente la data di scadenza sul prodotto. Max 5MB.</span>
					</label>
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

		<!-- Shipping -->
		<div class="card bg-base-100 shadow">
			<div class="card-body">
				<h2 class="card-title text-lg">Spedizione</h2>

				<div class="form-control">
					<label class="label">
						<span class="label-text">Metodo di Consegna</span>
					</label>
					<select class="select select-bordered" bind:value={shippingMethod}>
						<option value="PICKUP">Ritiro in sede</option>
						<option value="SELLER_SHIPS">Spedisco io</option>
						<option value="BUYER_ARRANGES">Organizza l'acquirente</option>
						<option value="DIGITAL_FORWARDERS">Digital Freight Forwarders (Coming Soon)</option>
					</select>
					{#if shippingMethod === 'DIGITAL_FORWARDERS'}
						<label class="label">
							<span class="label-text-alt text-warning">
								📦 Funzionalità in arrivo! Integreremo servizi di spedizione digitali per logistica avanzata.
							</span>
						</label>
					{/if}
				</div>

				{#if shippingMethod !== 'PICKUP' && shippingMethod !== 'DIGITAL_FORWARDERS'}
					<div class="form-control">
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
					</div>
				{/if}

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
				disabled={loading}
			>
				{#if loading}
					<span class="loading loading-spinner"></span>
					{#if uploadingImages}
						Caricamento immagini...
					{:else}
						Pubblicazione...
					{/if}
				{:else}
					Pubblica Prodotto
				{/if}
			</button>
			<a href="/seller/dashboard" class="btn btn-ghost">Annulla</a>
		</div>
	</form>
</div>
