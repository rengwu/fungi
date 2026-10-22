// Development-only packaging. The resulting directory requires neither Node
// nor Go on the user's machine. Uses the downloaded, pinned Electron runtime.
import { cp, mkdir, readFile, writeFile, stat } from 'node:fs/promises';
import { execFileSync } from 'node:child_process';
import { createRequire } from 'node:module';
import { parseArgs } from 'node:util';
import { createHash } from 'node:crypto';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const {values}=parseArgs({options:{platform:{type:'string',default:process.platform},arch:{type:'string',default:process.arch}}});
const {platform,arch}=values,goos={darwin:'darwin',win32:'windows',linux:'linux'}[platform],goarch={arm64:'arm64',x64:'amd64'}[arch];
if(!goos||!goarch)throw new Error('Targets: darwin, win32 or linux; arm64 or x64');
const pkg=JSON.parse(await readFile(path.join(root,'renderer/package.json'),'utf8'));
const version=pkg.devDependencies.electron;
const stamp=new Date().toISOString().replace(/[:.]/g,'-');
const dest=path.join(root,'dist',`fungi-${platform}-${arch}-${stamp}`);
const renderer=path.join(dest,'renderer');
await mkdir(renderer,{recursive:true});
const env={...process.env,GOOS:goos,GOARCH:goarch,CGO_ENABLED:'0'},suffix=platform==='win32'?'.exe':'';
const binary=path.join(dest,'fungi'+suffix);
execFileSync('go',['build','-trimpath','-o',binary,'.'],{cwd:root,stdio:'inherit',env});
const runtime=path.join(renderer,'runtime');
let runtimeArchiveSHA256=null;
if(platform===process.platform&&arch===process.arch){
  await cp(path.join(root,'renderer/node_modules/electron/dist'),runtime,{recursive:true,verbatimSymlinks:true});
}else{
  // Reuse Electron's pinned build-time downloader/extractor and its published
  // npm checksums. Cross-packaging does not establish target runtime support.
  const require=createRequire(path.join(root,'renderer/package.json'));
  const {downloadArtifact}=require('@electron/get'),{extract}=require('@electron-internal/extract-zip');
  const archive=await downloadArtifact({version,artifactName:'electron',platform,arch,checksums:require('electron/checksums.json')});
  runtimeArchiveSHA256=createHash('sha256').update(await readFile(archive)).digest('hex');
  await extract(archive,{dir:runtime});
}
const runtimeVersion=(await readFile(path.join(runtime,'version'),'utf8')).trim().replace(/^v/,'');
if(runtimeVersion!==version)throw new Error(`Runtime version ${runtimeVersion} differs from pinned ${version}`);
await stat(path.join(runtime,platform==='darwin'?'Electron.app/Contents/MacOS/Electron':platform==='win32'?'electron.exe':'electron'));
for(const file of ['package.json','main.cjs','assets.cjs','preload.cjs','form.html','form.js'])await cp(path.join(root,'renderer',file),path.join(renderer,file));
// Keep upstream signatures untouched; this development archive is not a
// notarized public release. Its runtime and app code have separate locations.
await mkdir(path.join(dest,'examples'));
for(const file of ['feedback.json','review.html'])await cp(path.join(root,'examples',file),path.join(dest,'examples',file));
await cp(path.join(root,'examples/jstack-grill-form'),path.join(dest,'examples/jstack-grill-form'),{recursive:true});
await cp(path.join(root,'examples/compare'),path.join(dest,'examples/compare'),{recursive:true});
execFileSync('go',['build','-trimpath','-o',path.join(dest,'examples/helper'+suffix),'./examples/helper'],{cwd:root,stdio:'inherit',env});
await cp(path.join(root,'README.md'),path.join(dest,'README.md'));
await cp(path.join(root,'docs/wiki/verification/desktop.md'),path.join(dest,'VERIFY.md'));
await cp(path.join(root,'skills'),path.join(dest,'skills'),{recursive:true,verbatimSymlinks:true});
await cp(path.join(root,'docs/wiki'),path.join(dest,'docs/wiki'),{recursive:true,verbatimSymlinks:true});
await writeFile(path.join(dest,'BUILD.json'),JSON.stringify({created:stamp,platform,arch,electron:version,binarySHA256:createHash('sha256').update(await readFile(binary)).digest('hex'),runtimeArchiveSHA256,publicRelease:false,runtimeVerified:false},null,2)+'\n');
console.log(dest);
