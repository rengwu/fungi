const { app, BrowserWindow, ipcMain, protocol, session, dialog, Menu } = require('electron');
const fs = require('node:fs/promises');
const path = require('node:path');
const readline = require('node:readline');
const crypto = require('node:crypto');
const {serveAsset}=require('./assets.cjs');

protocol.registerSchemesAsPrivileged([{scheme:'fungi',privileges:{standard:true,secure:true,supportFetchAPI:true,corsEnabled:true,stream:true}}]);
app.setName('Fungi');
const privateData=process.argv.find(a=>a.startsWith('--fungi-data='));
if(privateData)app.setPath('userData',privateData.slice('--fungi-data='.length));
// Electron's Windows initialization replaces process.stdin with an EOF-only
// stream. Read the inherited pipe directly without changing the protocol.
const input=process.platform==='win32'?require('node:fs').createReadStream(null,{fd:0,autoClose:false}):process.stdin;
input.on('error',error=>{console.error(error.message);app.exit(1);});
const lines = readline.createInterface({input});
let launch, win, terminal = false, closing = false, sequence = 0;
let queue = Promise.resolve();
const pending = new Map();
function send(value) { process.stdout.write(JSON.stringify(value)+'\n'); }
function request(value) {
  const id = String(++sequence);
  const line=JSON.stringify({...value,id})+'\n';
  if(Buffer.byteLength(line)>launch.limits.jsonBytes)throw new Error('Response envelope exceeds the JSON limit');
  return new Promise((resolve,reject)=>{
    pending.set(id,{resolve,reject});
    process.stdout.write(line);
  });
}
lines.on('line', async line => {
  try {
    const value=JSON.parse(line);
    if (!launch) {
      launch=value;
      if (launch.protocol!==1) throw new Error('Unsupported Fungi protocol');
      await openWindow();
      return;
    }
    const task=pending.get(value.id);
    if (task) {
      pending.delete(value.id);
      if (value.ok) task.resolve(value.result);
      else task.reject(new Error(value.error?.message || 'Unable to save'));
    }
  } catch(error) { console.error(error.message); app.exit(1); }
});
lines.on('close',()=>{
  if (!terminal) { console.error('Fungi owner disconnected; last acknowledged draft is recoverable.'); app.exit(1); }
});
function serialized(action) {
  const result=queue.then(action);
  queue=result.catch(()=>{});
  return result;
}
function assertSender(event) {
  if (!win || event.sender!==win.webContents || event.senderFrame!==win.webContents.mainFrame) throw new Error('Unsupported frame');
}
async function attachments(items = []) {
  if (!Array.isArray(items)) throw new Error('attachments must be an array');
  let total=0;const keys=new Set();const written=[];
  try {
    for (const item of items) {
      if (typeof item.key!=='string' || !item.key || keys.has(item.key) || typeof item.name!=='string' || !item.name) throw new Error('Attachments need unique keys and names');
      keys.add(item.key);
      const bytes=Buffer.from(item.bytes);
      total+=bytes.length;
      if (bytes.length>launch.limits.artifactBytes || total>launch.limits.totalArtifactBytes) throw new Error('Attachment size limit exceeded');
      const file=crypto.randomUUID();
      await fs.writeFile(path.join(launch.outputDir,file),bytes,{mode:0o600,flag:'wx'});
      written.push({key:item.key,name:path.basename(item.name),mediaType:item.mediaType||'application/octet-stream',path:file});
    }
    return written;
  } catch(e) { await cleanup(written);throw e; }
}
async function cleanup(items) { await Promise.all(items.map(x=>fs.rm(path.join(launch.outputDir,x.path),{force:true}))); }
async function persist(type, data, items, outcome) {
  if (Buffer.byteLength(JSON.stringify(data ?? null))>launch.limits.jsonBytes) throw new Error('Response JSON exceeds 1 MiB');
  const files=await attachments(items);
  try { return await request({type,data,artifacts:files,outcome}); }
  finally { await cleanup(files); }
}
async function openWindow() {
  await app.whenReady();
  Menu.setApplicationMenu(Menu.buildFromTemplate([
    ...(process.platform==='darwin'?[{role:'appMenu'}]:[{label:'File',submenu:[{role:'close'}]}]),
    {role:'editMenu'},{role:'windowMenu'}
  ]));
  const sourceRoot=await fs.realpath(path.dirname(launch.source));
  const builtinRoot=await fs.realpath(__dirname);
  const form=path.extname(launch.source).toLowerCase()==='.json' ? JSON.parse(await fs.readFile(launch.source,'utf8')) : null;
  // No Node.js, file:// browsing, remote navigation, or permission prompts in a
  // skill's page. Local assets are served beneath one explicit directory.
  protocol.handle('fungi',request=>serveAsset(request,{app:builtinRoot,skill:sourceRoot}));
  session.defaultSession.setPermissionRequestHandler((_wc,_p,callback)=>callback(false));
  session.defaultSession.setPermissionCheckHandler(()=>false);
  session.defaultSession.webRequest.onBeforeRequest((details,callback)=>{
    callback({cancel:details.resourceType==='subFrame'||!['fungi:','data:','blob:','devtools:'].some(p=>details.url.startsWith(p))});
  });
  win=new BrowserWindow({
    width:820,height:780,minWidth:420,minHeight:420,show:false,title:form?.title || 'Fungi',
    backgroundColor:'#18191c',
    webPreferences:{preload:path.join(__dirname,'preload.cjs'),nodeIntegration:false,contextIsolation:true,sandbox:true,webSecurity:true}
  });
  win.setMenu(null);
  win.webContents.setWindowOpenHandler(()=>({action:'deny'}));
  win.webContents.on('will-navigate',event=>event.preventDefault());
  // Inline srcdoc previews need no document request. Keep file/URL frame
  // navigation blocked; sandboxing the preview is the authoring skill's job.
  win.webContents.on('will-frame-navigate',event=>{
    if(event.isMainFrame||event.url!=='about:srcdoc')event.preventDefault();
  });
  win.webContents.on('render-process-gone',()=>{if (!terminal) serialized(()=>request({type:'result',outcome:'failed',error:{code:'renderer_crashed',message:'Page renderer stopped'}})).finally(()=>app.exit(1));});
  ipcMain.handle('fungi:ready',async event=>{
    assertSender(event);
    const restored=[];
    for (const a of launch.draftArtifacts||[]) restored.push({...a,bytes:new Uint8Array(await fs.readFile(a.path)),path:undefined});
    return {input:launch.input,draft:launch.draft,attachments:restored,form,limits:launch.limits,delivery:launch.delivery};
  });
  ipcMain.handle('fungi:save', (event,data,items)=>{
    assertSender(event);
    return serialized(()=>persist('draft',data,items));
  });
  ipcMain.handle('fungi:submit',(event,data,items)=>{
    assertSender(event);
    return serialized(async()=>{
      const result=await persist('result',data,items,'submitted');
      terminal=true;
      setImmediate(()=>win.close());
      return result;
    });
  });
  ipcMain.handle('fungi:cancel',event=>{
    assertSender(event);
    return serialized(async()=>{
      const result=await request({type:'result',outcome:'cancelled'});
      terminal=true;setImmediate(()=>win.close());return result;
    });
  });
  // Reference form autosaves in order. Normal close waits for those writes.
  // Custom pages can opt into acknowledged saveDraft; no DOM capture is implied.
  win.on('close',event=>{
    if (terminal) return;
    event.preventDefault();
    if (closing) return;
    closing=true;
    serialized(()=>request({type:'result',outcome:'closed'})).then(()=>{
      terminal=true;win.close();
    }).catch(async error=>{
      closing=false;
      await dialog.showMessageBox(win,{type:'error',message:'Fungi could not save the close outcome',detail:error.message,buttons:['Keep open']});
    });
  });
  win.on('closed',()=>app.quit());
  const url=form?'fungi://app/form.html':'fungi://skill/'+encodeURIComponent(path.basename(launch.source));
  await win.loadURL(url);
  win.show();win.focus();
  send({type:'ready'});
}
