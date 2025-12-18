<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api, type User, type Location, type SocialLinks } from '$lib/api';
	import { auth, isAuthenticated, currentUser, isBusiness } from '$lib/stores/auth';

	let loading = true;
	let saving = false;
	let error = '';
	let success = '';
	let locations: Location[] = [];

	// Form data
	let firstName = '';
	let lastName = '';
	let phone = '';
	let city = '';
	let province = '';
	let postalCode = '';
	let businessName = '';
	let vatNumber = '';
	let socialLinks: SocialLinks = {};
	let hasMultipleLocations = false;
	let accountType: 'PRIVATE' | 'BUSINESS' = 'PRIVATE';

	// Billing info
	let fiscalCode = '';
	let sdiCode = '';
	let pecEmail = '';
	let billingCountry = 'IT';

	// Account type change modal
	let showAccountTypeModal = false;
	let pendingAccountType: 'PRIVATE' | 'BUSINESS' | null = null;

	// New location form
	let showLocationForm = false;
	let newLocation = {
		name: '',
		address_street: '',
		address_city: '',
		address_province: '',
		address_postal_code: '',
		phone: '',
		email: '',
		is_primary: false
	};

	// Redirect if not authenticated
	$: if ($isAuthenticated !== undefined && !$isAuthenticated) {
		goto('/login');
	}

	// Initialize form from user data
	$: if ($currentUser) {
		firstName = $currentUser.first_name || '';
		lastName = $currentUser.last_name || '';
		phone = $currentUser.phone || '';
		city = $currentUser.city || '';
		province = $currentUser.province || '';
		postalCode = $currentUser.postal_code || '';
		businessName = $currentUser.business_name || '';
		vatNumber = $currentUser.vat_number || '';
		socialLinks = $currentUser.social_links || {};
		hasMultipleLocations = $currentUser.has_multiple_locations || false;
		accountType = $currentUser.account_type || 'PRIVATE';
		fiscalCode = $currentUser.fiscal_code || '';
		sdiCode = $currentUser.sdi_code || '';
		pecEmail = $currentUser.pec_email || '';
		billingCountry = $currentUser.billing_country || 'IT';
	}

	async function loadProfile() {
		loading = true;
		try {
			const result = await api.getProfile();
			locations = result.locations || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore caricamento';
		}
		loading = false;
	}

	async function saveProfile() {
		error = '';
		success = '';
		saving = true;

		try {
			const profileData: Record<string, unknown> = {
				first_name: firstName,
				last_name: lastName,
				phone: phone || undefined,
				city: city || undefined,
				province: province || undefined,
				postal_code: postalCode || undefined,
				social_links: socialLinks
			};

			// Add business fields if business account
			if (accountType === 'BUSINESS') {
				profileData.business_name = businessName || undefined;
				profileData.vat_number = vatNumber || undefined;
				profileData.has_multiple_locations = hasMultipleLocations;
				profileData.fiscal_code = fiscalCode || undefined;
				profileData.sdi_code = sdiCode || undefined;
				profileData.pec_email = pecEmail || undefined;
				profileData.billing_country = billingCountry || undefined;
			}

			const updated = await api.updateProfile(profileData);
			auth.updateUser(updated);
			success = 'Profilo aggiornato!';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore salvataggio';
		}
		saving = false;
	}

	function requestAccountTypeChange(newType: 'PRIVATE' | 'BUSINESS') {
		pendingAccountType = newType;
		showAccountTypeModal = true;
	}

	async function confirmAccountTypeChange() {
		if (!pendingAccountType) return;

		// Validation for BUSINESS
		if (pendingAccountType === 'BUSINESS') {
			if (!businessName.trim()) {
				error = 'Inserisci la Ragione Sociale per passare ad account Business';
				showAccountTypeModal = false;
				return;
			}
			if (!vatNumber.trim()) {
				error = 'Inserisci la Partita IVA per passare ad account Business';
				showAccountTypeModal = false;
				return;
			}
		}

		error = '';
		success = '';
		saving = true;
		showAccountTypeModal = false;

		try {
			const profileData: Record<string, unknown> = {
				account_type: pendingAccountType
			};

			// Include business data when switching to BUSINESS
			if (pendingAccountType === 'BUSINESS') {
				profileData.business_name = businessName;
				profileData.vat_number = vatNumber;
				profileData.fiscal_code = fiscalCode || undefined;
				profileData.sdi_code = sdiCode || undefined;
				profileData.pec_email = pecEmail || undefined;
				profileData.billing_country = billingCountry || undefined;
			}

			const updated = await api.updateProfile(profileData);
			auth.updateUser(updated);
			accountType = pendingAccountType;
			success = pendingAccountType === 'BUSINESS'
				? 'Account convertito in Business!'
				: 'Account convertito in Privato!';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore cambio tipo account';
		}
		saving = false;
		pendingAccountType = null;
	}

	async function uploadAvatar(event: Event) {
		const input = event.target as HTMLInputElement;
		if (!input.files || input.files.length === 0) return;

		const file = input.files[0];
		if (file.size > 2 * 1024 * 1024) {
			error = 'Immagine troppo grande (max 2MB)';
			return;
		}

		error = '';
		saving = true;
		try {
			const result = await api.uploadAvatar(file);
			if ($currentUser) {
				auth.updateUser({ ...$currentUser, avatar_url: result.avatar_url });
			}
			success = 'Avatar aggiornato!';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore upload';
		}
		saving = false;
	}

	async function addLocation() {
		if (!newLocation.name || !newLocation.address_street || !newLocation.address_city || !newLocation.address_postal_code) {
			error = 'Compila tutti i campi obbligatori';
			return;
		}

		error = '';
		saving = true;
		try {
			const loc = await api.createLocation({
				name: newLocation.name,
				address_street: newLocation.address_street,
				address_city: newLocation.address_city,
				address_province: newLocation.address_province || undefined,
				address_postal_code: newLocation.address_postal_code,
				phone: newLocation.phone || undefined,
				email: newLocation.email || undefined,
				is_primary: newLocation.is_primary
			});
			locations = [...locations, loc];
			showLocationForm = false;
			newLocation = {
				name: '',
				address_street: '',
				address_city: '',
				address_province: '',
				address_postal_code: '',
				phone: '',
				email: '',
				is_primary: false
			};
			success = 'Sede aggiunta!';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore creazione sede';
		}
		saving = false;
	}

	async function deleteLocation(id: string, name: string) {
		if (!confirm(`Eliminare la sede "${name}"?`)) return;

		try {
			await api.deleteLocation(id);
			locations = locations.filter(l => l.id !== id);
			success = 'Sede eliminata';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore eliminazione';
		}
	}

	onMount(loadProfile);
</script>

<svelte:head>
	<title>Il Mio Profilo - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<h1 class="text-3xl font-bold mb-6">Il Mio Profilo</h1>

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

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else}
		<div class="space-y-6">
			<!-- Avatar Section -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Foto Profilo</h2>
					<div class="flex items-center gap-6">
						<div class="avatar placeholder">
							{#if $currentUser?.avatar_url}
								<div class="w-24 rounded-full">
									<img src={$currentUser.avatar_url} alt="Avatar" />
								</div>
							{:else}
								<div class="bg-primary text-primary-content rounded-full w-24">
									<span class="text-3xl">{firstName?.charAt(0) || 'U'}</span>
								</div>
							{/if}
						</div>
						<div>
							<input
								type="file"
								accept="image/jpeg,image/png,image/webp"
								class="file-input file-input-bordered file-input-sm"
								on:change={uploadAvatar}
							/>
							<p class="text-sm text-base-content/60 mt-1">JPG, PNG o WebP. Max 2MB.</p>
						</div>
					</div>
				</div>
			</div>

			<!-- Account Type Section -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Tipo Account</h2>

					<div class="flex items-center justify-between">
						<div>
							<p class="text-lg font-semibold">
								{accountType === 'BUSINESS' ? 'Account Business' : 'Account Privato'}
							</p>
							<p class="text-sm text-base-content/70">
								{#if accountType === 'BUSINESS'}
									Puoi vendere come azienda con P.IVA, fatturazione elettronica e sedi multiple.
								{:else}
									Puoi vendere come privato. Passa a Business per funzionalità aziendali.
								{/if}
							</p>
						</div>
						<div>
							{#if accountType === 'BUSINESS'}
								<button
									class="btn btn-outline btn-sm"
									on:click={() => requestAccountTypeChange('PRIVATE')}
								>
									Passa a Privato
								</button>
							{:else}
								<button
									class="btn btn-primary btn-sm"
									on:click={() => requestAccountTypeChange('BUSINESS')}
								>
									Passa a Business
								</button>
							{/if}
						</div>
					</div>

					{#if accountType === 'BUSINESS'}
						<div class="alert alert-info mt-4">
							<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
							<span>Come Business puoi avere <strong>{locations.length}</strong> sedi di ritiro configurate.</span>
						</div>
					{/if}
				</div>
			</div>

			<!-- Personal Info -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Informazioni Personali</h2>

					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div class="form-control">
							<label class="label" for="firstName">
								<span class="label-text">Nome</span>
							</label>
							<input
								type="text"
								id="firstName"
								bind:value={firstName}
								class="input input-bordered"
							/>
						</div>

						<div class="form-control">
							<label class="label" for="lastName">
								<span class="label-text">Cognome</span>
							</label>
							<input
								type="text"
								id="lastName"
								bind:value={lastName}
								class="input input-bordered"
							/>
						</div>

						<div class="form-control">
							<label class="label" for="phone">
								<span class="label-text">Telefono</span>
							</label>
							<input
								type="tel"
								id="phone"
								bind:value={phone}
								class="input input-bordered"
								placeholder="+39 xxx xxx xxxx"
							/>
						</div>

						<div class="form-control">
							<label class="label" for="email">
								<span class="label-text">Email</span>
							</label>
							<input
								type="email"
								value={$currentUser?.email || ''}
								class="input input-bordered"
								disabled
							/>
						</div>
					</div>

					<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
						<div class="form-control">
							<label class="label" for="city">
								<span class="label-text">Città</span>
							</label>
							<input
								type="text"
								id="city"
								bind:value={city}
								class="input input-bordered"
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
								maxlength="2"
							/>
						</div>

						<div class="form-control">
							<label class="label" for="postalCode">
								<span class="label-text">CAP</span>
							</label>
							<input
								type="text"
								id="postalCode"
								bind:value={postalCode}
								class="input input-bordered"
								maxlength="5"
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Business Info (if business account or switching to business) -->
			{#if accountType === 'BUSINESS' || pendingAccountType === 'BUSINESS'}
				<div class="card bg-base-100 shadow">
					<div class="card-body">
						<h2 class="card-title text-lg">Informazioni Aziendali</h2>

						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div class="form-control">
								<label class="label" for="businessName">
									<span class="label-text">Ragione Sociale *</span>
								</label>
								<input
									type="text"
									id="businessName"
									bind:value={businessName}
									class="input input-bordered"
									placeholder="Nome Azienda Srl"
								/>
							</div>

							<div class="form-control">
								<label class="label" for="vatNumber">
									<span class="label-text">Partita IVA *</span>
								</label>
								<input
									type="text"
									id="vatNumber"
									bind:value={vatNumber}
									class="input input-bordered"
									placeholder="IT12345678901"
								/>
							</div>
						</div>

						<h3 class="font-semibold mt-6 mb-2">Dati Fatturazione</h3>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div class="form-control">
								<label class="label" for="fiscalCode">
									<span class="label-text">Codice Fiscale</span>
								</label>
								<input
									type="text"
									id="fiscalCode"
									bind:value={fiscalCode}
									class="input input-bordered"
									placeholder="RSSMRA80A01H501U"
								/>
							</div>

							<div class="form-control">
								<label class="label" for="billingCountry">
									<span class="label-text">Paese</span>
								</label>
								<select id="billingCountry" bind:value={billingCountry} class="select select-bordered">
									<option value="IT">Italia</option>
									<option value="DE">Germania</option>
									<option value="FR">Francia</option>
									<option value="ES">Spagna</option>
									<option value="AT">Austria</option>
									<option value="NL">Paesi Bassi</option>
									<option value="BE">Belgio</option>
									<option value="CH">Svizzera</option>
								</select>
							</div>

							{#if billingCountry === 'IT'}
								<div class="form-control">
									<label class="label" for="sdiCode">
										<span class="label-text">Codice SDI</span>
									</label>
									<input
										type="text"
										id="sdiCode"
										bind:value={sdiCode}
										class="input input-bordered"
										placeholder="0000000"
										maxlength="7"
									/>
									<label class="label">
										<span class="label-text-alt">Codice Univoco 7 caratteri</span>
									</label>
								</div>

								<div class="form-control">
									<label class="label" for="pecEmail">
										<span class="label-text">PEC</span>
									</label>
									<input
										type="email"
										id="pecEmail"
										bind:value={pecEmail}
										class="input input-bordered"
										placeholder="azienda@pec.it"
									/>
									<label class="label">
										<span class="label-text-alt">Richiesto se non hai SDI</span>
									</label>
								</div>
							{/if}
						</div>

						<div class="form-control mt-4">
							<label class="label cursor-pointer justify-start gap-3">
								<input
									type="checkbox"
									bind:checked={hasMultipleLocations}
									class="checkbox checkbox-primary"
								/>
								<span class="label-text">Ho più sedi di ritiro</span>
							</label>
						</div>
					</div>
				</div>
			{/if}

			<!-- Social Links -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Social Media</h2>

					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div class="form-control">
							<label class="label" for="instagram">
								<span class="label-text">Instagram</span>
							</label>
							<input
								type="text"
								id="instagram"
								bind:value={socialLinks.instagram}
								class="input input-bordered"
								placeholder="@username"
							/>
						</div>

						<div class="form-control">
							<label class="label" for="facebook">
								<span class="label-text">Facebook</span>
							</label>
							<input
								type="text"
								id="facebook"
								bind:value={socialLinks.facebook}
								class="input input-bordered"
								placeholder="facebook.com/..."
							/>
						</div>

						<div class="form-control">
							<label class="label" for="website">
								<span class="label-text">Sito Web</span>
							</label>
							<input
								type="url"
								id="website"
								bind:value={socialLinks.website}
								class="input input-bordered"
								placeholder="https://..."
							/>
						</div>

						<div class="form-control">
							<label class="label" for="linkedin">
								<span class="label-text">LinkedIn</span>
							</label>
							<input
								type="text"
								id="linkedin"
								bind:value={socialLinks.linkedin}
								class="input input-bordered"
								placeholder="linkedin.com/in/..."
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Save Button -->
			<button
				class="btn btn-primary w-full md:w-auto"
				on:click={saveProfile}
				disabled={saving}
			>
				{#if saving}
					<span class="loading loading-spinner"></span>
				{/if}
				Salva Modifiche
			</button>

			<!-- Locations Section -->
			{#if accountType === 'BUSINESS' || locations.length > 0}
				<div class="card bg-base-100 shadow">
					<div class="card-body">
						<div class="flex justify-between items-center">
							<h2 class="card-title text-lg">Sedi di Ritiro</h2>
							<button
								class="btn btn-primary btn-sm"
								on:click={() => showLocationForm = !showLocationForm}
							>
								{showLocationForm ? 'Annulla' : '+ Aggiungi Sede'}
							</button>
						</div>

						{#if showLocationForm}
							<div class="bg-base-200 p-4 rounded-lg mt-4 space-y-4">
								<h3 class="font-semibold">Nuova Sede</h3>

								<div class="form-control">
									<label class="label"><span class="label-text">Nome Sede *</span></label>
									<input
										type="text"
										bind:value={newLocation.name}
										class="input input-bordered"
										placeholder="es. Sede Centrale"
									/>
								</div>

								<div class="form-control">
									<label class="label"><span class="label-text">Indirizzo *</span></label>
									<input
										type="text"
										bind:value={newLocation.address_street}
										class="input input-bordered"
										placeholder="Via Roma 1"
									/>
								</div>

								<div class="grid grid-cols-3 gap-4">
									<div class="form-control">
										<label class="label"><span class="label-text">Città *</span></label>
										<input
											type="text"
											bind:value={newLocation.address_city}
											class="input input-bordered"
										/>
									</div>
									<div class="form-control">
										<label class="label"><span class="label-text">Provincia</span></label>
										<input
											type="text"
											bind:value={newLocation.address_province}
											class="input input-bordered"
											maxlength="2"
										/>
									</div>
									<div class="form-control">
										<label class="label"><span class="label-text">CAP *</span></label>
										<input
											type="text"
											bind:value={newLocation.address_postal_code}
											class="input input-bordered"
											maxlength="5"
										/>
									</div>
								</div>

								<div class="grid grid-cols-2 gap-4">
									<div class="form-control">
										<label class="label"><span class="label-text">Telefono</span></label>
										<input
											type="tel"
											bind:value={newLocation.phone}
											class="input input-bordered"
										/>
									</div>
									<div class="form-control">
										<label class="label"><span class="label-text">Email</span></label>
										<input
											type="email"
											bind:value={newLocation.email}
											class="input input-bordered"
										/>
									</div>
								</div>

								<div class="form-control">
									<label class="label cursor-pointer justify-start gap-3">
										<input
											type="checkbox"
											bind:checked={newLocation.is_primary}
											class="checkbox checkbox-primary"
										/>
										<span class="label-text">Sede principale</span>
									</label>
								</div>

								<button class="btn btn-primary" on:click={addLocation} disabled={saving}>
									{#if saving}<span class="loading loading-spinner"></span>{/if}
									Aggiungi Sede
								</button>
							</div>
						{/if}

						{#if locations.length > 0}
							<div class="space-y-4 mt-4">
								{#each locations as location}
									<div class="border rounded-lg p-4 flex justify-between items-start">
										<div>
											<div class="font-semibold flex items-center gap-2">
												{location.name}
												{#if location.is_primary}
													<span class="badge badge-primary badge-sm">Principale</span>
												{/if}
											</div>
											<p class="text-sm text-base-content/70">
												{location.address_street}, {location.address_postal_code} {location.address_city}
												{#if location.address_province}({location.address_province}){/if}
											</p>
											{#if location.phone}
												<p class="text-sm">Tel: {location.phone}</p>
											{/if}
										</div>
										<button
											class="btn btn-ghost btn-sm text-error"
											on:click={() => deleteLocation(location.id, location.name)}
										>
											Elimina
										</button>
									</div>
								{/each}
							</div>
						{:else if !showLocationForm}
							<p class="text-base-content/60 text-center py-4">
								Nessuna sede configurata
							</p>
						{/if}
					</div>
				</div>
			{/if}

			<!-- Eco Stats -->
			{#if $currentUser}
				<div class="card bg-base-100 shadow">
					<div class="card-body">
						<h2 class="card-title text-lg">Le Tue Statistiche Eco</h2>
						<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
							<div class="stat bg-success/10 rounded-lg p-4">
								<div class="stat-title text-sm">CO2 Risparmiata</div>
								<div class="stat-value text-success text-xl">{$currentUser.total_co2_saved.toFixed(1)} kg</div>
							</div>
							<div class="stat bg-info/10 rounded-lg p-4">
								<div class="stat-title text-sm">Acqua Risparmiata</div>
								<div class="stat-value text-info text-xl">{$currentUser.total_water_saved.toFixed(0)} L</div>
							</div>
							<div class="stat bg-warning/10 rounded-lg p-4">
								<div class="stat-title text-sm">Eco Crediti</div>
								<div class="stat-value text-warning text-xl">{$currentUser.eco_credits}</div>
							</div>
							<div class="stat bg-primary/10 rounded-lg p-4">
								<div class="stat-title text-sm">Livello</div>
								<div class="stat-value text-primary text-xl">{$currentUser.eco_level}</div>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Account Type Change Modal -->
{#if showAccountTypeModal}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 bg-black/50"
		style="z-index: 9998;"
		on:click={() => { showAccountTypeModal = false; pendingAccountType = null; }}
		on:keypress={() => {}}
		role="button"
		tabindex="-1"
	></div>
	<!-- Modal content -->
	<div
		class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-base-100 rounded-lg shadow-xl p-6 w-full max-w-md mx-4"
		style="z-index: 9999;"
	>
			{#if pendingAccountType === 'BUSINESS'}
				<h3 class="font-bold text-lg">Passa ad Account Business</h3>
				<p class="py-4">
					Stai per convertire il tuo account in <strong>Business</strong>.
				</p>
				<div class="bg-base-200 rounded-lg p-4 mb-4">
					<p class="font-semibold mb-2">Con un account Business potrai:</p>
					<ul class="list-disc list-inside text-sm space-y-1">
						<li>Vendere come azienda con Partita IVA</li>
						<li>Ricevere fatture elettroniche</li>
						<li>Gestire più sedi di ritiro</li>
						<li>Mostrare il nome aziendale nei prodotti</li>
					</ul>
				</div>

				<!-- Campi obbligatori per Business -->
				<div class="space-y-4 mb-4" style="position: relative; z-index: 100;">
					<div class="form-control w-full">
						<label class="label" for="modal-businessName">
							<span class="label-text">Ragione Sociale *</span>
						</label>
						<input
							type="text"
							id="modal-businessName"
							bind:value={businessName}
							class="input input-bordered w-full"
							placeholder="Nome Azienda Srl"
							autocomplete="organization"
							style="pointer-events: auto; position: relative; z-index: 101;"
						/>
					</div>
					<div class="form-control w-full">
						<label class="label" for="modal-vatNumber">
							<span class="label-text">Partita IVA *</span>
						</label>
						<input
							type="text"
							id="modal-vatNumber"
							bind:value={vatNumber}
							class="input input-bordered w-full"
							placeholder="IT12345678901"
							autocomplete="off"
							style="pointer-events: auto; position: relative; z-index: 101;"
						/>
					</div>
				</div>

				{#if !businessName.trim() || !vatNumber.trim()}
					<div class="alert alert-warning">
						<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
						</svg>
						<span>Compila Ragione Sociale e Partita IVA per continuare.</span>
					</div>
				{/if}
			{:else}
				<h3 class="font-bold text-lg">Passa ad Account Privato</h3>
				<p class="py-4">
					Stai per convertire il tuo account in <strong>Privato</strong>.
				</p>
				<div class="alert alert-warning mb-4">
					<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
					<div>
						<p class="font-semibold">Attenzione!</p>
						<p class="text-sm">Passando a Privato perderai l'accesso alle funzionalità Business. I dati aziendali e le sedi rimarranno salvati nel caso volessi tornare a Business.</p>
					</div>
				</div>
			{/if}

			<div class="modal-action">
				<button class="btn btn-ghost" on:click={() => { showAccountTypeModal = false; pendingAccountType = null; }}>
					Annulla
				</button>
				<button
					class="btn {pendingAccountType === 'BUSINESS' ? 'btn-primary' : 'btn-warning'}"
					on:click={confirmAccountTypeChange}
					disabled={saving || (pendingAccountType === 'BUSINESS' && (!businessName.trim() || !vatNumber.trim()))}
				>
					{#if saving}
						<span class="loading loading-spinner"></span>
					{/if}
					Conferma
				</button>
			</div>
		</div>
{/if}
