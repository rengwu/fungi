// Bounded native-window checks through the public CLI. Requires a desktop
// large enough for the default 820x780 window and native window decorations.
import {mkdtemp,mkdir,writeFile,readFile} from 'node:fs/promises';
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import path from 'node:path';
import assert from 'node:assert/strict';
import {fileURLToPath} from 'node:url';

assert.equal(process.argv.length,3,'usage: node scripts/window-smoke.mjs /absolute/package/fungi');
const binary=path.resolve(process.argv[2]);
const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
await mkdir(path.join(root,'verification-output'),{recursive:true});
const directory=await mkdtemp(path.join(root,'verification-output','window-smoke-'));
const source=path.join(directory,'index.html');
await writeFile(source,`<!doctype html><meta charset="utf-8"><title>Fungi window size verification</title>
<h1>Window size verification</h1><p>This disposable window closes automatically.</p>
<script type="module">
const {draft}=await fungi.ready();
await new Promise(resolve=>requestAnimationFrame(()=>requestAnimationFrame(resolve)));
const size={width:innerWidth,height:innerHeight,outerWidth,outerHeight,devicePixelRatio};
if(draft) await fungi.submit({initial:draft,current:size});
else {await fungi.saveDraft(size);await fungi.cancel();}
</script>`);
const env={...process.env,FUNGI_STATE_DIR:path.join(directory,'sessions')};
delete env.CODEX_THREAD_ID;
const evidence={binary,directory,platform:process.platform,checks:[]};
const run=async args=>{
  const {stdout}=await promisify(execFile)(binary,args,{env,timeout:30000,maxBuffer:1024*1024});
  return JSON.parse(stdout);
};
let baseline;
try {
  for(const [name,flags,width,height] of [
    ['defaults',[],820,780],
    ['both',['--width','960','--height','640'],960,640],
    ['width only',['--width','700'],700,780],
    ['height only',['--height','600'],820,600],
    ['minimum',['--width','420','--height','420'],420,420],
  ]) {
    const closed=await run(['open',source,...flags,'--delivery','wait']);
    assert.equal(closed.outcome,'cancelled');
    const result=await run(['reopen',closed.session,'--delivery','wait']);
    assert.equal(result.outcome,'submitted');
    evidence.checks.push({name,closed,result});
    const {initial,current}=result.data;
    baseline??=initial;
    // Native frames differ by platform; compare the actual viewport change.
    assert.equal(initial.width-baseline.width,width-820,name+' width');
    assert.equal(initial.height-baseline.height,height-780,name+' height');
    assert.deepEqual(current,initial,name+' reopen');
    if(name==='both') {
      const override=await run(['reopen',closed.session,'--width','740','--delivery','wait']);
      evidence.checks.push({name:'reopen width override',result:override});
      assert.equal(override.outcome,'submitted');
      assert.equal(override.data.current.width-baseline.width,740-820);
      assert.equal(override.data.current.height,initial.height);
      assert.deepEqual(override.data.initial,initial);
    }
  }
  evidence.passed=true;
} finally {
  evidence.build=JSON.parse(await readFile(path.join(path.dirname(binary),'BUILD.json'),'utf8'));
  await writeFile(path.join(directory,'evidence.json'),JSON.stringify(evidence,null,2)+'\n');
  console.log(JSON.stringify({evidence:path.join(directory,'evidence.json'),passed:!!evidence.passed,checks:evidence.checks.length}));
}
