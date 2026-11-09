// Reactive watch-progress store, saved in this browser's localStorage.
// Progress is per browser (per computer): it is deliberately not sent to the
// server. All storage access is wrapped so a private window or blocked
// storage simply means progress isn't remembered.
import {
	STORAGE_KEY,
	complete,
	parseProgress,
	record,
	type LessonProgress,
	type ProgressMap,
	type ProgressUpdate
} from './progress.ts';

function read(): ProgressMap {
	try {
		return parseProgress(localStorage.getItem(STORAGE_KEY));
	} catch {
		return {};
	}
}

function write(map: ProgressMap) {
	try {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(map));
	} catch {
		// storage full or unavailable; keep the in-memory copy
	}
}

class ProgressStore {
	map = $state<ProgressMap>({});

	constructor() {
		if (typeof window === 'undefined') return;
		this.map = read();
		// Keep several open tabs in sync.
		window.addEventListener('storage', (e) => {
			if (e.key === STORAGE_KEY || e.key === null) this.map = parseProgress(e.newValue);
		});
	}

	get(id: string): LessonProgress | undefined {
		return this.map[id];
	}

	save(id: string, update: ProgressUpdate) {
		this.map = record(this.map, id, update, Date.now());
		write(this.map);
	}

	finish(id: string, update: ProgressUpdate) {
		this.map = complete(this.map, id, update, Date.now());
		write(this.map);
	}

	forget(id: string) {
		const { [id]: _removed, ...rest } = this.map;
		void _removed;
		this.map = rest;
		write(this.map);
	}
}

export const progress = new ProgressStore();
