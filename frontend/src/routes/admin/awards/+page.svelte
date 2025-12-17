<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';

	// Types
	interface Task {
		id: string;
		award_id?: string;
		user_id?: string;
		task_type: string;
		title: string;
		description: string;
		status: string;
		priority: string;
		due_date?: string;
		completed_at?: string;
		notes: string;
		created_at: string;
		award?: Award;
		user?: UserMinimal;
	}

	interface Award {
		id: string;
		user_id: string;
		award_type: string;
		period_type: string;
		title: string;
		description: string;
		co2_saved: number;
		youtube_url: string;
		interview_status: string;
		interview_scheduled_at?: string;
		is_featured: boolean;
		is_public: boolean;
		created_at: string;
		user?: UserMinimal;
	}

	interface UserMinimal {
		id: string;
		business_name: string;
		first_name: string;
		last_name: string;
		avatar_url: string;
		city: string;
	}

	interface TasksResponse {
		urgent: Task[];
		this_week: Task[];
		later: Task[];
		completed: Task[];
		stats: {
			total_pending: number;
			total_completed: number;
			overdue_count: number;
		};
	}

	interface LeaderboardEntry {
		rank: number;
		user_id: string;
		business_name: string;
		first_name: string;
		last_name: string;
		total_co2_saved: number;
		total_products_sold: number;
		city: string;
	}

	// State
	let tasks: TasksResponse | null = null;
	let leaderboard: LeaderboardEntry[] = [];
	let loading = true;
	let error = '';
	let success = '';
	let activeTab: 'tasks' | 'awards' | 'leaderboard' = 'tasks';
	let selectedPeriod = 'MONTHLY';

	// Create Award Modal
	let showCreateAward = false;
	let newAward = {
		user_id: '',
		award_type: 'ECO_CHAMPION',
		period_type: 'MONTHLY',
		title: '',
		description: '',
		co2_saved: 0,
		interview_status: 'PENDING'
	};

	// Access control
	$: if ($isAuthenticated !== undefined && !$isAuthenticated) {
		goto('/login');
	}
	$: if ($isAdmin !== undefined && !$isAdmin) {
		goto('/');
	}

	// API helpers
	function getToken(): string {
		return localStorage.getItem('access_token') || '';
	}

	async function fetchTasks() {
		try {
			const response = await fetch('/api/v1/admin/awards/tasks', {
				headers: { Authorization: `Bearer ${getToken()}` }
			});
			const data = await response.json();
			if (!response.ok) throw new Error(data.error);
			tasks = data;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore caricamento tasks';
		}
	}

	async function fetchLeaderboard() {
		try {
			const response = await fetch(`/api/v1/leaderboard?period=${selectedPeriod}&limit=20`);
			const data = await response.json();
			if (!response.ok) throw new Error(data.error);
			leaderboard = data.entries || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore caricamento classifica';
		}
	}

	async function completeTask(taskId: string) {
		const notes = prompt('Note di completamento (opzionale):');
		try {
			const response = await fetch(`/api/v1/admin/awards/tasks/${taskId}/complete`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${getToken()}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ notes: notes || '' })
			});
			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error);
			}
			success = 'Task completato!';
			fetchTasks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore completamento';
		}
	}

	async function updateTaskPriority(taskId: string, priority: string) {
		try {
			const response = await fetch(`/api/v1/admin/awards/tasks/${taskId}`, {
				method: 'PUT',
				headers: {
					Authorization: `Bearer ${getToken()}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ priority })
			});
			if (!response.ok) throw new Error('Errore aggiornamento');
			fetchTasks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore';
		}
	}

	async function createAward() {
		if (!newAward.user_id || !newAward.title) {
			error = 'User ID e Titolo sono obbligatori';
			return;
		}
		try {
			const response = await fetch('/api/v1/admin/awards', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${getToken()}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(newAward)
			});
			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error);
			}
			success = 'Premio creato! I tasks sono stati generati automaticamente.';
			showCreateAward = false;
			newAward = {
				user_id: '',
				award_type: 'ECO_CHAMPION',
				period_type: 'MONTHLY',
				title: '',
				description: '',
				co2_saved: 0,
				interview_status: 'PENDING'
			};
			fetchTasks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore creazione';
		}
	}

	async function quickCreateAward(entry: LeaderboardEntry, awardType: string) {
		const periodNames: Record<string, string> = {
			WEEKLY: 'Settimana',
			MONTHLY: 'Mese',
			YEARLY: 'Anno',
			ALLTIME: 'Sempre'
		};
		const awardTitles: Record<string, string> = {
			ECO_CHAMPION: `Eco-Champion ${periodNames[selectedPeriod]}`,
			ECO_RUNNER_UP: `Secondo Classificato ${periodNames[selectedPeriod]}`,
			ECO_THIRD: `Terzo Classificato ${periodNames[selectedPeriod]}`,
			NEW_ENTRY: 'New Entry del Mese',
			RECORD_BREAKER: 'Record Breaker'
		};

		const displayName = entry.business_name || `${entry.first_name} ${entry.last_name}`;

		if (!confirm(`Creare premio "${awardTitles[awardType]}" per ${displayName}?`)) return;

		try {
			const response = await fetch('/api/v1/admin/awards', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${getToken()}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					user_id: entry.user_id,
					award_type: awardType,
					period_type: selectedPeriod,
					title: awardTitles[awardType],
					description: `${entry.total_co2_saved.toFixed(2)} kg CO2 salvata, ${entry.total_products_sold} prodotti venduti`,
					co2_saved: entry.total_co2_saved,
					interview_status: ['ECO_CHAMPION', 'NEW_ENTRY', 'RECORD_BREAKER'].includes(awardType) ? 'PENDING' : 'NOT_REQUIRED'
				})
			});
			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error);
			}
			success = `Premio creato per ${displayName}! Tasks generati automaticamente.`;
			fetchTasks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore';
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('it-IT', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}

	function getTaskIcon(taskType: string): string {
		const icons: Record<string, string> = {
			INTERVIEW_CONTACT: '📞',
			INTERVIEW_SCHEDULE: '📅',
			INTERVIEW_RECORD: '🎬',
			YOUTUBE_EDIT: '✂️',
			YOUTUBE_PUBLISH: '📺',
			SOCIAL_POST: '📱',
			HALL_OF_FAME_UPDATE: '🏆',
			EMAIL_WINNER: '📧',
			BADGE_ASSIGN: '🏅',
			OTHER: '📋'
		};
		return icons[taskType] || '📋';
	}

	function getPriorityClass(priority: string): string {
		const classes: Record<string, string> = {
			URGENT: 'badge-error',
			HIGH: 'badge-warning',
			NORMAL: 'badge-info',
			LOW: 'badge-ghost'
		};
		return classes[priority] || 'badge-ghost';
	}

	function getAwardTypeLabel(type: string): string {
		const labels: Record<string, string> = {
			ECO_CHAMPION: 'Eco-Champion',
			ECO_RUNNER_UP: 'Secondo',
			ECO_THIRD: 'Terzo',
			NEW_ENTRY: 'New Entry',
			RECORD_BREAKER: 'Record',
			TOP_WEEK: 'Top Settimana',
			ECO_LEGEND: 'Leggenda',
			MILESTONE: 'Milestone'
		};
		return labels[type] || type;
	}

	onMount(async () => {
		if ($isAdmin) {
			await Promise.all([fetchTasks(), fetchLeaderboard()]);
			loading = false;
		}
	});

	$: if ($isAdmin && !loading) {
		fetchLeaderboard();
	}
</script>

<svelte:head>
	<title>Admin - Premi e Contenuti - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-6">
		<div>
			<h1 class="text-3xl font-bold">Premi & Contenuti</h1>
			<p class="text-base-content/70">Gestisci awards, interviste e contenuti YouTube</p>
		</div>
		<div class="flex gap-2">
			<a href="/admin" class="btn btn-ghost">
				← Moderazione
			</a>
			<button class="btn btn-primary" on:click={() => showCreateAward = true}>
				+ Nuovo Premio
			</button>
		</div>
	</div>

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

	<!-- Stats -->
	{#if tasks}
		<div class="grid grid-cols-4 gap-4 mb-8">
			<div class="stat bg-error/20 rounded-lg">
				<div class="stat-title">Urgenti</div>
				<div class="stat-value text-error">{tasks.urgent.length}</div>
			</div>
			<div class="stat bg-warning/20 rounded-lg">
				<div class="stat-title">Questa Settimana</div>
				<div class="stat-value text-warning">{tasks.this_week.length}</div>
			</div>
			<div class="stat bg-info/20 rounded-lg">
				<div class="stat-title">In Attesa</div>
				<div class="stat-value text-info">{tasks.stats.total_pending}</div>
			</div>
			<div class="stat bg-success/20 rounded-lg">
				<div class="stat-title">Completati</div>
				<div class="stat-value text-success">{tasks.stats.total_completed}</div>
			</div>
		</div>
	{/if}

	<!-- Tabs -->
	<div class="tabs tabs-boxed mb-6">
		<button
			class="tab {activeTab === 'tasks' ? 'tab-active' : ''}"
			on:click={() => activeTab = 'tasks'}
		>
			Tasks da Fare
		</button>
		<button
			class="tab {activeTab === 'leaderboard' ? 'tab-active' : ''}"
			on:click={() => activeTab = 'leaderboard'}
		>
			Classifica (Assegna Premi)
		</button>
		<button
			class="tab {activeTab === 'awards' ? 'tab-active' : ''}"
			on:click={() => activeTab = 'awards'}
		>
			Cronologia Completati
		</button>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if activeTab === 'tasks' && tasks}
		<!-- URGENT TASKS -->
		{#if tasks.urgent.length > 0}
			<div class="mb-8">
				<h2 class="text-xl font-bold text-error mb-4 flex items-center gap-2">
					🚨 Urgenti ({tasks.urgent.length})
				</h2>
				<div class="space-y-3">
					{#each tasks.urgent as task}
						<div class="card bg-error/10 border border-error/30">
							<div class="card-body py-4">
								<div class="flex justify-between items-start">
									<div class="flex gap-3">
										<span class="text-2xl">{getTaskIcon(task.task_type)}</span>
										<div>
											<h3 class="font-semibold">{task.title}</h3>
											{#if task.description}
												<p class="text-sm text-base-content/70">{task.description}</p>
											{/if}
											{#if task.due_date}
												<p class="text-xs text-error">Scadenza: {formatDate(task.due_date)}</p>
											{/if}
										</div>
									</div>
									<div class="flex gap-2">
										<select
											class="select select-xs"
											value={task.priority}
											on:change={(e) => updateTaskPriority(task.id, e.currentTarget.value)}
										>
											<option value="URGENT">Urgente</option>
											<option value="HIGH">Alta</option>
											<option value="NORMAL">Normale</option>
											<option value="LOW">Bassa</option>
										</select>
										<button
											class="btn btn-success btn-sm"
											on:click={() => completeTask(task.id)}
										>
											✓ Fatto
										</button>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- THIS WEEK TASKS -->
		{#if tasks.this_week.length > 0}
			<div class="mb-8">
				<h2 class="text-xl font-bold text-warning mb-4 flex items-center gap-2">
					📅 Questa Settimana ({tasks.this_week.length})
				</h2>
				<div class="space-y-3">
					{#each tasks.this_week as task}
						<div class="card bg-base-100 shadow-sm">
							<div class="card-body py-4">
								<div class="flex justify-between items-start">
									<div class="flex gap-3">
										<span class="text-2xl">{getTaskIcon(task.task_type)}</span>
										<div>
											<h3 class="font-semibold">{task.title}</h3>
											{#if task.description}
												<p class="text-sm text-base-content/70">{task.description}</p>
											{/if}
											{#if task.due_date}
												<p class="text-xs text-warning">Scadenza: {formatDate(task.due_date)}</p>
											{/if}
										</div>
									</div>
									<div class="flex gap-2 items-center">
										<span class="badge {getPriorityClass(task.priority)}">{task.priority}</span>
										<button
											class="btn btn-success btn-sm"
											on:click={() => completeTask(task.id)}
										>
											✓ Fatto
										</button>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- LATER TASKS -->
		{#if tasks.later.length > 0}
			<div class="mb-8">
				<h2 class="text-xl font-bold text-base-content/70 mb-4">
					📋 Da Fare ({tasks.later.length})
				</h2>
				<div class="space-y-2">
					{#each tasks.later as task}
						<div class="card bg-base-100 shadow-sm">
							<div class="card-body py-3">
								<div class="flex justify-between items-center">
									<div class="flex gap-3 items-center">
										<span class="text-xl">{getTaskIcon(task.task_type)}</span>
										<span>{task.title}</span>
									</div>
									<div class="flex gap-2 items-center">
										<span class="badge {getPriorityClass(task.priority)} badge-sm">{task.priority}</span>
										<button
											class="btn btn-ghost btn-sm"
											on:click={() => completeTask(task.id)}
										>
											✓
										</button>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		{#if tasks.urgent.length === 0 && tasks.this_week.length === 0 && tasks.later.length === 0}
			<div class="text-center py-12 bg-base-100 rounded-lg">
				<p class="text-6xl mb-4">🎉</p>
				<p class="text-lg">Nessun task in sospeso!</p>
				<p class="text-base-content/70">Usa la classifica per assegnare nuovi premi.</p>
			</div>
		{/if}

	{:else if activeTab === 'leaderboard'}
		<!-- Period selector -->
		<div class="flex gap-2 mb-6">
			{#each ['WEEKLY', 'MONTHLY', 'YEARLY', 'ALLTIME'] as period}
				<button
					class="btn btn-sm {selectedPeriod === period ? 'btn-primary' : 'btn-ghost'}"
					on:click={() => { selectedPeriod = period; fetchLeaderboard(); }}
				>
					{period === 'WEEKLY' ? 'Settimana' : period === 'MONTHLY' ? 'Mese' : period === 'YEARLY' ? 'Anno' : 'Sempre'}
				</button>
			{/each}
		</div>

		<!-- Leaderboard -->
		<div class="overflow-x-auto">
			<table class="table">
				<thead>
					<tr>
						<th>Rank</th>
						<th>Nome</th>
						<th>Città</th>
						<th>CO2 Salvata</th>
						<th>Prodotti</th>
						<th>Azioni</th>
					</tr>
				</thead>
				<tbody>
					{#each leaderboard as entry, i}
						<tr class="{i < 3 ? 'bg-warning/10' : ''}">
							<td>
								<span class="font-bold text-lg">
									{#if entry.rank === 1}🥇{:else if entry.rank === 2}🥈{:else if entry.rank === 3}🥉{:else}{entry.rank}{/if}
								</span>
							</td>
							<td>
								<div class="font-semibold">
									{entry.business_name || `${entry.first_name} ${entry.last_name}`}
								</div>
							</td>
							<td class="text-sm text-base-content/70">{entry.city || '-'}</td>
							<td class="font-mono">{entry.total_co2_saved.toFixed(2)} kg</td>
							<td>{entry.total_products_sold}</td>
							<td>
								<div class="dropdown dropdown-end">
									<button tabindex="0" class="btn btn-primary btn-sm">
										🏆 Assegna Premio
									</button>
									<ul tabindex="0" class="dropdown-content z-10 menu p-2 shadow bg-base-100 rounded-box w-52">
										{#if entry.rank === 1}
											<li><button on:click={() => quickCreateAward(entry, 'ECO_CHAMPION')}>🥇 Eco-Champion</button></li>
										{/if}
										{#if entry.rank === 2}
											<li><button on:click={() => quickCreateAward(entry, 'ECO_RUNNER_UP')}>🥈 Secondo</button></li>
										{/if}
										{#if entry.rank === 3}
											<li><button on:click={() => quickCreateAward(entry, 'ECO_THIRD')}>🥉 Terzo</button></li>
										{/if}
										<li><button on:click={() => quickCreateAward(entry, 'NEW_ENTRY')}>🌟 New Entry</button></li>
										<li><button on:click={() => quickCreateAward(entry, 'RECORD_BREAKER')}>🚀 Record</button></li>
										<li><button on:click={() => quickCreateAward(entry, 'TOP_WEEK')}>📈 Top Week</button></li>
									</ul>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		{#if leaderboard.length === 0}
			<div class="text-center py-12 bg-base-100 rounded-lg">
				<p class="text-lg text-base-content/70">Nessun dato in classifica per questo periodo</p>
			</div>
		{/if}

	{:else if activeTab === 'awards' && tasks}
		<!-- Completed tasks -->
		<h2 class="text-xl font-bold text-success mb-4">Completati ({tasks.completed.length})</h2>
		{#if tasks.completed.length > 0}
			<div class="space-y-2">
				{#each tasks.completed.slice(0, 20) as task}
					<div class="card bg-success/10">
						<div class="card-body py-3">
							<div class="flex justify-between items-center">
								<div class="flex gap-3 items-center">
									<span class="text-xl">{getTaskIcon(task.task_type)}</span>
									<span class="line-through text-base-content/70">{task.title}</span>
								</div>
								<span class="text-sm text-base-content/70">
									{task.completed_at ? formatDate(task.completed_at) : ''}
								</span>
							</div>
							{#if task.notes}
								<p class="text-xs text-base-content/60 ml-9">Note: {task.notes}</p>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<div class="text-center py-12 bg-base-100 rounded-lg">
				<p class="text-lg text-base-content/70">Nessun task completato</p>
			</div>
		{/if}
	{/if}
</div>

<!-- Create Award Modal -->
{#if showCreateAward}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg mb-4">Crea Nuovo Premio</h3>

			<div class="form-control mb-4">
				<label class="label" for="user_id">
					<span class="label-text">User ID</span>
				</label>
				<input
					type="text"
					id="user_id"
					class="input input-bordered"
					placeholder="UUID dell'utente"
					bind:value={newAward.user_id}
				/>
			</div>

			<div class="form-control mb-4">
				<label class="label" for="award_type">
					<span class="label-text">Tipo Premio</span>
				</label>
				<select id="award_type" class="select select-bordered" bind:value={newAward.award_type}>
					<option value="ECO_CHAMPION">Eco-Champion</option>
					<option value="ECO_RUNNER_UP">Secondo Classificato</option>
					<option value="ECO_THIRD">Terzo Classificato</option>
					<option value="NEW_ENTRY">New Entry</option>
					<option value="RECORD_BREAKER">Record Breaker</option>
					<option value="TOP_WEEK">Top Settimana</option>
					<option value="ECO_LEGEND">Eco Legend</option>
					<option value="MILESTONE">Milestone</option>
				</select>
			</div>

			<div class="form-control mb-4">
				<label class="label" for="period_type">
					<span class="label-text">Periodo</span>
				</label>
				<select id="period_type" class="select select-bordered" bind:value={newAward.period_type}>
					<option value="WEEKLY">Settimanale</option>
					<option value="MONTHLY">Mensile</option>
					<option value="YEARLY">Annuale</option>
					<option value="ALLTIME">Di Sempre</option>
				</select>
			</div>

			<div class="form-control mb-4">
				<label class="label" for="title">
					<span class="label-text">Titolo</span>
				</label>
				<input
					type="text"
					id="title"
					class="input input-bordered"
					placeholder="Es: Eco-Champion Dicembre 2024"
					bind:value={newAward.title}
				/>
			</div>

			<div class="form-control mb-4">
				<label class="label" for="description">
					<span class="label-text">Descrizione</span>
				</label>
				<textarea
					id="description"
					class="textarea textarea-bordered"
					placeholder="Descrizione del premio..."
					bind:value={newAward.description}
				></textarea>
			</div>

			<div class="form-control mb-4">
				<label class="label" for="co2_saved">
					<span class="label-text">CO2 Salvata (kg)</span>
				</label>
				<input
					type="number"
					id="co2_saved"
					class="input input-bordered"
					step="0.01"
					bind:value={newAward.co2_saved}
				/>
			</div>

			<div class="form-control mb-4">
				<label class="label" for="interview_status">
					<span class="label-text">Stato Intervista</span>
				</label>
				<select id="interview_status" class="select select-bordered" bind:value={newAward.interview_status}>
					<option value="NOT_REQUIRED">Non Richiesta</option>
					<option value="PENDING">In Attesa</option>
				</select>
			</div>

			<div class="modal-action">
				<button class="btn btn-ghost" on:click={() => showCreateAward = false}>Annulla</button>
				<button class="btn btn-primary" on:click={createAward}>Crea Premio</button>
			</div>
		</div>
		<div class="modal-backdrop" on:click={() => showCreateAward = false} on:keypress={() => {}}></div>
	</div>
{/if}
