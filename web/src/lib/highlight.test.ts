import { test } from 'node:test';
import assert from 'node:assert/strict';
import { highlight, queryTerms } from './highlight.ts';

const show = (text: string, q: string) =>
	highlight(text, q)
		.map((s) => (s.match ? `[${s.text}]` : s.text))
		.join('');

test('queryTerms', () => {
	assert.deepEqual(queryTerms('  Mid-term  REVIEW '), ['mid', 'term', 'review']);
	assert.deepEqual(queryTerms('class 8 maths'), ['8', 'maths']);
	assert.deepEqual(queryTerms('1.2'), ['1.2']);
	assert.deepEqual(queryTerms('सन्धि'), ['सन्धि']);
});

test('highlight', () => {
	assert.equal(show('Polynomials in Maths', 'poly'), '[Poly]nomials in Maths');
	assert.equal(show('Mid-term Review', 'mid term'), '[Mid]-[term] Review');
	assert.equal(show('Common Factors', 'mon com'), '[Common] Factors', 'overlapping matches merge');
	assert.equal(show('banana', 'an'), 'b[anan]a', 'repeated matches');
	assert.equal(show('Sets', ''), 'Sets');
	assert.equal(show('Sets', 'xyz'), 'Sets');
	assert.equal(show('', 'x'), '');
});
