import {readFile} from 'node:fs/promises';
import vm from 'node:vm';
import test from 'node:test';
import assert from 'node:assert/strict';

const source=await readFile(new URL('../renderer/wails/bridge.js',import.meta.url),'utf8');
function bridge(invoke) {
  const window={go:{main:{Bridge:{Ready:()=>invoke('ready'),Persist:(type,data)=>invoke(type,data)}}},runtime:{Quit(){}}};
  window.top=window;
  vm.runInNewContext(source,{window,Blob,Uint8Array,btoa,atob,setTimeout(){}});
  return window.fungi;
}
const ready={attachments:[],limits:{jsonBytes:1048576,artifactBytes:20971520,totalArtifactBytes:104857600}};

test('slow attachment encoding cannot reorder saves, submit or cancel',async()=>{
  const calls=[];
  const api=bridge(async(channel,data)=>{
    if(channel==='ready')return structuredClone(ready);
    calls.push([channel,data?.version]);
  });
  let release;
  class SlowBlob extends Blob {
    async arrayBuffer(){await new Promise(resolve=>release=resolve);return super.arrayBuffer();}
  }
  const first=api.saveDraft({version:1},[{key:'image',name:'image',blob:new SlowBlob(['first'])}]);
  const second=api.saveDraft({version:2});
  const third=api.submit({version:3});
  const fourth=api.cancel();
  await new Promise(resolve=>setImmediate(resolve));
  assert.deepEqual(calls,[],'later mutations overtook encoding');
  release();
  await Promise.all([first,second,third,fourth]);
  assert.deepEqual(calls,[['draft',1],['draft',2],['submitted',3],['cancelled',undefined]]);
});

test('rejected save does not poison a later submit retry',async()=>{
  const api=bridge(async(channel)=>{
    if(channel==='ready')return structuredClone(ready);
    if(channel==='draft')throw new Error('disk full');
    return {outcome:'submitted'};
  });
  await assert.rejects(api.saveDraft({notes:'keep me'}),/disk full/);
  assert.equal((await api.submit({notes:'keep me'})).outcome,'submitted');
});
