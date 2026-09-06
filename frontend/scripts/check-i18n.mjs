// Guards the two failure modes that only surface in a production build: a locale
// missing a key another locale has, and a message the compiler rejects. An
// unescaped "@" is the usual culprit — it reads as a linked message. The dev
// runtime recovers from that, the build does not, so the check compiles every
// message the same way the build does.
import { readFileSync } from 'node:fs';
import { baseCompile } from '@intlify/message-compiler';

const locales = ['en', 'ru'];

const load = (code) =>
  JSON.parse(readFileSync(new URL(`../src/i18n/${code}.json`, import.meta.url), 'utf8'));

const flatten = (value, prefix = '') =>
  Object.entries(value).flatMap(([key, child]) =>
    child && typeof child === 'object'
      ? flatten(child, `${prefix}${key}.`)
      : [[`${prefix}${key}`, child]]
  );

const entries = Object.fromEntries(locales.map((code) => [code, flatten(load(code))]));
const problems = [];

const [reference, ...others] = locales;
const referenceKeys = new Set(entries[reference].map(([key]) => key));

for (const code of others) {
  const keys = new Set(entries[code].map(([key]) => key));
  for (const key of referenceKeys) if (!keys.has(key)) problems.push(`${code}: missing "${key}"`);
  for (const key of keys)
    if (!referenceKeys.has(key)) problems.push(`${reference}: missing "${key}"`);
}

for (const code of locales) {
  for (const [key, message] of entries[code]) {
    baseCompile(message, {
      warnHtmlMessage: false,
      onError: (error) => problems.push(`${code}: "${key}" — ${error.message}`),
    });
  }
}

if (problems.length) {
  console.error('i18n check failed:');
  for (const problem of problems) console.error(`  ${problem}`);
  process.exit(1);
}

console.log(`i18n ok: ${locales.length} locales, ${referenceKeys.size} keys`);
