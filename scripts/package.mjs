// Build-time tools only. Each package contains a Go CLI and a native Wails
// helper; the OS supplies the web engine. Build on the target OS/architecture.
import {cp,mkdir,readFile,writeFile,readdir} from 'node:fs/promises';
import {execFileSync} from 'node:child_process';
import {parseArgs} from 'node:util';
import {createHash} from 'node:crypto';
import path from 'node:path';
import {fileURLToPath} from 'node:url';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const {values}=parseArgs({options:{platform:{type:'string',default:process.platform},arch:{type:'string',default:process.arch},release:{type:'string'}}});
const {platform,arch}=values,goos={darwin:'darwin',win32:'windows',linux:'linux'}[platform],goarch={arm64:'arm64',x64:'amd64'}[arch];
if(!goos||!goarch)throw new Error('Targets: darwin, win32 or linux; arm64 or x64');
if(platform!==process.platform||arch!==process.arch)throw new Error('Build Wails packages on the target OS/CPU; use the platform CI builders.');
if(values.release&&!/^v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9]+(?:[.-][A-Za-z0-9]+)*)?$/.test(values.release))throw new Error('--release requires vX.Y.Z, optionally with a prerelease suffix');
const stamp=new Date().toISOString().replace(/[:.]/g,'-');
const dest=values.release?path.join(root,'dist',`release-${values.release}`,`fungi-${goos}-${goarch}`):path.join(root,'dist',`fungi-${goos}-${goarch}-${stamp}`);
await mkdir(path.dirname(dest),{recursive:true});await mkdir(dest);
const suffix=platform==='win32'?'.exe':'',binary=path.join(dest,'fungi'+suffix);
const version=values.release??'0.1.0-dev';
execFileSync('go',['build','-trimpath','-ldflags',`-s -w -X main.buildVersion=${version}`,'-o',binary,'.'],{cwd:root,stdio:'inherit',env:{...process.env,CGO_ENABLED:'0'}});
const renderer=path.join(dest,'renderer');
const helper=path.join(renderer,platform==='darwin'?'Fungi.app/Contents/MacOS/fungi-renderer':'fungi-renderer'+suffix);
await mkdir(path.dirname(helper),{recursive:true});
const tags=platform==='linux'?'production,webkit2_41':platform==='win32'?'production,wv2runtime.error':'production';
execFileSync('go',['build','-tags',tags,'-trimpath','-ldflags','-s -w','-o',helper,'./wails'],{cwd:path.join(root,'renderer'),stdio:'inherit',env:{...process.env,CGO_ENABLED:platform==='win32'?'0':'1'}});
if(platform==='darwin'){
  await writeFile(path.join(renderer,'Fungi.app/Contents/Info.plist'),`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict><key>CFBundleExecutable</key><string>fungi-renderer</string><key>CFBundleIdentifier</key><string>dev.fungi.renderer</string><key>CFBundleName</key><string>Fungi</string><key>CFBundlePackageType</key><string>APPL</string><key>CFBundleShortVersionString</key><string>${version.replace(/^v/,'').split('-')[0]}</string><key>CFBundleVersion</key><string>1</string><key>NSHighResolutionCapable</key><true/></dict></plist>\n`);
}
const modules=execFileSync('go',['list','-m','-f','{{.Path}}|{{.Version}}|{{.Dir}}','all'],{cwd:path.join(root,'renderer'),encoding:'utf8'}).trim().split('\n');
const wails=modules.find(line=>line.startsWith('github.com/wailsapp/wails/v2|')).split('|')[1];
const goRoot=execFileSync('go',['env','GOROOT'],{encoding:'utf8'}).trim();
await mkdir(path.join(dest,'licenses','go'),{recursive:true});
await cp(path.join(goRoot,'LICENSE'),path.join(dest,'licenses','go','LICENSE'));
for(const line of modules){
  const [name,,directory]=line.split('|');if(!directory||name==='fungi/renderer')continue;
  const notices=(await readdir(directory)).filter(name=>/^(licen[cs]e|copying|notice|authors)(\.|$)/i.test(name));
  if(!notices.length)continue;
  const destination=path.join(dest,'licenses',name.replaceAll('/','_'));await mkdir(destination,{recursive:true});
  for(const notice of notices)await cp(path.join(directory,notice),path.join(destination,notice),{recursive:true});
}
await mkdir(path.join(dest,'examples'));
for(const file of ['feedback.json','review.html'])await cp(path.join(root,'examples',file),path.join(dest,'examples',file));
await cp(path.join(root,'examples/compare'),path.join(dest,'examples/compare'),{recursive:true});
await cp(path.join(root,'docs/PORTABLE.md'),path.join(dest,'README.md'));
const hash=async file=>createHash('sha256').update(await readFile(file)).digest('hex');
const identity={format:2,version,platform:goos,arch:goarch,wails,binarySHA256:await hash(binary),rendererSHA256:await hash(helper)};
await writeFile(path.join(dest,'BUILD.json'),JSON.stringify({...identity,created:stamp,publicRelease:false,runtimeVerified:false},null,2)+'\n');
if(values.release){
  await writeFile(path.join(dest,'RELEASE.json'),JSON.stringify(identity,null,2)+'\n');
  const archive=dest+(platform==='win32'?'.zip':'.tar.gz');
  execFileSync('go',['run','./scripts/archive',dest,archive],{cwd:root,stdio:'inherit'});
  console.log(archive);
}else console.log(dest);
