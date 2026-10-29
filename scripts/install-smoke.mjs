// Real download/bootstrap and version switching, using test-owned local release
// files. No publication, model callback, persistent PATH change or human polling.
import {createServer} from 'node:http';
import {createReadStream} from 'node:fs';
import {mkdtemp,mkdir,readFile,writeFile,stat} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {createHash} from 'node:crypto';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
import assert from 'node:assert/strict';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const positionals=process.argv.slice(2);
assert.equal(positionals.length,2,'usage: install-smoke.mjs ARCHIVE_A ARCHIVE_B');
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
const invoke=async args=>(await promisify(execFile)(binary,args,{cwd:directory,env,timeout:120000,maxBuffer:2**20})).stdout;
const run=async args=>JSON.parse(await invoke(args));
async function discovery(version){
  const base=path.join(installation,'versions',version);
  const index=await run(['describe']);
  for(const topic of index.documentation.topics){
    const help=await invoke(['help',topic]);
    const markdown=await readFile(path.join(base,'docs/manual',topic+'.md'),'utf8');
    assert.ok(help.startsWith(markdown.split('\n')[0].replace(/^# /,'')+'\n'));
    assert.doesNotMatch(help,/^#{1,6} |^\| |^```/m);
    for(const [index,block] of markdown.split('```').entries()){
      if(index%2===0)continue;
      const code=block.slice(block.indexOf('\n')+1).replace(/\n$/,'');
      for(const line of code.split('\n'))assert.ok(help.includes('  '+line+'\n'),'help changed an example line');
    }
  }
  const result=await run(['describe','examples']);
  assert.equal(index.protocol,2);assert.equal(index.engine,'headless');
  assert.equal(result.capability.guide,'fungi help examples');
  return {version,topics:index.documentation.topics};
}
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
const evidence={directory,platform:process.platform,arch:process.arch,archiveA,checksumA,archiveB,checksumB,requests,scope:'bootstrap, active external helper, draft/artifact, update and rollback',childPathEmpty:true};
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
  evidence.discoveryBefore=await discovery(versionA);
  assert.deepEqual(requests,[`/download/${versionA}/${path.basename(archiveA)}.sha256`,`/download/${versionA}/${path.basename(archiveA)}`]);
  const gate=path.join(assets,'gate');
  const receipt=await run(['open','--exec',process.execPath,'--delivery','manual','--',path.join(root,'scripts/helper-fixture.mjs'),gate]);
  evidence.receipt=receipt;
  evidence.draft=await untilFile(path.join(state,receipt.session,'draft.json'));
  assert.equal(evidence.draft.data.marker,'fixture');
  assert.equal(path.basename(path.dirname(receipt.collect[0])),'bin');
  evidence.update=await run(['update','--archive',archiveB,'--sha256',checksumB]);
  evidence.after=await run(['--version']);assert.equal(evidence.after.version,versionB);
  evidence.discoveryAfter=await discovery(versionB);
  evidence.engine=await run(['doctor']);assert.equal(evidence.engine.available,true);
  const manifest=JSON.parse(await readFile(path.join(installation,'versions',versionB,'RELEASE.json'),'utf8'));
  assert.equal(manifest.format,3);assert.equal(manifest.protocol,2);assert.ok(!manifest.wails&&!manifest.electron&&!manifest.rendererSHA256);
  assert.equal(await stat(path.join(installation,'versions',versionB,'renderer')).then(()=>true,()=>false),false,'engine must not contain a renderer');
  evidence.manifest=manifest;
  if(receipt){
    await writeFile(gate,'submitted');
    const saved=await untilFile(path.join(state,receipt.session,'result.json'));
    const result=JSON.parse((await promisify(execFile)(receipt.collect[0],receipt.collect.slice(1),{env:{...env,...receipt.collectEnv},timeout:30000})).stdout);
    assert.deepEqual(result,saved);assert.equal(result.outcome,'submitted');assert.equal(result.data.notes.marker,'fixture');
    const bytes=await readFile(result.artifacts[0].path);assert.equal(createHash('sha256').update(bytes).digest('hex'),result.artifacts[0].sha256);
    evidence.result=result;evidence.artifactVerified=true;
  }
  evidence.rollback=await run(['rollback']);assert.equal((await run(['--version'])).version,versionA);
  evidence.discoveryRollback=await discovery(versionA);
  evidence.passed=true;
}catch(error){evidence.passed=false;evidence.error={message:error.message,code:error.code,signal:error.signal,killed:error.killed,stdout:error.stdout,stderr:error.stderr};throw error;}
finally{
  await writeFile(path.join(assets,'gate'),'cancelled');
  await new Promise(resolve=>server.close(resolve));
  await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
  console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),...evidence},null,2));
}
