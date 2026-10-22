import {mkdtemp,writeFile,mkdir,symlink,rm,realpath} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import path from 'node:path';
import test from 'node:test';
import assert from 'node:assert/strict';
import assets from '../renderer/assets.cjs';

test('local assets stream, seek and reject escapes or special routes',async()=>{
  const dir=await mkdtemp(path.join(tmpdir(),'fungi-assets-'));
  try{
    await mkdir(path.join(dir,'skill'));const root=await realpath(path.join(dir,'skill'));
    await writeFile(path.join(root,'sample.wav'),'0123456789');
    await writeFile(path.join(root,'empty.txt'),'');
    await writeFile(path.join(dir,'private.txt'),'private');
    await symlink(path.join(dir,'private.txt'),path.join(root,'escape.txt'));
    const get=(name,init={})=>assets.serveAsset(new Request('fungi://skill/'+name,init),{skill:root});
    let response=await get('sample.wav');
    assert.equal(response.headers.get('content-type'),'audio/wav');
    assert.equal(await response.text(),'0123456789');
    for(const [range,expected] of [['bytes=2-5','2345'],['bytes=7-','789'],['bytes=-3','789'],['bytes=8-50','89']]){
      response=await get('sample.wav',{headers:{Range:range}});
      assert.equal(response.status,206);assert.equal(await response.text(),expected);
    }
    for(const range of ['bytes=99-','bytes=5-2','bytes=-0','bytes=1-2,4-5','bytes=','bytes=99999999999999999-']){
      response=await get('sample.wav',{headers:{Range:range}});assert.equal(response.status,416,range);
      assert.equal(response.headers.get('content-range'),'bytes */10');
    }
    response=await get('sample.wav',{method:'HEAD'});
    assert.equal(response.headers.get('content-length'),'10');assert.equal(await response.text(),'');
    assert.equal((await get('sample.wav',{method:'POST'})).status,405);
    assert.equal(await (await get('empty.txt')).text(),'');
    assert.equal((await get('escape.txt')).status,403);
    assert.equal((await get('%2e%2e%2fprivate.txt')).status,403);
    assert.equal((await get('missing')).status,404);
    assert.equal((await get('')).status,404);
    assert.equal((await assets.serveAsset(new Request('fungi://constructor/sample.wav'),{skill:root})).status,404);
  }finally{await rm(dir,{recursive:true,force:true});}
});
