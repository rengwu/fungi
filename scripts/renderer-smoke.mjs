// Runs the real CLI/owner/Electron bridge. No browser automation dependency,
// callback, model polling, or synthetic input in a user's session.
import {mkdtemp,writeFile,readFile,cp,mkdir,symlink,readdir} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {createHash} from 'node:crypto';
import path from 'node:path';
import assert from 'node:assert/strict';
import {fileURLToPath} from 'node:url';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const binary=path.resolve(process.argv[2]??path.join(root,'fungi'));
const parent=path.join(root,'verification-output');await mkdir(parent,{recursive:true});
const directory=await mkdtemp(path.join(parent,'renderer-smoke-')),assets=path.join(directory,'assets');
await mkdir(assets);
await cp(new URL('./fixtures/rendering.html',import.meta.url),path.join(assets,'index.html'));
await writeFile(path.join(assets,'module.mjs'),'export const answer=42;');
await writeFile(path.join(assets,'worker.mjs'),"import {answer} from './module.mjs';postMessage(answer);");
await writeFile(path.join(assets,'data.json'),'{"answer":42}');
await writeFile(path.join(assets,'answer.wasm'),Buffer.from('0061736d010000000105016000017f03020100070a0106616e7377657200000a06010400412a0b','hex'));
const wave=Buffer.alloc(44+16000);wave.write('RIFF');wave.writeUInt32LE(wave.length-8,4);wave.write('WAVEfmt ',8);wave.writeUInt32LE(16,16);
wave.writeUInt16LE(1,20);wave.writeUInt16LE(1,22);wave.writeUInt32LE(8000,24);wave.writeUInt32LE(16000,28);wave.writeUInt16LE(2,32);wave.writeUInt16LE(16,34);wave.write('data',36);wave.writeUInt32LE(16000,40);
await writeFile(path.join(assets,'silence.wav'),wave);
await writeFile(path.join(directory,'private.txt'),'must stay outside the skill root');
await symlink(path.join(directory,'private.txt'),path.join(assets,'escape.txt'));
const env={...process.env,FUNGI_STATE_DIR:path.join(directory,'sessions')};
// Windows environment names are case-insensitive; avoid retaining Path too.
for(const key of Object.keys(env))if(key.toLowerCase()==='path')delete env[key];
env.PATH='';
delete env.CODEX_THREAD_ID;
const run=async args=>{
  const {stdout}=await promisify(execFile)(binary,args,{env,timeout:45000,maxBuffer:2*1024*1024});
  return JSON.parse(stdout);
};
try {
const closed=await run(['open',path.join(assets,'index.html'),'--delivery','wait']);
assert.equal(closed.outcome,'cancelled');
const result=await run(['reopen',closed.session,'--delivery','wait']);
assert.equal(result.outcome,'submitted');
const bytes=await readFile(result.artifacts[0].path);
assert.equal(createHash('sha256').update(bytes).digest('hex'),result.artifacts[0].sha256);
const evidence={binary,directory,closed,result,artifactVerified:true,childPathEmpty:true};
await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),checks:result.data.checks},null,2));
assert.ok(result.data.checks.every(check=>check.passed),'renderer checks failed');
} catch(error) {
  // Preserve actual launch diagnostics from this disposable fixture in CI.
  for(const session of await readdir(env.FUNGI_STATE_DIR).catch(()=>[])){
    const log=await readFile(path.join(env.FUNGI_STATE_DIR,session,'worker.log'),'utf8').catch(()=>'');
    if(log)console.error(log);
  }
  throw error;
}
