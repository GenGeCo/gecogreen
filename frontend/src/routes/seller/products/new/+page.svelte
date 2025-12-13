<script lang="ts">
	import { goto } from '$app/navigation';
	import { api, type CreateProductRequest } from '$lib/api';
	import { isAuthenticated, currentUser } from '$lib/stores/auth';

	let loading = false;
	let error = '';
	let imageFiles: FileList | null = null;
	let uploadingImages = false;

	// Form data
	let title = '';
	let description = '';
	let price = 0;
	let originalPrice: number | undefined = undefined;
	let quantity = 1;
	let listingType: 'SALE' | 'GIFT' = 'SALE';
	let shippingMethod: 'PICKUP' | 'SELLER_SHIPS' | 'BUYER_ARRANGES' = 'PICKUP';
	let shippingCost = 0;
	let expiryDate = '';
	let city = '';
	let province = '';

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
				price: listingType === 'GIFT' ? 0 : price,
				original_price: originalPrice,
				quantity,
				listing_type: listingType,
				shipping_method: shippingMethod,
				shipping_cost: shippingMethod === 'PICKUP' ? 0 : shippingCost,
				expiry_date: expiryDate || undefined,
				city: city || undefined,
				province: province || undefined
			};

			const product = await api.createProduct(productData);

			// Upload images if any
			if (imageFiles && imageFiles.length > 0) {
				uploadingImages = true;
				for (let i = 0; i < imageFiles.length && i < 5; i++) {
					try {
						await api.uploadProductImage(product.id, imageFiles[i]);
					} catch (e) {
						console.error('Error uploading image:', e);
					}
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
					<label class="label" for="images">
						<span class="label-text">Immagini (max 5)</span>
					</label>
					<input
						type="file"
						id="images"
						bind:files={imageFiles}
						accept="image/*"
						multiple
						class="file-input file-input-bordered"
					/>
					<label class="label">
						<span class="label-text-alt">JPG, PNG o WebP. Max 5MB per immagine.</span>
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

				<div class="form-control">
					<label class="label" for="quantity">
						<span class="label-text">Quantità Disponibile *</span>
					</label>
					<input
						type="number"
						id="quantity"
						bind:value={quantity}
						class="input input-bordered w-32"
						min="1"
						required
					/>
				</div>

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
					</select>
				</div>

				{#if shippingMethod !== 'PICKUP'}
					<div class="form-control">
						<label class="label" for="shippingCost">
							<span class="label-text">Costo Spedizione</span>
						</label>
						<input
							type="number"
							id="shippingCost"
							bind:value={shippingCost}
							class="input input-bordered w-32"
							min="0"
							step="0.01"
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
