#!/usr/bin/env npx --package=ts-node -- ts-node-esm --swc

import "zx/globals";
import { spinner } from "zx/experimental";

await $`cat package.json | grep name`;

try {
  await $`exit 1`;
} catch (p) {
  console.log(`Exit code: ${p.exitCode}`);
  console.error(`Error: ${p.stderr}`);
}

let bear = await question("What kind of bear is best? ");

console.info(bear);

// With a message.
await spinner("working...", () => $`sleep 2`);
