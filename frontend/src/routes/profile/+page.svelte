<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api, type User, type Location, type SocialLinks } from '$lib/api';
	import { auth, isAuthenticated, currentUser, isBusiness } from '$lib/stores/auth';

	let loading = true;
	let saving = false;
	let error = '';
	let success = '';
	let locationMessage = ''; // Local message for locations section
	let locations: Location[] = [];
	let editingLocationId: string | null = null;

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

	// Business photos
	let businessPhotos: string[] = [];

	// Account type change modal
	let showAccountTypeModal = false;
	let pendingAccountType: 'PRIVATE' | 'BUSINESS' | null = null;
	let dialogRef: HTMLDialogElement;

	// Track if form has been initialized
	let formInitialized = false;

	// Field-level validation errors (real-time)
	let fieldErrors: Record<string, string> = {};

	// Validation functions
	const nameRegex = /^[\p{L}\s'-]{2,50}$/u;

	function validateFirstName(value: string): string {
		if (!value.trim()) return 'Nome obbligatorio';
		if (!nameRegex.test(value.trim())) return '2-50 caratteri, solo lettere';
		return '';
	}

	function validateLastName(value: string): string {
		if (!value.trim()) return 'Cognome obbligatorio';
		if (!nameRegex.test(value.trim())) return '2-50 caratteri, solo lettere';
		return '';
	}

	function validatePhone(value: string): string {
		if (!value) return '';
		const clean = value.replace(/[\s\-\.]/g, '');
		if (!/^\+?[0-9]{6,15}$/.test(clean)) return 'Formato: +39 123 456 7890';
		return '';
	}

	function validatePostalCode(value: string): string {
		if (!value) return '';
		if (!/^[A-Za-z0-9\s-]{3,10}$/.test(value)) return 'CAP non valido';
		return '';
	}

	function validateBusinessName(value: string): string {
		if (accountType !== 'BUSINESS') return '';
		if (!value.trim() || value.trim().length < 2) return 'Min 2 caratteri';
		return '';
	}

	function validateVatNumber(value: string): string {
		if (accountType !== 'BUSINESS') return '';
		if (!value.trim()) return 'P.IVA obbligatoria';
		const clean = value.replace(/[\s\-\.]/g, '');
		if (!/^[A-Za-z0-9]{5,20}$/.test(clean)) return 'Formato non valido';
		return '';
	}

	function validateSdiCode(value: string): string {
		if (accountType !== 'BUSINESS' || billingCountry !== 'IT') return '';
		if (!value && !pecEmail) return 'SDI o PEC richiesto';
		if (value && !/^[A-Za-z0-9]{7}$/.test(value)) return 'Esattamente 7 caratteri';
		return '';
	}

	function validatePecEmail(value: string): string {
		if (accountType !== 'BUSINESS' || billingCountry !== 'IT') return '';
		if (!value && !sdiCode) return 'PEC o SDI richiesto';
		if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return 'Email non valida';
		return '';
	}

	function validateInstagram(value: string): string {
		if (!value) return '';
		// Accept @username or plain username (letters, numbers, underscores, dots)
		const clean = value.replace(/^@/, '');
		if (!/^[a-zA-Z0-9._]{1,30}$/.test(clean)) return 'Username non valido (es. @username)';
		return '';
	}

	function validateFacebook(value: string): string {
		if (!value) return '';
		// Accept URL or page name
		if (value.includes('facebook.com')) {
			if (!/^https?:\/\/(www\.)?facebook\.com\/[a-zA-Z0-9.]+\/?$/.test(value)) {
				return 'URL non valido (es. facebook.com/pagina)';
			}
		} else if (!/^[a-zA-Z0-9.]{3,50}$/.test(value)) {
			return 'Nome pagina non valido';
		}
		return '';
	}

	function validateWebsite(value: string): string {
		if (!value) return '';
		if (!/^https?:\/\/.+\..+/.test(value)) return 'URL non valido (es. https://sito.it)';
		return '';
	}

	function validateLinkedin(value: string): string {
		if (!value) return '';
		if (value.includes('linkedin.com')) {
			if (!/^https?:\/\/(www\.)?linkedin\.com\/(in|company)\/[a-zA-Z0-9-]+\/?$/.test(value)) {
				return 'URL non valido (es. linkedin.com/in/nome)';
			}
		} else if (!/^[a-zA-Z0-9-]{3,50}$/.test(value)) {
			return 'Username non valido';
		}
		return '';
	}

	// Real-time validation on input
	function validateField(field: string) {
		switch(field) {
			case 'firstName': fieldErrors.firstName = validateFirstName(firstName); break;
			case 'lastName': fieldErrors.lastName = validateLastName(lastName); break;
			case 'phone': fieldErrors.phone = validatePhone(phone); break;
			case 'postalCode': fieldErrors.postalCode = validatePostalCode(postalCode); break;
			case 'businessName': fieldErrors.businessName = validateBusinessName(businessName); break;
			case 'vatNumber': fieldErrors.vatNumber = validateVatNumber(vatNumber); break;
			case 'sdiCode':
				fieldErrors.sdiCode = validateSdiCode(sdiCode);
				fieldErrors.pecEmail = validatePecEmail(pecEmail);
				break;
			case 'pecEmail':
				fieldErrors.pecEmail = validatePecEmail(pecEmail);
				fieldErrors.sdiCode = validateSdiCode(sdiCode);
				break;
			case 'instagram': fieldErrors.instagram = validateInstagram(socialLinks.instagram || ''); break;
			case 'facebook': fieldErrors.facebook = validateFacebook(socialLinks.facebook || ''); break;
			case 'website': fieldErrors.website = validateWebsite(socialLinks.website || ''); break;
			case 'linkedin': fieldErrors.linkedin = validateLinkedin(socialLinks.linkedin || ''); break;
		}
		fieldErrors = fieldErrors; // Trigger reactivity
	}

	// Check if form has any errors
	function hasErrors(): boolean {
		return Object.values(fieldErrors).some(e => e !== '');
	}

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

	// Redirect if not authenticated (only in browser, not during SSR)
	$: if (typeof window !== 'undefined' && $isAuthenticated !== undefined && !$isAuthenticated) {
		goto('/login');
	}

	// Initialize form from user data (only once)
	$: if ($currentUser && !formInitialized) {
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
		businessPhotos = $currentUser.business_photos || [];
		formInitialized = true;
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

		// Validate all fields
		fieldErrors.firstName = validateFirstName(firstName);
		fieldErrors.lastName = validateLastName(lastName);
		fieldErrors.phone = validatePhone(phone);
		fieldErrors.postalCode = validatePostalCode(postalCode);
		// Validate social media
		fieldErrors.instagram = validateInstagram(socialLinks.instagram || '');
		fieldErrors.facebook = validateFacebook(socialLinks.facebook || '');
		fieldErrors.website = validateWebsite(socialLinks.website || '');
		fieldErrors.linkedin = validateLinkedin(socialLinks.linkedin || '');
		if (accountType === 'BUSINESS') {
			fieldErrors.businessName = validateBusinessName(businessName);
			fieldErrors.vatNumber = validateVatNumber(vatNumber);
			if (billingCountry === 'IT') {
				fieldErrors.sdiCode = validateSdiCode(sdiCode);
				fieldErrors.pecEmail = validatePecEmail(pecEmail);
			}
		}
		fieldErrors = fieldErrors;

		// Check for errors
		if (hasErrors()) {
			error = 'Correggi i campi evidenziati in rosso';
			return;
		}

		saving = true;

		try {
			const profileData: Record<string, unknown> = {
				first_name: firstName.trim(),
				last_name: lastName.trim(),
				phone: phone?.trim() || undefined,
				city: city?.trim() || undefined,
				province: province?.trim().toUpperCase() || undefined,
				postal_code: postalCode?.trim() || undefined,
				social_links: socialLinks
			};

			// Add business fields if business account
			if (accountType === 'BUSINESS') {
				profileData.business_name = businessName.trim() || undefined;
				profileData.vat_number = vatNumber.trim().toUpperCase() || undefined;
				profileData.has_multiple_locations = hasMultipleLocations;
				profileData.fiscal_code = fiscalCode?.trim().toUpperCase() || undefined;
				profileData.sdi_code = sdiCode?.trim().toUpperCase() || undefined;
				profileData.pec_email = pecEmail?.trim().toLowerCase() || undefined;
				profileData.billing_country = billingCountry || undefined;
				profileData.business_photos = businessPhotos;
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
		// Use native dialog - guaranteed to work
		if (dialogRef) {
			dialogRef.showModal();
		}
	}

	function closeModal() {
		showAccountTypeModal = false;
		pendingAccountType = null;
		if (dialogRef) {
			dialogRef.close();
		}
	}

	async function confirmAccountTypeChange() {
		if (!pendingAccountType) return;

		const newAccountType = pendingAccountType; // Save before closeModal clears it
		error = '';
		success = '';
		saving = true;
		closeModal();

		try {
			const updated = await api.updateProfile({
				account_type: newAccountType
			});
			auth.updateUser(updated);
			accountType = newAccountType;
			success = newAccountType === 'BUSINESS'
				? 'Account convertito in Business! Compila i dati aziendali qui sotto.'
				: 'Account convertito in Privato!';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore cambio tipo account';
		}
		saving = false;
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

	async function uploadBusinessPhoto(event: Event) {
		const input = event.target as HTMLInputElement;
		if (!input.files || input.files.length === 0) return;

		const file = input.files[0];
		if (file.size > 5 * 1024 * 1024) {
			error = 'Immagine troppo grande (max 5MB)';
			return;
		}

		if (businessPhotos.length >= 6) {
			error = 'Massimo 6 foto aziendali';
			return;
		}

		error = '';
		saving = true;
		try {
			const result = await api.uploadBusinessPhoto(file);
			businessPhotos = [...businessPhotos, result.photo_url];
			if ($currentUser) {
				auth.updateUser({ ...$currentUser, business_photos: businessPhotos });
			}
			success = 'Foto aggiunta!';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore upload';
		}
		saving = false;
		input.value = '';
	}

	function removeBusinessPhoto(index: number) {
		businessPhotos = businessPhotos.filter((_, i) => i !== index);
		// Note: We'll save this when user clicks "Salva"
	}

	async function addLocation() {
		if (!newLocation.name || !newLocation.address_street || !newLocation.address_city || !newLocation.address_postal_code) {
			locationMessage = 'error:Compila tutti i campi obbligatori';
			return;
		}

		locationMessage = '';
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
			locationMessage = 'success:Sede aggiunta!';
			setTimeout(() => locationMessage = '', 3000);
		} catch (e) {
			locationMessage = 'error:' + (e instanceof Error ? e.message : 'Errore creazione sede');
		}
		saving = false;
	}

	async function deleteLocation(id: string, name: string) {
		if (!confirm(`Eliminare la sede "${name}"?`)) return;

		try {
			await api.deleteLocation(id);
			locations = locations.filter(l => l.id !== id);
			locationMessage = 'success:Sede eliminata';
			setTimeout(() => locationMessage = '', 3000);
		} catch (e) {
			locationMessage = 'error:' + (e instanceof Error ? e.message : 'Errore eliminazione');
		}
	}

	function startEditLocation(loc: Location) {
		editingLocationId = loc.id;
		newLocation = {
			name: loc.name,
			address_street: loc.address_street,
			address_city: loc.address_city,
			address_province: loc.address_province || '',
			address_postal_code: loc.address_postal_code,
			phone: loc.phone || '',
			email: loc.email || '',
			is_primary: loc.is_primary
		};
		showLocationForm = true;
	}

	function cancelEditLocation() {
		editingLocationId = null;
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
	}

	async function updateLocation() {
		if (!editingLocationId) return;
		if (!newLocation.name || !newLocation.address_street || !newLocation.address_city || !newLocation.address_postal_code) {
			locationMessage = 'error:Compila tutti i campi obbligatori';
			return;
		}

		locationMessage = '';
		saving = true;
		try {
			const updated = await api.updateLocation(editingLocationId, {
				name: newLocation.name,
				address_street: newLocation.address_street,
				address_city: newLocation.address_city,
				address_province: newLocation.address_province || undefined,
				address_postal_code: newLocation.address_postal_code,
				phone: newLocation.phone || undefined,
				email: newLocation.email || undefined,
				is_primary: newLocation.is_primary
			});
			locations = locations.map(l => l.id === editingLocationId ? updated : l);
			cancelEditLocation();
			locationMessage = 'success:Sede aggiornata!';
			setTimeout(() => locationMessage = '', 3000);
		} catch (e) {
			locationMessage = 'error:' + (e instanceof Error ? e.message : 'Errore aggiornamento sede');
		}
		saving = false;
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

				</div>
			</div>

			<!-- Personal Info -->
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-lg">Informazioni Personali</h2>

					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div class="form-control">
							<label class="label" for="firstName">
								<span class="label-text">Nome *</span>
							</label>
							<input
								type="text"
								id="firstName"
								bind:value={firstName}
								on:input={() => validateField('firstName')}
								class="input input-bordered"
								class:input-error={fieldErrors.firstName}
							/>
							{#if fieldErrors.firstName}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.firstName}</span></label>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="lastName">
								<span class="label-text">Cognome *</span>
							</label>
							<input
								type="text"
								id="lastName"
								bind:value={lastName}
								on:input={() => validateField('lastName')}
								class="input input-bordered"
								class:input-error={fieldErrors.lastName}
							/>
							{#if fieldErrors.lastName}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.lastName}</span></label>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="phone">
								<span class="label-text">Telefono</span>
							</label>
							<input
								type="tel"
								id="phone"
								bind:value={phone}
								on:input={() => validateField('phone')}
								class="input input-bordered"
								class:input-error={fieldErrors.phone}
								placeholder="+39 xxx xxx xxxx"
							/>
							{#if fieldErrors.phone}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.phone}</span></label>
							{/if}
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
								on:input={() => validateField('postalCode')}
								class="input input-bordered"
								class:input-error={fieldErrors.postalCode}
								maxlength="5"
							/>
							{#if fieldErrors.postalCode}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.postalCode}</span></label>
							{/if}
						</div>
					</div>
				</div>
			</div>

			<!-- Business Info (if business account) -->
			{#if accountType === 'BUSINESS'}
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
									on:input={() => validateField('businessName')}
									class="input input-bordered"
									class:input-error={fieldErrors.businessName}
									placeholder="Nome Azienda Srl"
								/>
								{#if fieldErrors.businessName}
									<label class="label"><span class="label-text-alt text-error">{fieldErrors.businessName}</span></label>
								{/if}
							</div>

							<div class="form-control">
								<label class="label" for="vatNumber">
									<span class="label-text">Partita IVA *</span>
								</label>
								<input
									type="text"
									id="vatNumber"
									bind:value={vatNumber}
									on:input={() => validateField('vatNumber')}
									class="input input-bordered"
									class:input-error={fieldErrors.vatNumber}
									placeholder="IT12345678901"
								/>
								{#if fieldErrors.vatNumber}
									<label class="label"><span class="label-text-alt text-error">{fieldErrors.vatNumber}</span></label>
								{/if}
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
										on:input={() => validateField('sdiCode')}
										class="input input-bordered"
										class:input-error={fieldErrors.sdiCode}
										placeholder="0000000"
										maxlength="7"
									/>
									{#if fieldErrors.sdiCode}
										<label class="label"><span class="label-text-alt text-error">{fieldErrors.sdiCode}</span></label>
									{:else}
										<label class="label"><span class="label-text-alt">Codice Univoco 7 caratteri</span></label>
									{/if}
								</div>

								<div class="form-control">
									<label class="label" for="pecEmail">
										<span class="label-text">PEC</span>
									</label>
									<input
										type="email"
										id="pecEmail"
										bind:value={pecEmail}
										on:input={() => validateField('pecEmail')}
										class="input input-bordered"
										class:input-error={fieldErrors.pecEmail}
										placeholder="azienda@pec.it"
									/>
									{#if fieldErrors.pecEmail}
										<label class="label"><span class="label-text-alt text-error">{fieldErrors.pecEmail}</span></label>
									{:else}
										<label class="label"><span class="label-text-alt">Richiesto se non hai SDI</span></label>
									{/if}
								</div>
							{/if}
						</div>

						<!-- Business Photos -->
						<h3 class="font-semibold mt-6 mb-2">Foto Aziendali</h3>
						<p class="text-sm text-base-content/70 mb-4">
							Aggiungi foto del tuo negozio, prodotti o attività (max 6 foto, 5MB ciascuna)
						</p>

						<div class="grid grid-cols-2 md:grid-cols-3 gap-4">
							{#each businessPhotos as photo, index}
								<div class="relative aspect-video bg-base-200 rounded-lg overflow-hidden group">
									<img src={photo} alt="Foto aziendale {index + 1}" class="w-full h-full object-cover" />
									<button
										type="button"
										class="absolute top-2 right-2 btn btn-circle btn-sm btn-error opacity-0 group-hover:opacity-100 transition-opacity"
										on:click={() => removeBusinessPhoto(index)}
									>
										<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
										</svg>
									</button>
								</div>
							{/each}

							{#if businessPhotos.length < 6}
								<label class="aspect-video bg-base-200 rounded-lg border-2 border-dashed border-base-300 flex flex-col items-center justify-center cursor-pointer hover:bg-base-300 transition-colors">
									<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-base-content/50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
									</svg>
									<span class="text-sm text-base-content/50 mt-1">Aggiungi foto</span>
									<input
										type="file"
										accept="image/*"
										class="hidden"
										on:change={uploadBusinessPhoto}
									/>
								</label>
							{/if}
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
								on:input={() => validateField('instagram')}
								class="input input-bordered"
								class:input-error={fieldErrors.instagram}
								placeholder="@username"
							/>
							{#if fieldErrors.instagram}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.instagram}</span></label>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="facebook">
								<span class="label-text">Facebook</span>
							</label>
							<input
								type="text"
								id="facebook"
								bind:value={socialLinks.facebook}
								on:input={() => validateField('facebook')}
								class="input input-bordered"
								class:input-error={fieldErrors.facebook}
								placeholder="pagina o facebook.com/pagina"
							/>
							{#if fieldErrors.facebook}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.facebook}</span></label>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="website">
								<span class="label-text">Sito Web</span>
							</label>
							<input
								type="text"
								id="website"
								bind:value={socialLinks.website}
								on:input={() => validateField('website')}
								class="input input-bordered"
								class:input-error={fieldErrors.website}
								placeholder="https://www.miosito.it"
							/>
							{#if fieldErrors.website}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.website}</span></label>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="linkedin">
								<span class="label-text">LinkedIn</span>
							</label>
							<input
								type="text"
								id="linkedin"
								bind:value={socialLinks.linkedin}
								on:input={() => validateField('linkedin')}
								class="input input-bordered"
								class:input-error={fieldErrors.linkedin}
								placeholder="linkedin.com/in/nome"
							/>
							{#if fieldErrors.linkedin}
								<label class="label"><span class="label-text-alt text-error">{fieldErrors.linkedin}</span></label>
							{/if}
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
							{#if !showLocationForm}
								<button
									class="btn btn-primary btn-sm"
									on:click={() => { editingLocationId = null; showLocationForm = true; }}
								>
									+ Aggiungi Sede
								</button>
							{/if}
						</div>

						<!-- Local message for locations -->
						{#if locationMessage}
							<div class="alert {locationMessage.startsWith('success:') ? 'alert-success' : 'alert-error'} mt-2">
								<span>{locationMessage.replace(/^(success:|error:)/, '')}</span>
								<button class="btn btn-sm btn-ghost" on:click={() => locationMessage = ''}>×</button>
							</div>
						{/if}

						{#if showLocationForm}
							<div class="bg-base-200 p-4 rounded-lg mt-4 space-y-4">
								<div class="flex justify-between items-center">
									<h3 class="font-semibold">{editingLocationId ? 'Modifica Sede' : 'Nuova Sede'}</h3>
									<button class="btn btn-ghost btn-sm" on:click={cancelEditLocation}>Annulla</button>
								</div>

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

								<button
									class="btn btn-primary"
									on:click={editingLocationId ? updateLocation : addLocation}
									disabled={saving}
								>
									{#if saving}<span class="loading loading-spinner"></span>{/if}
									{editingLocationId ? 'Salva Modifiche' : 'Aggiungi Sede'}
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
										<div class="flex gap-2">
											<button
												class="btn btn-ghost btn-sm"
												on:click={() => startEditLocation(location)}
											>
												Modifica
											</button>
											<button
												class="btn btn-ghost btn-sm text-error"
												on:click={() => deleteLocation(location.id, location.name)}
											>
												Elimina
											</button>
										</div>
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

<!-- Account Type Change Modal - Native HTML dialog -->
<dialog
	bind:this={dialogRef}
	class="modal"
	on:close={closeModal}
>
	<div class="modal-box">
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
			<p class="text-sm text-base-content/70">
				Dopo la conferma potrai inserire i dati aziendali nella sezione profilo.
			</p>
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
			<button class="btn btn-ghost" on:click={closeModal}>
				Annulla
			</button>
			<button
				class="btn {pendingAccountType === 'BUSINESS' ? 'btn-primary' : 'btn-warning'}"
				on:click={confirmAccountTypeChange}
				disabled={saving}
			>
				{#if saving}
					<span class="loading loading-spinner"></span>
				{/if}
				Conferma
			</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
