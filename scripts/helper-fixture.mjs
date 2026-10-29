// Automated protocol fixture only. No UI, model calls or human waiting.
import {createInterface} from 'node:readline';
import {readFile,writeFile,copyFile} from 'node:fs/promises';
import path from 'node:path';
import assert from 'node:assert/strict';

const reader=createInterface({input:process.stdin})[Symbol.asyncIterator]();
const receive=async()=>{const line=await reader.next();assert.equal(line.done,false,'owner disconnected');return JSON.parse(line.value);};
const send=value=>process.stdout.write(JSON.stringify(value)+'\n');
const persist=async value=>{send(value);const ack=await receive();assert.equal(ack.id,value.id);assert.equal(ack.ok,true,JSON.stringify(ack));return ack;};
const launch=await receive();
assert.equal(launch.protocol,2);assert.equal(process.env.CODEX_THREAD_ID,undefined);
assert.equal('source' in launch,false);assert.equal('width' in launch,false);
const data=launch.draft??{marker:launch.input?.marker??'fixture'};
const artifacts=[];
for(const [index,item] of launch.draftArtifacts.entries()){
  const name=`restored-${index}`;await copyFile(item.path,path.join(launch.outputDir,name));
  artifacts.push({key:item.key,name:item.name,path:name,mediaType:item.mediaType});
}
if(!artifacts.length){
  await writeFile(path.join(launch.outputDir,'proof.txt'),'Fungi headless artifact\n');
  artifacts.push({key:'proof',name:'proof.txt',path:'proof.txt',mediaType:'text/plain'});
}
await persist({id:'draft',type:'draft',data,artifacts});
send({type:'ready'});
const gate=process.argv[2];let outcome='submitted';
if(gate){
  await writeFile(gate+'.ready',launch.session);
  const deadline=Date.now()+60000;
  for(;;){
    try{outcome=(await readFile(gate,'utf8')).trim();break;}
    catch(error){if(error.code!=='ENOENT'||Date.now()>deadline)throw error;}
    await new Promise(resolve=>setTimeout(resolve,25));
  }
}
await persist({id:'result',type:'result',outcome,data:{input:launch.input,restoredDraft:launch.draft,notes:data},artifacts});
if(gate)await writeFile(gate+'.done',launch.session);
process.exit(0);
