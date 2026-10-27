import {readFile} from 'node:fs/promises';
import vm from 'node:vm';
import test from 'node:test';
import assert from 'node:assert/strict';

const source=await readFile(new URL('../renderer/wails/bridge.js',import.meta.url),'utf8');
function bridge(invoke) {
  const window={go:{main:{Bridge:{Ready:()=>invoke('ready'),Persist:(type,data,items)=>invoke(type,data,items)}}},runtime:{Quit(){}}};
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

test('protocol 1 readiness restores binary files and tolerates additive metadata',async()=>{
  let reads=0;
  const api=bridge(async()=>{
    reads++;
    return {...structuredClone(ready),input:{enabled:false,count:0},draft:null,futureMetadata:{optional:true},
      delivery:{mode:'manual'},attachments:[
        {key:'empty',name:'empty.bin',mediaType:'application/octet-stream'},
        {key:'bytes',name:'bytes.bin',mediaType:'application/octet-stream',bytes:'AP+A'}
      ]};
  });
  for(const name of ['ready','saveDraft','submit','cancel'])assert.equal(typeof api[name],'function',name);
  for(let i=0;i<2;i++){
    const state=await api.ready();
    assert.deepEqual(state.input,{enabled:false,count:0});
    assert.equal(state.draft,null);
    assert.equal(state.delivery.mode,'manual');
    assert.equal(state.attachments[0].blob.size,0);
    assert.equal(state.attachments[1].key,'bytes');
    assert.equal(state.attachments[1].name,'bytes.bin');
    assert.equal(state.attachments[1].blob.type,'application/octet-stream');
    assert.deepEqual(new Uint8Array(await state.attachments[1].blob.arrayBuffer()),new Uint8Array([0,255,128]));
  }
  assert.equal(reads,1,'readiness must be reusable');
});

test('protocol 1 submission preserves JSON values, attachment bytes and acknowledged result',async()=>{
  for(const data of [null,false,0,'',[],{nested:[false,0,null,'雪']}]){
    const expected={protocol:1,kind:'result',outcome:'submitted',data,artifacts:[],futureMetadata:true};
    const api=bridge(async(channel,payload,items)=>{
      if(channel==='ready')return structuredClone(ready);
      assert.equal(channel,'submitted');
      assert.deepEqual(payload,data);
      assert.deepEqual(JSON.parse(JSON.stringify(items)),[
        {key:'bytes',name:'sample.bin',mediaType:'application/octet-stream',bytes:'AP+A'}
      ]);
      return expected;
    });
    const result=await api.submit(data,[{key:'bytes',name:'sample.bin',blob:new Blob([new Uint8Array([0,255,128])],{type:'application/octet-stream'})}]);
    assert.deepEqual(result,expected);
  }
});
