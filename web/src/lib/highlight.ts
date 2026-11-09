// Splits text into plain and matched segments for highlighting search
// results. Pure TypeScript; mirrors the server's case-insensitive matching.

export interface Segment {
	text: string;
	match: boolean;
}

/** Search words from a query, lower-cased, punctuation removed. */
export function queryTerms(query: string): string[] {
	return query
		.toLowerCase()
		.split(/[^\p{L}\p{M}\p{N}.]+/u)
		.filter((t) => t.length > 0 && t !== 'class');
}

export function highlight(text: string, query: string): Segment[] {
	const terms = queryTerms(query);
	if (terms.length === 0 || text === '') return [{ text, match: false }];

	// Mark every character covered by any term (case-insensitive).
	const lower = text.toLowerCase();
	const marked: boolean[] = new Array(text.length).fill(false);
	for (const term of terms) {
		for (let i = lower.indexOf(term); i !== -1; i = lower.indexOf(term, i + 1)) {
			marked.fill(true, i, i + term.length);
		}
	}

	const segments: Segment[] = [];
	for (let i = 0; i < text.length;) {
		let j = i;
		while (j < text.length && marked[j] === marked[i]) j++;
		segments.push({ text: text.slice(i, j), match: marked[i] });
		i = j;
	}
	return segments;
}
