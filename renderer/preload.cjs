const { contextBridge, ipcRenderer } = require('electron');
let initialization, queue=Promise.resolve();
function configuration(){return initialization??=ipcRenderer.invoke('fungi:ready');}
// Order starts before asynchronous Blob encoding, not only at main-process IPC.
// Failed writes reject their own caller but must not prevent a later retry.
function serialized(action){
  const result=queue.then(action);
  queue=result.catch(()=>{});
  return result;
}
async function encode(items=[]) {
  if(!Array.isArray(items))throw new Error('attachments must be an array');
  const {limits}=await configuration();
  let total=0;
  const keys=new Set();
  for(const {key,name,blob} of items){
    if(typeof key!=='string'||!key||keys.has(key)||typeof name!=='string'||!name)throw new Error('Attachments need unique keys and names');
    keys.add(key);
    if (!(blob instanceof Blob)) throw new Error('attachment.blob must be a Blob or File');
    total+=blob.size;
    if(blob.size>limits.artifactBytes||total>limits.totalArtifactBytes)throw new Error('Attachment size limit exceeded');
  }
  return Promise.all(items.map(async({key,name,blob})=>{
    return {key,name,mediaType:blob.type,bytes:new Uint8Array(await blob.arrayBuffer())};
  }));
}
contextBridge.exposeInMainWorld('fungi',{
  ready:async()=>{
    const ready=await configuration();
    return {...ready,attachments:ready.attachments.map(a=>({key:a.key,name:a.name,blob:new Blob([a.bytes],{type:a.mediaType})}))};
  },
  saveDraft:(data,items)=>serialized(async()=>ipcRenderer.invoke('fungi:save',data,await encode(items))),
  submit:(data,items)=>serialized(async()=>ipcRenderer.invoke('fungi:submit',data,await encode(items))),
  cancel:()=>serialized(()=>ipcRenderer.invoke('fungi:cancel'))
});
