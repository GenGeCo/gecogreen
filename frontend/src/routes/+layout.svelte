<script lang="ts">
	import '../app.css';
</script>

<div class="min-h-screen flex flex-col">
	<!-- Navbar -->
	<header class="navbar bg-base-100 shadow-sm px-4 lg:px-8">
		<div class="flex-1">
			<a href="/" class="text-xl font-bold text-geco-gradient">
				GecoGreen
			</a>
		</div>
		<div class="flex-none gap-2">
			<a href="/login" class="btn btn-ghost btn-sm">Accedi</a>
			<a href="/register" class="btn btn-primary btn-sm">Registrati</a>
		</div>
	</header>

	<!-- Main Content -->
	<main class="flex-1">
		<slot />
	</main>

	<!-- Footer -->
	<footer class="footer footer-center p-4 bg-base-200 text-base-content">
		<aside>
			<p>GecoGreen - La piattaforma antispreco</p>
		</aside>
	</footer>
</div>
