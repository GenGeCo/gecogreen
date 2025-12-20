<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import type { AccountType } from '$lib/api';

	let email = '';
	let password = '';
	let confirmPassword = '';
	let firstName = '';
	let lastName = '';
	let accountType: AccountType = 'PRIVATE';
	let businessName = '';
	let vatNumber = '';
	let fiscalCode = '';
	let sdiCode = '';
	let pecEmail = '';
	let euVatId = '';
	let country = 'IT';
	let hasMultipleLocations = false;
	let city = '';
	let province = '';
	let postalCode = '';
	let error = '';
	let loading = false;

	// Redirect if already logged in
	$: if (typeof window !== 'undefined' && $isAuthenticated) {
		goto('/');
	}

	async function handleSubmit() {
		error = '';

		if (password !== confirmPassword) {
			error = 'Le password non coincidono';
			return;
		}

		if (password.length < 8) {
			error = 'La password deve essere almeno 8 caratteri';
			return;
		}

		if (!city) {
			error = 'La città è obbligatoria';
			return;
		}

		if (accountType === 'BUSINESS') {
			if (!businessName) {
				error = 'La ragione sociale è obbligatoria per account aziendali';
				return;
			}
			if (!vatNumber) {
				error = 'La partita IVA è obbligatoria per account aziendali';
				return;
			}
			if (country === 'IT' && !sdiCode && !pecEmail) {
				error = 'Per aziende italiane è obbligatorio il Codice Univoco SDI o la PEC';
				return;
			}
		}

		loading = true;

		const result = await auth.register({
			email,
			password,
			first_name: firstName,
			last_name: lastName,
			account_type: accountType,
			business_name: accountType === 'BUSINESS' ? businessName : undefined,
			vat_number: accountType === 'BUSINESS' ? vatNumber : undefined,
			fiscal_code: accountType === 'BUSINESS' ? (fiscalCode || undefined) : undefined,
			sdi_code: accountType === 'BUSINESS' && country === 'IT' ? (sdiCode || undefined) : undefined,
			pec_email: accountType === 'BUSINESS' && country === 'IT' ? (pecEmail || undefined) : undefined,
			eu_vat_id: accountType === 'BUSINESS' && country !== 'IT' ? (euVatId || undefined) : undefined,
			billing_country: country || 'IT',
			has_multiple_locations: accountType === 'BUSINESS' ? hasMultipleLocations : false,
			city,
			province: province || undefined,
			postal_code: postalCode || undefined
		});

		if (result.success) {
			goto('/');
		} else {
			error = result.error || 'Errore nella registrazione';
		}

		loading = false;
	}
</script>

<svelte:head>
	<title>Registrati - GecoGreen</title>
</svelte:head>

<div class="min-h-[80vh] flex items-center justify-center p-4">
	<div class="card w-full max-w-lg bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title text-2xl justify-center mb-4">Crea Account</h2>

			{#if error}
				<div class="alert alert-error mb-4">
					<span>{error}</span>
				</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="space-y-4">
				<!-- Account Type Selection -->
				<div class="form-control">
					<label class="label">
						<span class="label-text font-semibold">Tipo Account</span>
					</label>
					<div class="flex gap-4">
						<label class="label cursor-pointer gap-2 flex-1 justify-start border rounded-lg p-4 {accountType === 'PRIVATE' ? 'border-primary bg-primary/10' : 'border-base-300'}">
							<input
								type="radio"
								name="accountType"
								class="radio radio-primary"
								value="PRIVATE"
								bind:group={accountType}
							/>
							<div>
								<span class="label-text font-medium">Privato</span>
								<p class="text-xs opacity-70">Per uso personale</p>
							</div>
						</label>
						<label class="label cursor-pointer gap-2 flex-1 justify-start border rounded-lg p-4 {accountType === 'BUSINESS' ? 'border-primary bg-primary/10' : 'border-base-300'}">
							<input
								type="radio"
								name="accountType"
								class="radio radio-primary"
								value="BUSINESS"
								bind:group={accountType}
							/>
							<div>
								<span class="label-text font-medium">Azienda</span>
								<p class="text-xs opacity-70">Per attività commerciali</p>
							</div>
						</label>
					</div>
				</div>

				<!-- Name Fields -->
				<div class="grid grid-cols-2 gap-4">
					<div class="form-control">
						<label class="label" for="firstName">
							<span class="label-text">Nome</span>
						</label>
						<input
							type="text"
							id="firstName"
							bind:value={firstName}
							class="input input-bordered"
							placeholder="Mario"
							required
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
							placeholder="Rossi"
							required
						/>
					</div>
				</div>

				<!-- Business Fields (only shown for BUSINESS account) -->
				{#if accountType === 'BUSINESS'}
					<div class="bg-base-200 p-4 rounded-lg space-y-4">
						<h3 class="font-semibold text-sm">Dati Aziendali</h3>

						<div class="form-control">
							<label class="label" for="businessName">
								<span class="label-text">Ragione Sociale *</span>
							</label>
							<input
								type="text"
								id="businessName"
								bind:value={businessName}
								class="input input-bordered"
								placeholder="Nome Azienda S.r.l."
								required={accountType === 'BUSINESS'}
							/>
						</div>

						<div class="grid grid-cols-2 gap-4">
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
									required={accountType === 'BUSINESS'}
								/>
							</div>
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
									maxlength="16"
								/>
							</div>
						</div>

						<div class="form-control">
							<label class="label" for="country">
								<span class="label-text">Nazione *</span>
							</label>
							<select id="country" bind:value={country} class="select select-bordered">
								<option value="IT">Italia</option>
								<option value="DE">Germania</option>
								<option value="FR">Francia</option>
								<option value="ES">Spagna</option>
								<option value="AT">Austria</option>
								<option value="BE">Belgio</option>
								<option value="NL">Paesi Bassi</option>
								<option value="PT">Portogallo</option>
								<option value="GR">Grecia</option>
								<option value="PL">Polonia</option>
								<option value="OTHER">Altro paese UE</option>
							</select>
						</div>

						{#if country === 'IT'}
							<div class="alert alert-info text-sm">
								<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-5 h-5">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
								</svg>
								<span>Per la fatturazione elettronica italiana è necessario il <strong>Codice Univoco SDI</strong> oppure la <strong>PEC</strong>.</span>
							</div>

							<div class="grid grid-cols-2 gap-4">
								<div class="form-control">
									<label class="label" for="sdiCode">
										<span class="label-text">Codice Univoco SDI</span>
									</label>
									<input
										type="text"
										id="sdiCode"
										bind:value={sdiCode}
										class="input input-bordered"
										placeholder="XXXXXXX"
										maxlength="7"
									/>
									<label class="label">
										<span class="label-text-alt">7 caratteri alfanumerici</span>
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
										<span class="label-text-alt">Alternativa al Codice SDI</span>
									</label>
								</div>
							</div>
						{:else}
							<div class="form-control">
								<label class="label" for="euVatId">
									<span class="label-text">EU VAT ID</span>
								</label>
								<input
									type="text"
									id="euVatId"
									bind:value={euVatId}
									class="input input-bordered"
									placeholder="{country}123456789"
								/>
								<label class="label">
									<span class="label-text-alt">Partita IVA europea (es. DE123456789, FR12345678901)</span>
								</label>
							</div>
						{/if}

						<div class="form-control">
							<label class="label cursor-pointer justify-start gap-3">
								<input
									type="checkbox"
									bind:checked={hasMultipleLocations}
									class="checkbox checkbox-primary"
								/>
								<div>
									<span class="label-text">Ho più sedi di ritiro</span>
									<p class="text-xs opacity-70">Potrai aggiungere altre sedi dal profilo</p>
								</div>
							</label>
						</div>
					</div>
				{/if}

				<!-- Email -->
				<div class="form-control">
					<label class="label" for="email">
						<span class="label-text">Email</span>
					</label>
					<input
						type="email"
						id="email"
						bind:value={email}
						class="input input-bordered"
						placeholder="email@esempio.com"
						required
					/>
				</div>

				<!-- Password -->
				<div class="form-control">
					<label class="label" for="password">
						<span class="label-text">Password</span>
					</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						class="input input-bordered"
						placeholder="Minimo 8 caratteri"
						required
					/>
				</div>

				<div class="form-control">
					<label class="label" for="confirmPassword">
						<span class="label-text">Conferma Password</span>
					</label>
					<input
						type="password"
						id="confirmPassword"
						bind:value={confirmPassword}
						class="input input-bordered"
						placeholder="Ripeti la password"
						required
					/>
				</div>

				<!-- Location -->
				<div class="bg-base-200 p-4 rounded-lg space-y-4">
					<h3 class="font-semibold text-sm">Sede {accountType === 'BUSINESS' ? 'Principale' : ''}</h3>

					<div class="form-control">
						<label class="label" for="city">
							<span class="label-text">Città *</span>
						</label>
						<input
							type="text"
							id="city"
							bind:value={city}
							class="input input-bordered"
							placeholder="Milano"
							required
						/>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div class="form-control">
							<label class="label" for="province">
								<span class="label-text">Provincia</span>
							</label>
							<input
								type="text"
								id="province"
								bind:value={province}
								class="input input-bordered"
								placeholder="MI"
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
								placeholder="20100"
								maxlength="5"
							/>
						</div>
					</div>
				</div>

				<button type="submit" class="btn btn-primary w-full" disabled={loading}>
					{#if loading}
						<span class="loading loading-spinner"></span>
					{:else}
						Registrati
					{/if}
				</button>
			</form>

			<div class="divider">oppure</div>

			<p class="text-center">
				Hai già un account?
				<a href="/login" class="link link-primary">Accedi</a>
			</p>
		</div>
	</div>
</div>
