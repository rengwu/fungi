// Verify a packaged engine with an external helper, empty PATH and no display.
import {mkdtemp,mkdir,readFile,writeFile,copyFile,chmod,readdir} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {createHash} from 'node:crypto';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
import assert from 'node:assert/strict';

assert.equal(process.argv.length,3,'usage: engine-smoke.mjs /absolute/package/fungi');
const binary=path.resolve(process.argv[2]),root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const parent=path.join(root,'verification-output');await mkdir(parent,{recursive:true});
const directory=await mkdtemp(path.join(parent,'engine-smoke-'));
const standalone=path.join(directory,process.platform==='win32'?'fungi.exe':'fungi');
await copyFile(binary,standalone);await chmod(standalone,0o755);
const helper=path.join(root,'scripts/helper-fixture.mjs'),state=path.join(directory,'sessions');
const env={...process.env,FUNGI_STATE_DIR:state};
for(const key of Object.keys(env))if(['path','codex_thread_id','display','wayland_display','fungi_renderer_dir'].includes(key.toLowerCase()))delete env[key];
env.PATH='';
const execute=promisify(execFile);
const invoke=async args=>(await execute(standalone,args,{env,cwd:directory,timeout:30000,maxBuffer:2**20})).stdout;
const run=async args=>JSON.parse(await invoke(args));
const evidence={directory,sourceBinary:binary,platform:process.platform,arch:process.arch,standalone:true,childPathEmpty:true,displayRemoved:true,checks:[]};
async function fileReady(file){
  // Bounded machine-fixture synchronization; never polls a person's session.
  const deadline=Date.now()+15000;
  for(;;){try{return await readFile(file,'utf8');}catch(error){if(error.code!=='ENOENT'||Date.now()>deadline)throw error;await new Promise(resolve=>setTimeout(resolve,20));}}
}
async function verifyArtifacts(result){
  assert.equal(result.artifacts.length,1);
  const item=result.artifacts[0],bytes=await readFile(item.path);
  assert.equal(item.key,'proof');assert.equal(bytes.toString(),'Fungi headless artifact\n');
  assert.equal(item.sha256,createHash('sha256').update(bytes).digest('hex'));
  assert.equal(item.bytes,bytes.length);
}
let pendingGate;
try{
  evidence.binarySHA256=createHash('sha256').update(await readFile(standalone)).digest('hex');
  evidence.version=await run(['--version']);assert.equal(evidence.version.protocol,2);
  evidence.doctor=await run(['doctor']);assert.equal(evidence.doctor.available,true);assert.equal(evidence.doctor.engine,'headless');
  assert.deepEqual(await readdir(state),[]);evidence.checks.push('standalone engine and storage/lock doctor');
  const index=await run(['describe']);assert.equal(index.engine,'headless');
  assert.ok(!index.topics.includes('html')&&!index.topics.includes('form'));
  for(const topic of index.documentation.topics)assert.ok((await invoke(['help',topic])).trim().length>0);
  evidence.checks.push('offline discovery and manuals');
  const input=path.join(directory,'input with spaces.json'),data={marker:'headless smoke',nested:{values:[null,true,1.5]},literal:'$(not a command)'};
  await writeFile(input,JSON.stringify(data));
  const launch=['open','--exec',process.execPath,'--input',input,'--delivery','wait','--',helper];
  const results=await Promise.all([run(launch),run(launch)]);
  assert.notEqual(results[0].session,results[1].session);
  for(const result of results){assert.equal(result.outcome,'submitted');assert.deepEqual(result.data.input,data);await verifyArtifacts(result);assert.deepEqual(await run(['collect',result.session]),result);}
  evidence.results=results;evidence.checks.push('concurrent external helpers, literal JSON, durable artifacts and collection');
  const gate=path.join(directory,'gate with spaces');pendingGate=gate;
  const receipt=await run(['open','--exec',process.execPath,'--input',input,'--delivery','manual','--',helper,gate]);
  assert.equal(receipt.delivery.mode,'manual');assert.equal(await fileReady(gate+'.ready'),receipt.session);
  await writeFile(gate,'closed');assert.equal(await fileReady(gate+'.done'),receipt.session);
  const closed=await run(['collect',receipt.session]);assert.equal(closed.outcome,'closed');
  const draft=JSON.parse(await readFile(path.join(state,receipt.session,'draft.json'),'utf8'));assert.equal(draft.data.marker,data.marker);
  await writeFile(gate,'submitted');
  const reopened=await run(['reopen',receipt.session,'--delivery','wait']);pendingGate=undefined;
  assert.notEqual(reopened.session,receipt.session);assert.deepEqual(reopened.data.restoredDraft,draft.data);await verifyArtifacts(reopened);
  assert.deepEqual(await run(['collect',receipt.session]),closed);
  evidence.recovery={receipt,closed,reopened};evidence.checks.push('manual completion and draft/artifact reopen without changing the original');
  for(const args of [['open','page.html'],['open','form.json'],['open','--exec',process.execPath,'--width','600']]){
    await assert.rejects(invoke(args),error=>error.code===2&&JSON.parse(error.stdout).error.code==='invalid_arguments');
  }
  evidence.checks.push('removed UI entry points rejected');evidence.passed=true;
}catch(error){evidence.passed=false;evidence.error={message:error.message,stdout:error.stdout,stderr:error.stderr};throw error;}
finally{
  if(pendingGate)await writeFile(pendingGate,'cancelled');
  await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
  console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),passed:evidence.passed,checks:evidence.checks}));
}
