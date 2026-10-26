// Verify the actual public installer and updater without GitHub credentials,
// changing the user's PATH, or installing system runtimes. Use a ready desktop.
import {mkdtemp,mkdir,readFile,writeFile,rm} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {createHash} from 'node:crypto';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
import assert from 'node:assert/strict';

const version=process.argv[2];
assert.match(version??'',/^v\d+\.\d+\.\d+-alpha\.[1-9]\d*$/,'usage: public-release-smoke.mjs vX.Y.Z-alpha.N');
const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const parent=path.join(root,'verification-output');await mkdir(parent,{recursive:true});
const directory=await mkdtemp(path.join(parent,'public-release-'));
const install=path.join(directory,'install with spaces');
const env={...process.env,FUNGI_STATE_DIR:path.join(directory,'sessions')};
for(const key of Object.keys(env))if(['gh_token','github_token','codex_thread_id','fungi_renderer_dir','fungi_release_base_url','fungi_install_dir'].includes(key.toLowerCase()))delete env[key];
const execute=promisify(execFile);
const evidence={version,platform:process.platform,arch:process.arch,anonymous:true,persistentPathChanged:false,runtimeSetup:false};
try{
  const name=process.platform==='win32'?'install.ps1':'install.sh';
  const url=`https://github.com/rengwu/fungi/releases/download/${version}/${name}`;
  const response=await fetch(url,{signal:AbortSignal.timeout(120000)});
  assert.equal(response.status,200,url);
  const script=Buffer.from(await response.arrayBuffer());
  const checksum=await fetch(url+'.sha256',{signal:AbortSignal.timeout(120000)});
  assert.equal(checksum.status,200);
  const [digest,filename]=(await checksum.text()).trim().split(/\s+/);
  assert.equal(filename,name);assert.equal(createHash('sha256').update(script).digest('hex'),digest);
  const installer=path.join(directory,name);await writeFile(installer,script);
  let result;
  if(process.platform==='win32'){
    const powershell=path.join(process.env.SystemRoot,'System32','WindowsPowerShell','v1.0','powershell.exe');
    result=await execute(powershell,['-NoProfile','-File',installer,'-Version',version,'-Directory',install,'-NoPath'],{env,timeout:600000});
  }else{
    result=await execute('/bin/sh',[installer,'--version',version,'--dir',install],{env,timeout:600000});
  }
  evidence.installer={url,sha256:digest,stdout:result.stdout,stderr:result.stderr};
  const binary=path.join(install,'bin',process.platform==='win32'?'fungi.exe':'fungi');
  const run=async args=>JSON.parse((await execute(binary,args,{env,timeout:600000})).stdout);
  evidence.identity=await run(['--version']);assert.equal(evidence.identity.version,version);
  evidence.build=JSON.parse(await readFile(path.join(install,'versions',version,'BUILD.json'),'utf8'));
  assert.equal(evidence.build.sourceDirty,false);assert.match(evidence.build.sourceCommit,/^[a-f0-9]{40}$/);
  evidence.doctor=await run(['doctor']);
  evidence.update=await run(['update','--version',version]);
  assert.equal((await run(['--version'])).version,version);
  evidence.examples=await run(['describe','examples']);
  assert.ok(evidence.examples.capability.items.every(item=>item.available));
  const args=[path.join(root,'scripts/renderer-smoke.mjs'),binary];
  if(process.platform==='linux')args.push('--allow-unavailable-webgl','--allow-unavailable-local-media');
  evidence.renderer=JSON.parse((await execute(process.execPath,args,{env,timeout:180000,maxBuffer:2**20})).stdout);
  evidence.passed=true;
}catch(error){evidence.error=error.stack;evidence.stdout=error.stdout;evidence.stderr=error.stderr;throw error;}
finally{
  await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
  await rm(install,{recursive:true,force:true});
  console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),passed:evidence.passed===true}));
}
