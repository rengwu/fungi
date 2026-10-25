// Diagnostics for the accepted Wails host behavior.
// Every page and localStorage key belongs to this disposable run.
import {mkdtemp,mkdir,writeFile,readFile} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import {randomUUID} from 'node:crypto';
import path from 'node:path';
import assert from 'node:assert/strict';

const [corePath,helperPath]=process.argv.slice(2);
assert.ok(corePath&&helperPath,'usage: node scripts/wails-probe.mjs CORE HELPER');
await mkdir('verification-output',{recursive:true});
const directory=await mkdtemp(path.resolve('verification-output/wails-probe-'));
const env={...process.env,FUNGI_STATE_DIR:path.join(directory,'sessions')};
delete env.CODEX_THREAD_ID;
if(process.platform==='linux'){
  env.GST_DEBUG='2,webkitmediaplayer:6,webkitwebsrc:6';
  env.GST_DEBUG_NO_COLOR='1';
}
const core=path.resolve(corePath),helper=path.resolve(helperPath),entry=path.join(directory,'index.html');
const key='fungi-disposable-'+randomUUID();
const wave=Buffer.alloc(16044);wave.write('RIFF');wave.writeUInt32LE(wave.length-8,4);wave.write('WAVEfmt ',8);wave.writeUInt32LE(16,16);
wave.writeUInt16LE(1,20);wave.writeUInt16LE(1,22);wave.writeUInt32LE(8000,24);wave.writeUInt32LE(16000,28);wave.writeUInt16LE(2,32);wave.writeUInt16LE(16,34);wave.write('data',36);wave.writeUInt32LE(16000,40);
await writeFile(path.join(directory,'sample.wav'),wave);
await writeFile(entry,`<!doctype html><meta charset="utf-8"><title>Disposable Wails diagnostics</title>
<h1>Disposable Wails diagnostics</h1><script type="module">
const ready=await fungi.ready(),key=${JSON.stringify(key)};
if(ready.input?.independent){
  const inheritedStorage=localStorage.getItem(key);
  localStorage.removeItem(key);
  await fungi.submit({inheritedStorage});
}else if(ready.draft){
  await fungi.submit({restoredEmpty:ready.attachments.length===1&&ready.attachments[0].blob.size===0,
    restoredType:ready.attachments[0]?.blob.type,diagnostics:ready.draft},ready.attachments);
}else{
  const native={browserOpenURL:typeof window.runtime?.BrowserOpenURL,
    clipboardGetText:typeof window.runtime?.ClipboardGetText,
    quit:typeof window.runtime?.Quit,boundPersist:typeof window.go?.main?.Bridge?.Persist,
    lowLevelInvoke:typeof window.external?.invoke};
  async function media(url){
    const audio=new Audio(url);audio.preload='auto';let timer;
    try{return await new Promise(resolve=>{
      timer=setTimeout(()=>resolve({ok:false,error:'timeout'}),3000);
      audio.onloadedmetadata=()=>resolve({ok:true,duration:audio.duration});
      audio.onerror=()=>resolve({ok:false,code:audio.error?.code,error:audio.error?.message});
    });}finally{clearTimeout(timer);audio.src='';}
  }
  const response=await fetch('./sample.wav'),blob=await response.blob(),url=URL.createObjectURL(blob);
  const localMedia=await media('./sample.wav'),blobMedia=await media(url);URL.revokeObjectURL(url);
  localStorage.setItem(key,'written-by-another-fungi-session');
  await fungi.saveDraft({native,localMedia,blobMedia,mediaType:blob.type},
    [{key:'empty',name:'empty.txt',blob:new Blob([],{type:'text/plain'})}]);
  window.runtime.Quit();
}
</script>`);
const run=async args=>JSON.parse((await promisify(execFile)(core,args,{env,timeout:30000,maxBuffer:1024*1024})).stdout);
const closed=await run(['open','--exec',helper,'--delivery','wait','--',entry]);
assert.equal(closed.outcome,'closed');
const restored=await run(['reopen',closed.session,'--delivery','wait']);
assert.equal(restored.outcome,'submitted');
assert.equal(restored.data.restoredEmpty,true);
assert.equal(restored.data.restoredType,'text/plain');
assert.equal((await readFile(restored.artifacts[0].path)).length,0);
const input=path.join(directory,'input.json');await writeFile(input,'{"independent":true}');
const independent=await run(['open','--exec',helper,'--input',input,'--delivery','wait','--',entry]);
const evidence={core,helper,directory,closed,restored,independent,
  scope:'native close hook, empty artifact round trip, media and cross-session storage diagnostics'};
await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),restored:restored.data,independent:independent.data},null,2));
