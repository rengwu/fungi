// Included inside the generated form's async scope by adapt.mjs.
let fungiSaving=Promise.resolve(), fungiGeneration=0, fungiBusy=false;
const imageMimeExtensions={'image/png':'png','image/jpeg':'jpg','image/gif':'gif','image/webp':'webp'};

async function restoreFungiDraft(state) {
  const draft=fungiReady.draft;
  if(!draft||draft.specDigest!==state.specDigest||draft.sessionId!==state.sessionId)return null;
  const restored=structuredClone(draft.answers);
  const files=new Map(fungiReady.attachments.map(a=>[a.key,a.blob]));
  for(const note of restored.flatMap(a=>[a,...(a.comments??[])])){
    for(const image of note.images??[]){
      const blob=files.get(image.artifact);
      if(!blob)throw new Error('Saved inline image is missing');
      image.dataUrl=await new Promise((resolve,reject)=>{
        const reader=new FileReader();
        reader.onload=()=>resolve(reader.result);reader.onerror=()=>reject(reader.error);
        reader.readAsDataURL(blob);
      });
      delete image.artifact;delete image.mimeType;
    }
  }
  return restored;
}
function fungiPayload() {
  const copy=structuredClone(answers), attachments=[];
  for(const note of copy.flatMap(a=>[a,...(a.comments??[])])){
    for(const image of note.images??[]){
      const match=image.dataUrl?.match(/^data:(image\/(?:png|jpeg|gif|webp));base64,([A-Za-z0-9+/]+={0,2})$/);
      if(!match)throw new Error('An inline image is not ready');
      const key='image-'+image.id;
      const bytes=Uint8Array.from(atob(match[2]),c=>c.charCodeAt(0));
      attachments.push({key,name:key+'.'+imageMimeExtensions[match[1]],blob:new Blob([bytes],{type:match[1]})});
      delete image.dataUrl;
      image.artifact=key;image.mimeType=match[1];
    }
  }
  return {data:{schemaVersion:1,sessionId:state.sessionId,specDigest:state.specDigest,title:state.spec.title,answers:copy},attachments};
}
function showFungiError(error) {
  message.className='error';message.textContent=error.message;message.hidden=false;
}
function saveFungiDraft() {
  submit.disabled=fungiBusy||!canSubmit()||!notesReady();
  if(state.archived||fungiBusy||!notesReady())return;
  const warning=document.getElementById('draft-message'), generation=++fungiGeneration;
  warning.textContent='Saving draft…';warning.hidden=false;
  let payload;
  try{payload=fungiPayload()}catch(error){showFungiError(error);return}
  fungiSaving=fungiSaving.catch(()=>{}).then(()=>fungi.saveDraft(payload.data,payload.attachments)).then(()=>{
    if(generation===fungiGeneration)warning.textContent='Draft saved in Fungi';
  }).catch(error=>{
    if(generation===fungiGeneration){warning.textContent='Draft not saved';showFungiError(error)}
    throw error;
  });
  fungiSaving.catch(()=>{});
}
async function submitFungiForm(event) {
  event.preventDefault();
  if(fungiBusy||!canSubmit()||!notesReady())return;
  fungiBusy=true;submit.disabled=true;submit.textContent='Submitting…';form.inert=true;message.hidden=true;
  try{
    const payload=fungiPayload();
    await fungiSaving.catch(()=>{});
    await fungi.submit(payload.data,payload.attachments);
  }catch(error){
    fungiBusy=false;form.inert=false;submit.textContent='Submit feedback';
    submit.disabled=!canSubmit()||!notesReady();showFungiError(error);message.focus();
  }
}
