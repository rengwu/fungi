// Build-time tools only. Packages contain one Go engine executable and docs.
import {cp,mkdir,readFile,writeFile} from 'node:fs/promises';
import {execFileSync} from 'node:child_process';
import {parseArgs} from 'node:util';
import {createHash} from 'node:crypto';
import path from 'node:path';
import {fileURLToPath} from 'node:url';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const {values}=parseArgs({options:{platform:{type:'string',default:process.platform},arch:{type:'string',default:process.arch},release:{type:'string'}}});
const goos={darwin:'darwin',win32:'windows',linux:'linux'}[values.platform],goarch={arm64:'arm64',x64:'amd64'}[values.arch];
if(!goos||!goarch)throw new Error('Targets: darwin, win32 or linux; arm64 or x64');
if(values.release&&!/^v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9]+(?:[.-][A-Za-z0-9]+)*)?$/.test(values.release))throw new Error('--release requires vX.Y.Z, optionally with a prerelease suffix');
const stamp=new Date().toISOString().replace(/[:.]/g,'-');
const dest=values.release?path.join(root,'dist',`release-${values.release}`,`fungi-${goos}-${goarch}`):path.join(root,'dist',`fungi-${goos}-${goarch}-${stamp}`);
await mkdir(path.dirname(dest),{recursive:true});await mkdir(dest);
const binary=path.join(dest,goos==='windows'?'fungi.exe':'fungi');
const version=values.release??'0.2.0-dev';
execFileSync('go',['build','-trimpath','-ldflags',`-s -w -X main.buildVersion=${version}`,'-o',binary,'.'],{cwd:root,stdio:'inherit',env:{...process.env,GOOS:goos,GOARCH:goarch,CGO_ENABLED:'0'}});
const goRoot=execFileSync('go',['env','GOROOT'],{encoding:'utf8'}).trim();
await mkdir(path.join(dest,'licenses','go'),{recursive:true});
try{await cp(path.join(goRoot,'LICENSE'),path.join(dest,'licenses','go','LICENSE'));}
catch(error){if(error.code!=='ENOENT')throw error;await cp(path.join(goRoot,'..','LICENSE'),path.join(dest,'licenses','go','LICENSE'));}
await cp(path.join(root,'docs/manual'),path.join(dest,'docs/manual'),{recursive:true});
await cp(path.join(root,'docs/PORTABLE.md'),path.join(dest,'README.md'));
const identity={format:3,protocol:2,version,platform:goos,arch:goarch,binarySHA256:createHash('sha256').update(await readFile(binary)).digest('hex')};
const sourceCommit=execFileSync('git',['rev-parse','HEAD'],{cwd:root,encoding:'utf8'}).trim();
const sourceDirty=execFileSync('git',['status','--porcelain'],{cwd:root,encoding:'utf8'}).trim()!=='';
await writeFile(path.join(dest,'BUILD.json'),JSON.stringify({...identity,created:stamp,sourceCommit,sourceDirty},null,2)+'\n');
if(values.release){
  await writeFile(path.join(dest,'RELEASE.json'),JSON.stringify(identity,null,2)+'\n');
  const archive=dest+(goos==='windows'?'.zip':'.tar.gz');
  execFileSync('go',['run','./scripts/archive',dest,archive],{cwd:root,stdio:'inherit'});
  console.log(archive);
}else console.log(dest);
