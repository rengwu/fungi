const { contextBridge, ipcRenderer } = require('electron');
async function encode(items=[]) {
  return Promise.all(items.map(async({key,name,blob})=>{
    if (!(blob instanceof Blob)) throw new Error('attachment.blob must be a Blob or File');
    if(blob.size>20*1024*1024) throw new Error('Attachment exceeds 20 MiB');
    return {key,name,mediaType:blob.type,bytes:new Uint8Array(await blob.arrayBuffer())};
  }));
}
contextBridge.exposeInMainWorld('fungi',{
  ready:async()=>{
    const ready=await ipcRenderer.invoke('fungi:ready');
    ready.attachments=ready.attachments.map(a=>({key:a.key,name:a.name,blob:new Blob([a.bytes],{type:a.mediaType})}));
    return ready;
  },
  saveDraft:async(data,items)=>ipcRenderer.invoke('fungi:save',data,await encode(items)),
  submit:async(data,items)=>ipcRenderer.invoke('fungi:submit',data,await encode(items)),
  cancel:()=>ipcRenderer.invoke('fungi:cancel')
});
