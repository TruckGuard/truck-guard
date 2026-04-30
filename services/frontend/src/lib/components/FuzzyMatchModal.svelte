<script lang="ts">
	import * as Dialog from "$lib/components/ui/dialog";
	import { Button } from "$lib/components/ui/button";
	import { Badge } from "$lib/components/ui/badge";
	import { ScrollArea } from "$lib/components/ui/scroll-area";
	import { Image as ImageIcon, Link, Clock, Hash, MapPin, Search } from "@lucide/svelte";
	import { enhance } from "$app/forms";
	import { goto } from "$app/navigation";
	import { toast } from "svelte-sonner";

	let {
		open = $bindable(false),
		data = null,
		onClose = () => {}
	} = $props<{
		open: boolean;
		data: any;
		onClose: () => void;
	}>();

	let linkingId = $state<number | null>(null);

	function getLevenshteinDistance(a: string, b: string): number {
		if (!a || !b) return 999;
		if (a.length === 0) return b.length;
		if (b.length === 0) return a.length;
		const matrix = [];
		for (let i = 0; i <= b.length; i++) matrix[i] = [i];
		for (let j = 0; j <= a.length; j++) matrix[0][j] = j;
		for (let i = 1; i <= b.length; i++) {
			for (let j = 1; j <= a.length; j++) {
				if (b.charAt(i - 1) === a.charAt(j - 1)) {
					matrix[i][j] = matrix[i - 1][j - 1];
				} else {
					matrix[i][j] = Math.min(
						matrix[i - 1][j - 1] + 1,
						Math.min(matrix[i][j - 1] + 1, matrix[i - 1][j] + 1)
					);
				}
			}
		}
		return matrix[b.length][a.length];
	}
</script>

<Dialog.Root bind:open onOpenChange={(v) => { if (!v) onClose(); }}>
	<Dialog.Content class="sm:max-w-[700px] gap-0 p-0 overflow-hidden">
		<!-- Header with incoming event details -->
		<div class="bg-muted/50 p-6 pb-4 border-b">
			<Dialog.Title class="text-xl font-semibold flex items-center gap-2 mb-2">
				<Search class="h-5 w-5 text-primary" />
				Знайдено нечіткі збіги номерного знаку
			</Dialog.Title>
			<Dialog.Description class="text-sm mb-4">
				Камера зафіксувала номер, який не має точного збігу, але схожий на існуючі активні перепустки. Виберіть перепустку зі списку для прив'язки.
			</Dialog.Description>

			{#if data}
				<div class="flex items-center gap-4 bg-background p-4 rounded-lg border shadow-sm">
					{#if data.image_key}
						<div class="w-32 h-20 rounded-md overflow-hidden bg-muted flex items-center justify-center shrink-0 border relative z-10 transition-all duration-300 hover:scale-[3] origin-top-left hover:z-50 hover:shadow-2xl hover:rounded-sm cursor-zoom-in">
							<img src="/api/images/{data.image_key}" alt="Plate crop" class="w-full h-full object-cover" />
						</div>
					{:else}
						<div class="w-32 h-20 rounded-md bg-muted flex flex-col items-center justify-center shrink-0 border text-muted-foreground">
							<ImageIcon class="h-6 w-6 mb-1 opacity-50" />
							<span class="text-[10px] uppercase font-medium">Без фото</span>
						</div>
					{/if}
					
					<div class="flex-1">
						<div class="text-xs text-muted-foreground uppercase tracking-wider font-semibold mb-1">Розпізнаний номер</div>
						<div class="text-3xl font-mono font-bold tracking-widest">{data.incoming_plate}</div>
					</div>
				</div>
			{/if}
		</div>

		<!-- Candidates List -->
		<ScrollArea class="max-h-[50vh] p-6 pt-4">
			{#if data?.candidates && data.candidates.length > 0}
				<div class="text-sm font-medium text-muted-foreground mb-3 flex items-center justify-between">
					<span>Можливі перепустки ({data.candidates.length})</span>
					<span class="text-xs font-normal">Сортування за схожістю</span>
				</div>
				<div class="space-y-3">
					{#each data.candidates as candidate}
						<div class="group flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 rounded-xl border bg-card hover:bg-accent/30 transition-colors relative overflow-hidden">
							<!-- Distance indicator line -->
							<div class="absolute left-0 top-0 bottom-0 w-1 {candidate.distance === 0 ? 'bg-green-500' : candidate.distance === 1 ? 'bg-amber-500' : 'bg-red-500'}"></div>
							
							<div class="flex-1 min-w-0 grid grid-cols-1 sm:grid-cols-2 gap-y-2 gap-x-4 pl-2">
								<div class="flex flex-col">
									<div class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1">
										<Hash class="h-3.5 w-3.5" />
										<span>Код перепустки</span>
									</div>
									<div class="font-semibold text-base">{candidate.code}</div>
								</div>
								
								<div class="flex flex-col">
									<div class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1">
										<Clock class="h-3.5 w-3.5" />
										<span>Час заїзду</span>
									</div>
									<div class="text-sm">
										{new Date(candidate.entry_time).toLocaleString('uk-UA', { 
											day: '2-digit', month: '2-digit', year: 'numeric',
											hour: '2-digit', minute: '2-digit'
										})}
									</div>
								</div>

								<div class="flex flex-col sm:col-span-2 mt-1">
									<div class="flex items-center gap-2 mb-1.5">
										<MapPin class="h-3.5 w-3.5 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">Номери ТЗ</span>
									</div>
									<div class="flex flex-wrap gap-2">
										{#if candidate.plate_front}
											{@const isMatchFront = getLevenshteinDistance(candidate.plate_front, data.incoming_plate) === candidate.distance}
											<Badge variant={isMatchFront ? "default" : "outline"} class="font-mono text-sm px-2 py-0.5">
												Перед: {candidate.plate_front}
											</Badge>
										{/if}
										{#if candidate.plate_back}
											{@const isMatchBack = getLevenshteinDistance(candidate.plate_back, data.incoming_plate) === candidate.distance}
											<Badge variant={isMatchBack ? "default" : "outline"} class="font-mono text-sm px-2 py-0.5">
												Зад: {candidate.plate_back}
											</Badge>
										{/if}
									</div>
								</div>
							</div>

							<div class="flex flex-row sm:flex-col items-center sm:items-end justify-between sm:justify-center gap-3 border-t sm:border-t-0 sm:border-l pt-3 sm:pt-0 sm:pl-4">
								<div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-muted/50 border text-xs font-medium">
									Відхилення: <span class={candidate.distance === 0 ? "text-green-600" : candidate.distance === 1 ? "text-amber-600" : "text-red-600"}>{candidate.distance}</span>
								</div>
								<form 
									action={`/permits/${candidate.id}?/linkEntity`}
									method="POST"
									use:enhance={() => {
										linkingId = candidate.id;
										return async ({ result, update }) => {
											linkingId = null;
											if (result.type === 'success') {
												open = false;
												toast.success("Подію успішно прив'язано до перепустки");
												await goto(`/permits/${candidate.id}`);
											} else {
												if (result.type === 'failure') {
													toast.error("Помилка: " + (result.data?.error || "Не вдалося прив'язати"));
												} else if (result.type === 'error') {
													toast.error("Системна помилка: " + (result.error?.message || "Помилка сервера"));
												}
												await update();
											}
										};
									}}
								>
									<input type="hidden" name="permit_id" value={candidate.id} />
									<input type="hidden" name="event_id" value={data?.event_id} />
									<input type="hidden" name="event_type" value="plate" />
									<Button type="submit" size="sm" variant="default" class="w-full sm:w-auto gap-1.5 shadow-sm" disabled={linkingId !== null}>
										{#if linkingId === candidate.id}
											<div class="h-3.5 w-3.5 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
											<span>Обробка...</span>
										{:else}
											<Link class="h-3.5 w-3.5" />
											<span>Прив'язати</span>
										{/if}
									</Button>
								</form>
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<div class="py-12 text-center text-muted-foreground flex flex-col items-center">
					<Search class="h-10 w-10 mb-3 opacity-20" />
					<p>Кандидатів не знайдено</p>
				</div>
			{/if}
		</ScrollArea>

		<div class="p-4 border-t bg-muted/30 flex justify-end">
			<Button variant="outline" onclick={() => open = false}>Скасувати</Button>
		</div>
	</Dialog.Content>
</Dialog.Root>
