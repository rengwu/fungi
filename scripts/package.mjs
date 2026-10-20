// Development-only packaging. The resulting directory requires neither Node
// nor Go on the user's machine. Uses the downloaded, pinned Electron runtime.
import { cp, mkdir, readFile, writeFile, stat } from 'node:fs/promises';
import { execFileSync } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
if(process.platform!=='darwin')throw new Error('This packaging recipe is verified only on macOS');
const stamp=new Date().toISOString().replace(/[:.]/g,'-');
const dest=path.join(root,'dist',`fungi-darwin-${process.arch}-${stamp}`);
const renderer=path.join(dest,'renderer');
await mkdir(renderer,{recursive:true});
execFileSync('go',['build','-trimpath','-o',path.join(dest,'fungi'),'.'],{cwd:root,stdio:'inherit'});
const runtime=path.join(root,'renderer/node_modules/electron/dist');
await stat(path.join(runtime,'Electron.app'));
await cp(runtime,path.join(renderer,'runtime'),{recursive:true,verbatimSymlinks:true});
for(const file of ['package.json','main.cjs','preload.cjs','form.html','form.js'])await cp(path.join(root,'renderer',file),path.join(renderer,file));
// Keep upstream signatures untouched; this development archive is not a
// notarized public release. Its runtime and app code have separate locations.
await mkdir(path.join(dest,'examples'));
for(const file of ['feedback.json','review.html'])await cp(path.join(root,'examples',file),path.join(dest,'examples',file));
execFileSync('go',['build','-trimpath','-o',path.join(dest,'examples/helper'),'./examples/helper'],{cwd:root,stdio:'inherit'});
await cp(path.join(root,'README.md'),path.join(dest,'README.md'));
await cp(path.join(root,'skills'),path.join(dest,'skills'),{recursive:true,verbatimSymlinks:true});
await cp(path.join(root,'docs/wiki'),path.join(dest,'docs/wiki'),{recursive:true,verbatimSymlinks:true});
const pkg=JSON.parse(await readFile(path.join(root,'renderer/package.json'),'utf8'));
await writeFile(path.join(dest,'BUILD.json'),JSON.stringify({created:stamp,platform:process.platform,arch:process.arch,electron:pkg.devDependencies.electron,publicRelease:false},null,2)+'\n');
console.log(dest);
