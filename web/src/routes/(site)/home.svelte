<script lang="ts">
	import type { PageProps } from './$types';

	let { data, params }: PageProps = $props();

	let mainAnimals = $state([
		{ src: '/images/animals/cat.svg', alt: 'cats', label: 'Cats' },
		{ src: '/images/animals/dog.svg', alt: 'dogs', label: 'Dogs' },
		{ src: '/images/animals/paw-heart.svg', alt: 'other animals', label: 'Other animals' },
		{ src: '/images/animals/shelter.svg', alt: 'rescues and shelters', label: 'Rescues & shelters' }
	]);

	let otherAnimals = $state([
		{ src: '/images/animals/rabbit.svg', alt: 'rabbits', label: 'Rabbits' },
		{ src: '/images/animals/bird-3.svg', alt: 'birds', label: 'Birds' },
		{ src: '/images/animals/horse-2.svg', alt: 'horses', label: 'Horses' },
		{ src: '/images/animals/mouse.svg', alt: 'small and furry', label: 'Small & furry' },
		{ src: '/images/animals/fish-2.svg', alt: 'scales, fisn and other', label: 'Scales, Fins & other' },
		{ src: '/images/animals/cow.svg', alt: 'barnyard', label: 'Barnyard' }
	]);
</script>

<main>
	<h1 class="mt-6 text-center text-4xl">Find and adopt your best pal with fluffly</h1>
	<h2 class="mt-4 text-center text-xl">Browse your favorite pets</h2>

	<div class="grid gap-4">
		<section class="grid grid-cols-[repeat(auto-fill,minmax(12rem,1fr))] justify-center gap-4">
			{#each mainAnimals as item}
				<div class="flex flex-col items-center justify-center gap-2 rounded-md border border-4 border-yellow-600 p-2">
					<img class="max-w-full object-cover" src={item.src} alt={item.alt} />
					<p class="text-lg">{item.label}</p>
				</div>
			{/each}
		</section>

		<section class="grid grid-cols-[repeat(auto-fill,minmax(12rem,1fr))] justify-center gap-4">
			{#each otherAnimals as item}
				<div class="flex flex-col items-center justify-center gap-2 rounded-md border border-4 border-yellow-600 p-2">
					<img class="max-w-full object-cover" src={item.src} alt={item.alt} />
					<p class="text-lg">{item.label}</p>
				</div>
			{/each}
		</section>
	</div>
</main>
