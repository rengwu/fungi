// Real download/bootstrap and version switching, using test-owned local release
// files. No publication, model callback, persistent PATH change or human polling.
import {createServer} from 'node:http';
import {createReadStream} from 'node:fs';
import {mkdtemp,mkdir,readFile,writeFile,stat} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify,parseArgs} from 'node:util';
import {createHash} from 'node:crypto';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
import assert from 'node:assert/strict';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const {values,positionals}=parseArgs({allowPositionals:true,options:{'core-only':{type:'boolean'}}});
assert.equal(positionals.length,2,'usage: install-smoke.mjs ARCHIVE_A ARCHIVE_B [--core-only]');
const [archiveA,archiveB]=positionals.map(p=>path.resolve(p));
const versionA=path.basename(path.dirname(archiveA)).replace(/^release-/,''),versionB=path.basename(path.dirname(archiveB)).replace(/^release-/,'');
const checksumA=(await readFile(archiveA+'.sha256','utf8')).trim().split(/\s+/)[0];
const checksumB=(await readFile(archiveB+'.sha256','utf8')).trim().split(/\s+/)[0];
const parent=path.join(root,'verification-output');await mkdir(parent,{recursive:true});
const directory=await mkdtemp(path.join(parent,'install-smoke-'));
const installation=path.join(directory,'install with spaces'),state=path.join(directory,'sessions'),assets=path.join(directory,'assets');await mkdir(assets);
const binary=path.join(installation,'bin',process.platform==='win32'?'fungi.exe':'fungi');
const env={...process.env,FUNGI_STATE_DIR:state};
for(const key of Object.keys(env))if(['path','codex_thread_id','fungi_renderer_dir'].includes(key.toLowerCase()))delete env[key];
env.PATH='';
const run=async args=>JSON.parse((await promisify(execFile)(binary,args,{env,timeout:120000,maxBuffer:2**20})).stdout);
const requests=[];
const server=createServer(async(req,res)=>{
  requests.push(req.url);
  const base=`/download/${versionA}/${path.basename(archiveA)}`;
  const file=req.url===base?archiveA:req.url===base+'.sha256'?archiveA+'.sha256':null;
  if(!file){res.writeHead(404);res.end();return;}
  try{res.setHeader('Content-Length',(await stat(file)).size);createReadStream(file).pipe(res);}catch{res.writeHead(500);res.end();}
});
await new Promise(resolve=>server.listen(0,'127.0.0.1',resolve));
const bootstrapEnv={...process.env,FUNGI_RELEASE_BASE_URL:`http://127.0.0.1:${server.address().port}`,FUNGI_INSTALL_DIR:installation};
const evidence={directory,platform:process.platform,arch:process.arch,archiveA,checksumA,archiveB,checksumB,requests,scope:values['core-only']?'bootstrap and version switching; UI omitted':'bootstrap, active Electron session, draft/artifact, update and rollback',childPathEmpty:true};
async function untilFile(file){
  // A finite automated fixture wait, never an agent-facing wait for a person.
  const deadline=Date.now()+60000;
  while(Date.now()<deadline){try{return JSON.parse(await readFile(file,'utf8'));}catch(error){if(error.code!=='ENOENT')throw error;await new Promise(resolve=>setTimeout(resolve,25));}}
  throw new Error('Fixture timed out: '+file);
}
try{
  let installer;
  const bootstrapStarted=Date.now();
  if(process.platform==='win32'){
    const powershell=path.join(process.env.SystemRoot,'System32','WindowsPowerShell','v1.0','powershell.exe');
    // Allow hosted Windows ZIP extraction headroom. This is a finite fixture
    // bound, not a timeout for a human interaction or a relaxed assertion.
    installer=await promisify(execFile)(powershell,['-NoProfile','-File',path.join(root,'scripts/install.ps1'),'-Version',versionA,'-Directory',installation,'-NoPath'],{env:bootstrapEnv,timeout:300000});
  }else{
    bootstrapEnv.PATH='/usr/bin:/bin:/usr/sbin:/sbin';
    installer=await promisify(execFile)('/bin/sh',[path.join(root,'scripts/install.sh'),'--version',versionA,'--dir',installation],{env:bootstrapEnv,timeout:180000});
  }
  evidence.bootstrap={stdout:installer.stdout,stderr:installer.stderr,elapsedMs:Date.now()-bootstrapStarted};
  evidence.before=await run(['--version']);assert.equal(evidence.before.version,versionA);
  assert.deepEqual(requests,[`/download/${versionA}/${path.basename(archiveA)}.sha256`,`/download/${versionA}/${path.basename(archiveA)}`]);
  let receipt;
  if(!values['core-only']){
    await writeFile(path.join(assets,'index.html'),`<!doctype html><meta charset="utf-8"><title>Fungi update verification</title>
      <h1>Update verification</h1><p>Synthetic fixture. This window submits after the test updates Fungi.</p>
      <script type="module">
      const initial=await fungi.ready();
      const canvas=document.createElement('canvas');canvas.width=canvas.height=32;document.body.append(canvas);
      const ctx=canvas.getContext('2d');ctx.fillStyle='#865ce6';ctx.fillRect(0,0,32,32);
      const blob=await new Promise(resolve=>canvas.toBlob(resolve));
      const data={notes:'Saved before update',restored:!!initial.draft};
      await fungi.saveDraft(data,[{key:'canvas',name:'canvas.png',blob}]);
      const deadline=Date.now()+90000;
      while(Date.now()<deadline){
        const response=await fetch('./gate.json',{cache:'no-store'});
        if(response.ok){const gate=await response.json();if(gate.submit){await fungi.submit({...data,activated:gate.version},[{key:'canvas',name:'canvas.png',blob}]);break;}}
        await new Promise(resolve=>setTimeout(resolve,50));
      }
      </script>`);
    receipt=await run(['open',path.join(assets,'index.html'),'--delivery','manual']);
    evidence.receipt=receipt;
    evidence.draft=await untilFile(path.join(state,receipt.session,'draft.json'));
    assert.equal(evidence.draft.data.notes,'Saved before update');
    assert.equal(path.basename(path.dirname(receipt.collect[0])),'bin');
  }
  evidence.update=await run(['update','--archive',archiveB,'--sha256',checksumB]);
  evidence.after=await run(['--version']);assert.equal(evidence.after.version,versionB);
  assert.ok(evidence.update.reusedRuntimeBytes>0,'matching runtime files were not shared');
  const engine=process.platform==='darwin'?'Electron.app/Contents/MacOS/Electron':process.platform==='win32'?'electron.exe':'electron';
  const runtimeFile=version=>path.join(installation,'versions',version,'renderer','runtime',engine);
  const [oldEngine,newEngine]=await Promise.all([stat(runtimeFile(versionA),{bigint:true}),stat(runtimeFile(versionB),{bigint:true})]);
  assert.equal(oldEngine.dev,newEngine.dev);assert.equal(oldEngine.ino,newEngine.ino);
  evidence.runtimeShared={bytes:evidence.update.reusedRuntimeBytes,engine,sameFile:true};
  if(receipt){
    await writeFile(path.join(assets,'gate.json'),JSON.stringify({submit:true,version:versionB}));
    const saved=await untilFile(path.join(state,receipt.session,'result.json'));
    const result=JSON.parse((await promisify(execFile)(receipt.collect[0],receipt.collect.slice(1),{env:{...env,...receipt.collectEnv},timeout:30000})).stdout);
    assert.deepEqual(result,saved);assert.equal(result.outcome,'submitted');assert.equal(result.data.activated,versionB);
    const bytes=await readFile(result.artifacts[0].path);assert.equal(createHash('sha256').update(bytes).digest('hex'),result.artifacts[0].sha256);
    evidence.result=result;evidence.artifactVerified=true;
  }
  evidence.rollback=await run(['rollback']);assert.equal((await run(['--version'])).version,versionA);
  evidence.passed=true;
}catch(error){evidence.passed=false;evidence.error={message:error.message,code:error.code,signal:error.signal,killed:error.killed,stdout:error.stdout,stderr:error.stderr};throw error;}
finally{
  await new Promise(resolve=>server.close(resolve));
  await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
  console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),...evidence},null,2));
}
