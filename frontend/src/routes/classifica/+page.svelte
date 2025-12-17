<script lang="ts">
	import { onMount } from 'svelte';

	interface LeaderboardEntry {
		rank: number;
		user_id: string;
		business_name: string;
		first_name: string;
		last_name: string;
		account_type: string;
		city: string;
		avatar_url: string;
		total_co2_saved: number;
		total_water_saved: number;
		total_products_sold: number;
		total_orders: number;
	}

	interface Award {
		id: string;
		award_type: string;
		title: string;
		co2_saved: number;
		youtube_url: string;
		user?: {
			id: string;
			business_name: string;
			first_name: string;
			last_name: string;
			avatar_url: string;
			city: string;
		};
	}

	interface HallOfFame {
		eco_champion?: Award;
		runner_up?: Award;
		third?: Award;
		new_entry?: Award;
		record?: Award;
		period: string;
		period_start: string;
		period_end: string;
	}

	interface CommunityStats {
		total_co2_saved: number;
		total_water_saved: number;
		total_trees_planted: number;
		total_products_saved: number;
		total_users: number;
	}

	let leaderboard: LeaderboardEntry[] = [];
	let hallOfFame: HallOfFame | null = null;
	let communityStats: CommunityStats | null = null;
	let loading = true;
	let selectedPeriod = 'MONTHLY';

	const periodLabels: Record<string, string> = {
		WEEKLY: 'Questa Settimana',
		MONTHLY: 'Questo Mese',
		YEARLY: 'Quest\'Anno',
		ALLTIME: 'Di Sempre'
	};

	async function fetchLeaderboard() {
		try {
			const response = await fetch(`/api/v1/leaderboard?period=${selectedPeriod}&limit=50`);
			const data = await response.json();
			if (response.ok) {
				leaderboard = data.entries || [];
			}
		} catch (e) {
			console.error('Error fetching leaderboard:', e);
		}
	}

	async function fetchHallOfFame() {
		try {
			const response = await fetch('/api/v1/leaderboard/hall-of-fame?period=MONTHLY');
			const data = await response.json();
			if (response.ok) {
				hallOfFame = data;
			}
		} catch (e) {
			console.error('Error fetching hall of fame:', e);
		}
	}

	async function fetchCommunityStats() {
		try {
			const response = await fetch('/api/v1/leaderboard/community-stats');
			const data = await response.json();
			if (response.ok) {
				communityStats = data;
			}
		} catch (e) {
			console.error('Error fetching community stats:', e);
		}
	}

	function getDisplayName(entry: LeaderboardEntry | Award['user']): string {
		if (!entry) return 'Anonimo';
		if ('business_name' in entry && entry.business_name) return entry.business_name;
		if ('first_name' in entry) return `${entry.first_name} ${entry.last_name?.charAt(0) || ''}.`;
		return 'Utente';
	}

	function formatCO2(kg: number): string {
		if (kg >= 1000) {
			return `${(kg / 1000).toFixed(1)}t`;
		}
		return `${kg.toFixed(1)}kg`;
	}

	function getRankEmoji(rank: number): string {
		if (rank === 1) return '🥇';
		if (rank === 2) return '🥈';
		if (rank === 3) return '🥉';
		return `#${rank}`;
	}

	function getRankClass(rank: number): string {
		if (rank === 1) return 'bg-gradient-to-r from-yellow-400/20 to-yellow-600/20 border-yellow-500/50';
		if (rank === 2) return 'bg-gradient-to-r from-gray-300/20 to-gray-400/20 border-gray-400/50';
		if (rank === 3) return 'bg-gradient-to-r from-amber-600/20 to-amber-700/20 border-amber-600/50';
		return '';
	}

	onMount(async () => {
		await Promise.all([
			fetchLeaderboard(),
			fetchHallOfFame(),
			fetchCommunityStats()
		]);
		loading = false;
	});
</script>

<svelte:head>
	<title>Classifica Eco-Campioni - GecoGreen</title>
	<meta name="description" content="Scopri chi sta salvando di più il pianeta con GecoGreen. Classifica settimanale, mensile e annuale degli Eco-Campioni." />
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<!-- Hero Section -->
	<div class="text-center mb-12">
		<h1 class="text-4xl font-bold mb-4">Classifica Eco-Campioni</h1>
		<p class="text-lg text-base-content/70 max-w-2xl mx-auto">
			Chi sta salvando di più il pianeta? Scopri i campioni della sostenibilità che ogni giorno
			combattono lo spreco alimentare.
		</p>
	</div>

	<!-- Community Stats -->
	{#if communityStats}
		<div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-12">
			<div class="stat bg-success/20 rounded-xl text-center">
				<div class="stat-title text-xs">CO2 Risparmiata</div>
				<div class="stat-value text-success text-2xl">{formatCO2(communityStats.total_co2_saved)}</div>
			</div>
			<div class="stat bg-info/20 rounded-xl text-center">
				<div class="stat-title text-xs">Acqua Risparmiata</div>
				<div class="stat-value text-info text-2xl">{(communityStats.total_water_saved / 1000).toFixed(0)}k L</div>
			</div>
			<div class="stat bg-primary/20 rounded-xl text-center">
				<div class="stat-title text-xs">Alberi Piantati</div>
				<div class="stat-value text-primary text-2xl">{communityStats.total_trees_planted}</div>
			</div>
			<div class="stat bg-warning/20 rounded-xl text-center">
				<div class="stat-title text-xs">Prodotti Salvati</div>
				<div class="stat-value text-warning text-2xl">{communityStats.total_products_saved}</div>
			</div>
			<div class="stat bg-secondary/20 rounded-xl text-center">
				<div class="stat-title text-xs">Eco-Guerrieri</div>
				<div class="stat-value text-secondary text-2xl">{communityStats.total_users}</div>
			</div>
		</div>
	{/if}

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else}
		<!-- Hall of Fame -->
		{#if hallOfFame && (hallOfFame.eco_champion || hallOfFame.runner_up || hallOfFame.third)}
			<div class="mb-12">
				<h2 class="text-2xl font-bold text-center mb-6">
					Hall of Fame - {hallOfFame.period}
				</h2>
				<div class="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
					<!-- Eco Champion -->
					{#if hallOfFame.eco_champion}
						<div class="card bg-gradient-to-br from-yellow-400/30 to-yellow-600/30 border-2 border-yellow-500 shadow-xl md:col-start-2 md:row-start-1 md:-mt-4">
							<div class="card-body items-center text-center">
								<div class="text-5xl mb-2">🏆</div>
								<div class="avatar">
									<div class="w-20 rounded-full ring ring-yellow-500 ring-offset-base-100 ring-offset-2">
										{#if hallOfFame.eco_champion.user?.avatar_url}
											<img src={hallOfFame.eco_champion.user.avatar_url} alt="Avatar" />
										{:else}
											<div class="bg-yellow-500 w-full h-full flex items-center justify-center text-3xl">
												{(hallOfFame.eco_champion.user?.first_name?.[0] || '?')}
											</div>
										{/if}
									</div>
								</div>
								<h3 class="text-xl font-bold mt-2">Eco-Champion</h3>
								<p class="font-semibold">{getDisplayName(hallOfFame.eco_champion.user)}</p>
								<p class="text-sm text-base-content/70">{hallOfFame.eco_champion.user?.city || ''}</p>
								<div class="badge badge-success gap-1 mt-2">
									{formatCO2(hallOfFame.eco_champion.co2_saved)} CO2
								</div>
								{#if hallOfFame.eco_champion.youtube_url}
									<a href={hallOfFame.eco_champion.youtube_url} target="_blank" class="btn btn-sm btn-ghost mt-2">
										📺 Guarda Intervista
									</a>
								{/if}
							</div>
						</div>
					{/if}

					<!-- Runner Up -->
					{#if hallOfFame.runner_up}
						<div class="card bg-gradient-to-br from-gray-300/30 to-gray-400/30 border border-gray-400 shadow-lg md:col-start-1 md:row-start-1 md:mt-8">
							<div class="card-body items-center text-center py-6">
								<div class="text-4xl mb-2">🥈</div>
								<div class="avatar">
									<div class="w-16 rounded-full ring ring-gray-400 ring-offset-base-100 ring-offset-2">
										{#if hallOfFame.runner_up.user?.avatar_url}
											<img src={hallOfFame.runner_up.user.avatar_url} alt="Avatar" />
										{:else}
											<div class="bg-gray-400 w-full h-full flex items-center justify-center text-2xl">
												{(hallOfFame.runner_up.user?.first_name?.[0] || '?')}
											</div>
										{/if}
									</div>
								</div>
								<h3 class="font-bold mt-2">Secondo</h3>
								<p class="font-semibold text-sm">{getDisplayName(hallOfFame.runner_up.user)}</p>
								<div class="badge badge-success badge-sm gap-1 mt-2">
									{formatCO2(hallOfFame.runner_up.co2_saved)} CO2
								</div>
							</div>
						</div>
					{/if}

					<!-- Third -->
					{#if hallOfFame.third}
						<div class="card bg-gradient-to-br from-amber-600/30 to-amber-700/30 border border-amber-600 shadow-lg md:col-start-3 md:row-start-1 md:mt-8">
							<div class="card-body items-center text-center py-6">
								<div class="text-4xl mb-2">🥉</div>
								<div class="avatar">
									<div class="w-16 rounded-full ring ring-amber-600 ring-offset-base-100 ring-offset-2">
										{#if hallOfFame.third.user?.avatar_url}
											<img src={hallOfFame.third.user.avatar_url} alt="Avatar" />
										{:else}
											<div class="bg-amber-600 w-full h-full flex items-center justify-center text-2xl">
												{(hallOfFame.third.user?.first_name?.[0] || '?')}
											</div>
										{/if}
									</div>
								</div>
								<h3 class="font-bold mt-2">Terzo</h3>
								<p class="font-semibold text-sm">{getDisplayName(hallOfFame.third.user)}</p>
								<div class="badge badge-success badge-sm gap-1 mt-2">
									{formatCO2(hallOfFame.third.co2_saved)} CO2
								</div>
							</div>
						</div>
					{/if}
				</div>

				<!-- Special Awards -->
				{#if hallOfFame.new_entry || hallOfFame.record}
					<div class="flex justify-center gap-4 mt-6 flex-wrap">
						{#if hallOfFame.new_entry}
							<div class="badge badge-lg badge-outline gap-2 p-4">
								🌟 New Entry: <span class="font-bold">{getDisplayName(hallOfFame.new_entry.user)}</span>
							</div>
						{/if}
						{#if hallOfFame.record}
							<div class="badge badge-lg badge-outline gap-2 p-4">
								🚀 Record: <span class="font-bold">{getDisplayName(hallOfFame.record.user)}</span>
							</div>
						{/if}
					</div>
				{/if}
			</div>

			<div class="divider"></div>
		{/if}

		<!-- Period Selector -->
		<div class="flex justify-center mb-8">
			<div class="tabs tabs-boxed">
				{#each Object.entries(periodLabels) as [period, label]}
					<button
						class="tab {selectedPeriod === period ? 'tab-active' : ''}"
						on:click={() => { selectedPeriod = period; fetchLeaderboard(); }}
					>
						{label}
					</button>
				{/each}
			</div>
		</div>

		<!-- Leaderboard -->
		<h2 class="text-2xl font-bold text-center mb-6">
			Classifica {periodLabels[selectedPeriod]}
		</h2>

		{#if leaderboard.length > 0}
			<div class="max-w-3xl mx-auto">
				<div class="space-y-3">
					{#each leaderboard as entry, i}
						<div class="card {getRankClass(entry.rank)} border {entry.rank <= 3 ? 'border-2' : 'border-base-300'} hover:shadow-lg transition-shadow">
							<div class="card-body py-4 flex-row items-center gap-4">
								<!-- Rank -->
								<div class="text-3xl font-bold w-12 text-center">
									{getRankEmoji(entry.rank)}
								</div>

								<!-- Avatar -->
								<div class="avatar">
									<div class="w-12 rounded-full bg-base-300">
										{#if entry.avatar_url}
											<img src={entry.avatar_url} alt="Avatar" />
										{:else}
											<div class="w-full h-full flex items-center justify-center text-xl font-bold">
												{entry.first_name?.[0] || '?'}
											</div>
										{/if}
									</div>
								</div>

								<!-- Info -->
								<div class="flex-1">
									<div class="font-bold">{getDisplayName(entry)}</div>
									<div class="text-sm text-base-content/70 flex gap-2">
										{#if entry.city}
											<span>📍 {entry.city}</span>
										{/if}
										<span class="badge badge-ghost badge-xs">
											{entry.account_type === 'BUSINESS' ? '🏢 Azienda' : '👤 Privato'}
										</span>
									</div>
								</div>

								<!-- Stats -->
								<div class="text-right">
									<div class="font-bold text-success">{formatCO2(entry.total_co2_saved)}</div>
									<div class="text-xs text-base-content/70">
										{entry.total_products_sold} prodotti
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else}
			<div class="text-center py-12 bg-base-100 rounded-xl shadow">
				<p class="text-6xl mb-4">📊</p>
				<p class="text-lg text-base-content/70">Nessun dato disponibile per questo periodo</p>
				<p class="text-sm text-base-content/50 mt-2">
					Inizia a vendere o comprare prodotti per entrare in classifica!
				</p>
			</div>
		{/if}

		<!-- Call to Action -->
		<div class="text-center mt-12">
			<div class="card bg-primary/10 max-w-xl mx-auto">
				<div class="card-body">
					<h3 class="text-xl font-bold">Vuoi entrare in classifica?</h3>
					<p class="text-base-content/70">
						Ogni prodotto venduto o acquistato ti fa guadagnare EcoCredits e ti avvicina al podio!
					</p>
					<div class="card-actions justify-center mt-4">
						<a href="/prodotti" class="btn btn-primary">
							Esplora Prodotti
						</a>
						<a href="/seller/products/new" class="btn btn-outline btn-primary">
							Inizia a Vendere
						</a>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
