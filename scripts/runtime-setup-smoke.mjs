// Linux-only integration: a disposable Ubuntu container starts without the
// webview libraries. Exercise the real bootstrap, including explicit setup.
import {createServer} from 'node:http';
import {createReadStream} from 'node:fs';
import {mkdtemp,mkdir,writeFile,readFile} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import path from 'node:path';
import assert from 'node:assert/strict';
const [archiveArg,target]=process.argv.slice(2);
assert.equal(process.platform,'linux');assert.ok(['ubuntu-22.04','ubuntu-24.04'].includes(target));
const archive=path.resolve(archiveArg),name=path.basename(archive);
await mkdir('verification-output',{recursive:true});
const directory=await mkdtemp(path.resolve('verification-output/runtime-setup-'));
const server=createServer((req,res)=>{
  const file=req.url===`/latest/download/${name}`?archive:req.url===`/latest/download/${name}.sha256`?archive+'.sha256':null;
  if(!file){res.writeHead(404);res.end();return;}createReadStream(file).pipe(res);
});
await new Promise(resolve=>server.listen(0,'127.0.0.1',resolve));
const shell=`set -eu
apt-get update > /proof/prerequisites.log 2>&1
apt-get install -y curl ca-certificates >> /proof/prerequisites.log 2>&1
mkdir /package
tar -xzf /archive.tar.gz -C /package
code=0
/package/fungi doctor > /proof/missing.json || code=$?
[ "$code" -eq 3 ]
printf '<!doctype html><title>Missing runtime check</title>' > /page.html
code=0
FUNGI_STATE_DIR=/sessions /package/fungi open /page.html --delivery wait > /proof/open-missing.json || code=$?
[ "$code" -eq 1 ]
grep -q '"kind":"result"' /proof/open-missing.json
grep -q 'runtime_missing' /proof/open-missing.json
code=0
sh /install.sh --dir /fungi > /proof/declined.json 2> /proof/declined.log || code=$?
[ "$code" -eq 3 ]
[ ! -f /fungi/installation.json ]
# Controlled package-manager failure: no actual system state is altered here.
mkdir /denied
printf '#!/bin/sh\\nexit 23\\n' > /denied/apt-get
chmod +x /denied/apt-get
code=0
PATH=/denied:$PATH sh /install.sh --setup-runtime --dir /fungi > /proof/failed.json 2> /proof/failed.log || code=$?
[ "$code" -eq 23 ]
[ ! -f /fungi/installation.json ]
sh /install.sh --setup-runtime --dir /fungi > /proof/installed.json 2> /proof/setup.log
/fungi/bin/fungi doctor > /proof/available.json
sh /install.sh --setup-runtime --dir /fungi > /proof/reused.json 2> /proof/reused.log
! grep -q 'Installing WebKitGTK' /proof/reused.log
`;
await writeFile(path.join(directory,'run.sh'),shell);
const evidence={scope:'real missing Linux libraries, ordinary launch failure, declined setup, injected package-manager failure, official Ubuntu runtime install and reuse; no desktop/UI claim',archive,target,directory};
try{
  const {stdout,stderr}=await promisify(execFile)('docker',['run','--rm','--network','host','-e',`FUNGI_RELEASE_BASE_URL=http://127.0.0.1:${server.address().port}`,'-v',`${archive}:/archive.tar.gz:ro`,'-v',`${path.resolve('scripts/install.sh')}:/install.sh:ro`,'-v',`${directory}:/proof`,`ubuntu:${target.replace('ubuntu-','')}`,'sh','/proof/run.sh'],{timeout:600000,maxBuffer:4*1024*1024});
  await writeFile(path.join(directory,'container.log'),stdout+'\n'+stderr);
  evidence.missing=JSON.parse(await readFile(path.join(directory,'missing.json'),'utf8'));
  evidence.available=JSON.parse(await readFile(path.join(directory,'available.json'),'utf8'));
  assert.equal(evidence.missing.error.code,'runtime_missing');assert.equal(evidence.available.available,true);
  evidence.passed=true;
}catch(error){evidence.passed=false;evidence.error={message:error.message,stdout:error.stdout,stderr:error.stderr};throw error;}
finally{await new Promise(resolve=>server.close(resolve));await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');console.log(JSON.stringify(evidence,null,2));}
