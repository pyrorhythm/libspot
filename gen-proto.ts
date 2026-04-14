#!/usr/bin/env bun

import { $ } from "bun";
import { cp, mkdir, readdir, rm, stat } from "fs/promises";
import { basename, dirname, extname, join } from "path";
import { fileURLToPath } from "url";

const SCRIPT_DIR = dirname(fileURLToPath(import.meta.url));
const ROOT = SCRIPT_DIR;
const IN = "proto";
const OUT = "gen";
const PROTO_DIR = join(ROOT, IN);
const OUT_DIR = join(ROOT, OUT);
const TMP_DIR = `/tmp/libspot-${Bun.randomUUIDv7()}`;

let goMod = Bun.file(join(SCRIPT_DIR, `go.mod`));
let contents = (await goMod.text()).split("\n")[0]
let rema = (contents ?? "").match(/([a-zA-Z0-9./-_]+)/g)
const MODULE = rema?.at(1);

if (MODULE === undefined) {
  throw Error('could not get go.mod');
}

const sedRules: [RegExp, string][] = [
  [/\/proto/g, ""],
  [/list\/v1\/model\/attributes/g, "list/v1/model"],
  [/player\/context_view\/cyclic_list/g, "player/context_view"],
  [/player\/context_view/g, "player"],
  [/timeline\/v1\/behavior/g, "timeline/v1"],
  [/context\/v1\/behavior/g, "context/v1"],
  // [],
];


async function walkDir(dir: string): Promise<string[]> {
  const protos: string[] = [];
  for (const entry of await readdir(dir)) {
    const full = join(dir, entry);
    const dirstat = await stat(full);
    if (dirstat.isDirectory()) {
      protos.push(...await walkDir(full));
    } else if (entry.endsWith(".proto")) {
      protos.push(full);
    }
  }
  return protos.sort();
}

async function processProto(protoPath: string): Promise<string> {
  const file = Bun.file(protoPath);
  const content = await file.text();
  const baseName = basename(protoPath, extname(protoPath));


  const match = content.match(/^package\p+(\S+);/m);
  if (!match) throw new Error(`No package in ${protoPath}`);

  let pkg = (match[1] ?? "").replace(/\./g, "/");

  if (baseName === "ledger" && dirname(protoPath).endsWith("behavior")) {
    const alias = pkg.split("/").pop() ?? pkg;
    const goPackage = `${MODULE}/${OUT}/${pkg};${alias}`;
    const goPackageLine = `option go_package="${goPackage}";`;

    await Bun.write(file, content.replace(/^syntax\p*=\p*"[^"]+";/m, `$&\n${goPackageLine}`))
    return protoPath
  }

  for (const [re, replacement] of sedRules) {
    if (re.test(pkg)) {
      pkg = pkg.replace(re, replacement);
    }
  }

  const alias = pkg.split("/").pop() ?? pkg;
  const goPackage = `${MODULE}/${OUT}/${pkg};${alias}`;
  const goPackageLine = `option go_package="${goPackage}";`;

  let patched = content;
  if (/^option go_package/m.test(content)) {
    patched = content.replace(/^option go_package="[^"]+";/m, goPackageLine);
  } else {
    patched = content.replace(/^syntax\p*=\p*"[^"]+";/m, `$&\n${goPackageLine}`);
  }

  await Bun.write(protoPath, patched);

  return protoPath;
}

async function rmmk(dir: string): Promise<void> {
  try {
    await rm(dir, { recursive: true });
    await mkdir(dir, { recursive: true });
  } catch (_) { }
}

await rmmk(TMP_DIR);
await rmmk(OUT_DIR);
await cp(join(PROTO_DIR, "."), TMP_DIR, { recursive: true });

try {
  var protos = await walkDir(TMP_DIR);
  var newProtos = new Array<string>();
  for (const proto of protos) {
    const newProto = await processProto(proto)
    if (newProto) newProtos.push(proto);
  }

  console.log(`Generating ${newProtos.length} protos`);

  // Run protoc with explicit cwd - no need for process.chdir()
  await $`protoc --proto_path=${TMP_DIR} --go_out=${OUT_DIR} --go_opt=default_api_level=API_OPAQUE --go_opt=module=${MODULE}/${OUT} ${newProtos}`.cwd(TMP_DIR);
} finally {
  // Always cleanup temp directory
  await rm(TMP_DIR, { recursive: true, force: true });
  console.log("OK")
}